package matcher

import (
	"math"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/complai/complai/services/go/recon-service/internal/domain"
)

// Run executes the 5-stage match pipeline:
// Stage 1: Exact match — invoice_number exact + vendor GSTIN exact + amount diff < 1 + date exact
// Stage 2: Fuzzy match — DL distance <=2 on invoice_number + same vendor + amount +/-5% + date +/-5 days
// Stage 3: AI-stub — placeholder (Phase 4+)
// Stage 4: Partial match — same vendor + amounts that sum to PR amount (split invoices)
// Stage 5: Unmatched — anything left goes to Missing-2B or Missing-PR
func Run(pr []domain.PurchaseRegisterEntry, gstr2b []domain.GSTR2BEntry, gstin, returnPeriod string, runID uuid.UUID) []domain.ReconMatch {
	var results []domain.ReconMatch

	// Build index maps for quick lookups
	prUsed := make(map[int]bool)
	gstr2bUsed := make(map[int]bool)

	// Pre-normalize invoice numbers for comparison
	prNorm := make([]string, len(pr))
	for i, p := range pr {
		prNorm[i] = normalizeInvoiceNumber(p.InvoiceNumber)
	}
	gstr2bNorm := make([]string, len(gstr2b))
	for i, g := range gstr2b {
		gstr2bNorm[i] = normalizeInvoiceNumber(g.InvoiceNumber)
	}

	// Stage 1: Exact match
	for pi, p := range pr {
		if prUsed[pi] {
			continue
		}
		for gi, g := range gstr2b {
			if gstr2bUsed[gi] {
				continue
			}
			if prNorm[pi] == gstr2bNorm[gi] &&
				p.VendorGSTIN == g.SupplierGSTIN &&
				p.TotalValue.Sub(g.TotalValue).Abs().LessThan(decimal.NewFromInt(1)) &&
				p.InvoiceDate == g.InvoiceDate {

				reasons := []string{"exact_invoice_number", "exact_gstin", "exact_amount", "exact_date"}
				m := buildMatch(runID, gstin, returnPeriod, &p, &g, domain.MatchTypeDirect, decimal.NewFromFloat(1.0), reasons)
				results = append(results, m)
				prUsed[pi] = true
				gstr2bUsed[gi] = true
				break
			}
		}
	}

	// Stage 2: Fuzzy match
	for pi, p := range pr {
		if prUsed[pi] {
			continue
		}
		bestGI := -1
		bestConf := decimal.Zero
		var bestReasons []string

		for gi, g := range gstr2b {
			if gstr2bUsed[gi] {
				continue
			}
			if p.VendorGSTIN != g.SupplierGSTIN {
				continue
			}

			dist := damerauLevenshtein(prNorm[pi], gstr2bNorm[gi])
			if dist > 2 {
				continue
			}

			amountDiffPct := amountDiffPercent(p.TotalValue, g.TotalValue)
			if amountDiffPct > 5.0 {
				continue
			}

			dateDiff := dateDiffDays(p.InvoiceDate, g.InvoiceDate)
			if dateDiff > 5 {
				continue
			}

			// Compute confidence score
			conf := computeFuzzyConfidence(dist, amountDiffPct, dateDiff)
			var reasons []string
			if dist == 0 {
				reasons = append(reasons, "exact_invoice_number")
			} else {
				reasons = append(reasons, "fuzzy_invoice_number")
			}
			reasons = append(reasons, "exact_gstin")
			if amountDiffPct == 0 {
				reasons = append(reasons, "exact_amount")
			} else {
				reasons = append(reasons, "approx_amount")
			}
			if dateDiff == 0 {
				reasons = append(reasons, "exact_date")
			} else {
				reasons = append(reasons, "approx_date")
			}

			if conf.GreaterThan(bestConf) {
				bestGI = gi
				bestConf = conf
				bestReasons = reasons
			}
		}

		if bestGI >= 0 {
			m := buildMatch(runID, gstin, returnPeriod, &p, &gstr2b[bestGI], domain.MatchTypeProbable, bestConf, bestReasons)
			results = append(results, m)
			prUsed[pi] = true
			gstr2bUsed[bestGI] = true
		}
	}

	// Stage 3: AI-stub — placeholder, no matches generated (Phase 4+)

	// Stage 4: Partial match — same vendor + amounts that sum to PR amount (split invoices)
	for pi, p := range pr {
		if prUsed[pi] {
			continue
		}

		// Find unmatched GSTR-2B entries from the same vendor
		var candidateGIs []int
		for gi, g := range gstr2b {
			if gstr2bUsed[gi] {
				continue
			}
			if p.VendorGSTIN == g.SupplierGSTIN {
				candidateGIs = append(candidateGIs, gi)
			}
		}

		if len(candidateGIs) < 2 {
			continue
		}

		// Try to find a subset whose total sums to the PR amount (+/-1)
		sumGIs := findSubsetSum(gstr2b, candidateGIs, p.TotalValue)
		if len(sumGIs) > 0 {
			// Create a partial match for each 2B entry in the subset
			for _, gi := range sumGIs {
				g := gstr2b[gi]
				reasons := []string{"partial_amount_match", "same_vendor", "split_invoice"}
				m := buildMatch(runID, gstin, returnPeriod, &p, &g, domain.MatchTypePartial, decimal.NewFromFloat(0.70), reasons)
				results = append(results, m)
				gstr2bUsed[gi] = true
			}
			prUsed[pi] = true
		}
	}

	// Stage 5: Unmatched
	for pi, p := range pr {
		if prUsed[pi] {
			continue
		}
		reasons := []string{"no_matching_2b_entry"}
		m := buildMatch(runID, gstin, returnPeriod, &p, nil, domain.MatchTypeMissing2B, decimal.Zero, reasons)
		results = append(results, m)
	}

	for gi, g := range gstr2b {
		if gstr2bUsed[gi] {
			continue
		}
		reasons := []string{"no_matching_pr_entry"}
		m := buildMatch(runID, gstin, returnPeriod, nil, &g, domain.MatchTypeMissingPR, decimal.Zero, reasons)
		results = append(results, m)
	}

	return results
}

func buildMatch(runID uuid.UUID, gstin, returnPeriod string, pr *domain.PurchaseRegisterEntry, gstr2b *domain.GSTR2BEntry, matchType domain.MatchType, confidence decimal.Decimal, reasons []string) domain.ReconMatch {
	m := domain.ReconMatch{
		ID:              uuid.New(),
		RunID:           runID,
		GSTIN:           gstin,
		ReturnPeriod:    returnPeriod,
		MatchType:       matchType,
		MatchConfidence: confidence,
		ReasonCodes:     reasons,
		Status:          domain.MatchStatusUnreviewed,
		PRAmount:        decimal.Zero,
		GSTR2BAmount:    decimal.Zero,
	}

	if pr != nil {
		m.PRInvoiceNumber = pr.InvoiceNumber
		m.PRInvoiceDate = pr.InvoiceDate
		m.PRVendorGSTIN = pr.VendorGSTIN
		m.PRAmount = pr.TotalValue
		m.PRHSN = pr.HSN
		m.PRSourceID = pr.SourceID
	}

	if gstr2b != nil {
		m.GSTR2BInvoiceNumber = gstr2b.InvoiceNumber
		m.GSTR2BInvoiceDate = gstr2b.InvoiceDate
		m.GSTR2BSupplierGSTIN = gstr2b.SupplierGSTIN
		m.GSTR2BAmount = gstr2b.TotalValue
		m.GSTR2BHSN = gstr2b.HSN
	}

	return m
}

// normalizeInvoiceNumber strips hyphens, slashes, spaces, and converts to uppercase.
func normalizeInvoiceNumber(inv string) string {
	inv = strings.ToUpper(inv)
	inv = strings.ReplaceAll(inv, "-", "")
	inv = strings.ReplaceAll(inv, "/", "")
	inv = strings.ReplaceAll(inv, " ", "")
	return inv
}

// damerauLevenshtein computes the Damerau-Levenshtein distance between two strings.
func damerauLevenshtein(s, t string) int {
	sLen := len(s)
	tLen := len(t)

	if sLen == 0 {
		return tLen
	}
	if tLen == 0 {
		return sLen
	}

	// Early exit if length difference > 2
	if abs(sLen-tLen) > 2 {
		return abs(sLen - tLen)
	}

	d := make([][]int, sLen+1)
	for i := range d {
		d[i] = make([]int, tLen+1)
		d[i][0] = i
	}
	for j := 0; j <= tLen; j++ {
		d[0][j] = j
	}

	for i := 1; i <= sLen; i++ {
		for j := 1; j <= tLen; j++ {
			cost := 0
			if s[i-1] != t[j-1] {
				cost = 1
			}

			d[i][j] = min3(
				d[i-1][j]+1,      // deletion
				d[i][j-1]+1,      // insertion
				d[i-1][j-1]+cost, // substitution
			)

			// Transposition
			if i > 1 && j > 1 && s[i-1] == t[j-2] && s[i-2] == t[j-1] {
				d[i][j] = min2(d[i][j], d[i-2][j-2]+cost)
			}
		}
	}

	return d[sLen][tLen]
}

func min2(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func min3(a, b, c int) int {
	return min2(min2(a, b), c)
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func amountDiffPercent(a, b decimal.Decimal) float64 {
	if a.IsZero() && b.IsZero() {
		return 0
	}
	diff := a.Sub(b).Abs()
	base := a.Abs()
	if base.IsZero() {
		base = b.Abs()
	}
	pct, _ := diff.Div(base).Mul(decimal.NewFromInt(100)).Float64()
	return pct
}

func dateDiffDays(dateA, dateB string) int {
	layouts := []string{"02-01-2006", "02/01/2006", "2006-01-02"}

	var tA, tB time.Time
	var errA, errB error

	for _, layout := range layouts {
		tA, errA = time.Parse(layout, dateA)
		if errA == nil {
			break
		}
	}
	for _, layout := range layouts {
		tB, errB = time.Parse(layout, dateB)
		if errB == nil {
			break
		}
	}

	if errA != nil || errB != nil {
		// If we can't parse dates, treat as no match on date
		return 999
	}

	diff := tA.Sub(tB).Hours() / 24
	return int(math.Abs(diff))
}

func computeFuzzyConfidence(editDist int, amountDiffPct float64, dateDiff int) decimal.Decimal {
	// Start at 0.95 (perfect fuzzy match) and reduce based on deviations
	conf := 0.95

	// Reduce for edit distance
	conf -= float64(editDist) * 0.05

	// Reduce for amount difference
	conf -= amountDiffPct * 0.02

	// Reduce for date difference
	conf -= float64(dateDiff) * 0.02

	if conf < 0.50 {
		conf = 0.50
	}

	return decimal.NewFromFloat(conf)
}

// findSubsetSum tries to find a subset of gstr2b entries (by index) whose TotalValue sums
// to the target amount (+/- 1). Uses a simple greedy approach for efficiency.
func findSubsetSum(entries []domain.GSTR2BEntry, candidateGIs []int, target decimal.Decimal) []int {
	// Greedy: sort candidates by amount descending and accumulate
	// For practical purposes (performance), limit to first 20 candidates
	limit := len(candidateGIs)
	if limit > 20 {
		limit = 20
	}
	candidates := candidateGIs[:limit]

	// Try all pairs first (most common split invoice case)
	for i := 0; i < len(candidates); i++ {
		for j := i + 1; j < len(candidates); j++ {
			sum := entries[candidates[i]].TotalValue.Add(entries[candidates[j]].TotalValue)
			if sum.Sub(target).Abs().LessThan(decimal.NewFromInt(1)) {
				return []int{candidates[i], candidates[j]}
			}
		}
	}

	// Try triples
	if len(candidates) >= 3 {
		for i := 0; i < len(candidates); i++ {
			for j := i + 1; j < len(candidates); j++ {
				for k := j + 1; k < len(candidates); k++ {
					sum := entries[candidates[i]].TotalValue.
						Add(entries[candidates[j]].TotalValue).
						Add(entries[candidates[k]].TotalValue)
					if sum.Sub(target).Abs().LessThan(decimal.NewFromInt(1)) {
						return []int{candidates[i], candidates[j], candidates[k]}
					}
				}
			}
		}
	}

	return nil
}
