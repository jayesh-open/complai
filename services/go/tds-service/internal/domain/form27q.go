package domain

import (
	"fmt"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

type NonResidentDetail struct {
	DeducteeID       string          `json:"deductee_id"`
	PAN              string          `json:"pan"`
	Name             string          `json:"name"`
	CountryCode      string          `json:"country_code"`
	TaxResidency     string          `json:"tax_residency"`
	DTAAArticle      string          `json:"dtaa_article"`
	DTAARate         decimal.Decimal `json:"dtaa_rate"`
	GrossAmount      decimal.Decimal `json:"gross_amount"`
	GrossAmountFC    decimal.Decimal `json:"gross_amount_fc"`
	CurrencyCode     string          `json:"currency_code"`
	ExchangeRate     decimal.Decimal `json:"exchange_rate"`
	TDSRate          decimal.Decimal `json:"tds_rate"`
	TDSAmount        decimal.Decimal `json:"tds_amount"`
	Surcharge        decimal.Decimal `json:"surcharge"`
	Cess             decimal.Decimal `json:"cess"`
	TotalTax         decimal.Decimal `json:"total_tax"`
	NatureOfRemittance string       `json:"nature_of_remittance"`
	DateOfPayment    time.Time       `json:"date_of_payment"`
	DateOfDeduction  time.Time       `json:"date_of_deduction"`
	ChallanNumber    string          `json:"challan_number"`
	BSRCode          string          `json:"bsr_code"`
}

type Form27QPayload struct {
	FormType      FormType            `json:"form_type"`
	FinancialYear string              `json:"financial_year"`
	Quarter       string              `json:"quarter"`
	Deductor      DeductorDetails     `json:"deductor"`
	Remittances   []NonResidentDetail `json:"remittances"`
	TotalTDS      decimal.Decimal     `json:"total_tds"`
	TotalPaid     decimal.Decimal     `json:"total_paid"`
	CreatedAt     time.Time           `json:"created_at"`
	Errors        []string            `json:"errors,omitempty"`
}

type Form27QInput struct {
	Deductor      DeductorDetails
	FinancialYear string
	Quarter       string
	Deductees     []Deductee
	Entries       []TDSEntry
	CountryCodes  map[string]string
	DTAAArticles  map[string]string
	DTAARates     map[string]decimal.Decimal
	CurrencyCodes map[string]string
	ExchangeRates map[string]decimal.Decimal
	ForeignAmounts map[string]decimal.Decimal
}

func GenerateForm27Q(input Form27QInput) (*Form27QPayload, error) {
	if input.Deductor.TAN == "" {
		return nil, fmt.Errorf("deductor TAN is required")
	}
	if input.FinancialYear == "" || input.Quarter == "" {
		return nil, fmt.Errorf("financial_year and quarter are required")
	}

	s195Entries := filterBySection(input.Entries, Section195)
	if len(s195Entries) == 0 {
		return nil, fmt.Errorf("no section 195 entries found for %s %s", input.FinancialYear, input.Quarter)
	}

	deducteeMap := make(map[string]*Deductee)
	for i := range input.Deductees {
		deducteeMap[input.Deductees[i].ID.String()] = &input.Deductees[i]
	}

	payload := &Form27QPayload{
		FormType:      FormType27Q,
		FinancialYear: input.FinancialYear,
		Quarter:       input.Quarter,
		Deductor:      input.Deductor,
		CreatedAt:     time.Now(),
	}

	totalTDS := decimal.Zero
	totalPaid := decimal.Zero
	var validationErrors []string

	for _, entry := range s195Entries {
		d, ok := deducteeMap[entry.DeducteeID.String()]
		if !ok {
			validationErrors = append(validationErrors, fmt.Sprintf("deductee %s not found", entry.DeducteeID))
			continue
		}
		if d.ResidentStatus != NonResident {
			validationErrors = append(validationErrors, fmt.Sprintf("deductee %s (%s) is not non-resident", d.Name, entry.DeducteeID))
		}

		deducteeIDStr := entry.DeducteeID.String()
		countryCode := lookupDefault(input.CountryCodes, deducteeIDStr, "")
		taxResidency := countryCode
		dtaaArticle := lookupDefault(input.DTAAArticles, deducteeIDStr, "")
		dtaaRate := lookupDecimal(input.DTAARates, deducteeIDStr, decimal.Zero)
		currencyCode := lookupDefault(input.CurrencyCodes, deducteeIDStr, "INR")
		exchangeRate := lookupDecimal(input.ExchangeRates, deducteeIDStr, decimal.NewFromInt(1))
		foreignAmount := lookupDecimal(input.ForeignAmounts, deducteeIDStr, decimal.Zero)

		detail := NonResidentDetail{
			DeducteeID:         deducteeIDStr,
			PAN:                d.PAN,
			Name:               d.Name,
			CountryCode:        countryCode,
			TaxResidency:       taxResidency,
			DTAAArticle:        dtaaArticle,
			DTAARate:           dtaaRate,
			GrossAmount:        entry.GrossAmount,
			GrossAmountFC:      foreignAmount,
			CurrencyCode:       currencyCode,
			ExchangeRate:       exchangeRate,
			TDSRate:            entry.TDSRate,
			TDSAmount:          entry.TDSAmount,
			Surcharge:          entry.Surcharge,
			Cess:               entry.Cess,
			TotalTax:           entry.TotalTax,
			NatureOfRemittance: entry.NatureOfPayment,
			DateOfPayment:      entry.TransactionDate,
			DateOfDeduction:    entry.TransactionDate,
			ChallanNumber:      entry.ChallanNumber,
			BSRCode:            entry.BSRCode,
		}
		payload.Remittances = append(payload.Remittances, detail)
		totalTDS = totalTDS.Add(entry.TotalTax)
		totalPaid = totalPaid.Add(entry.GrossAmount)
	}

	payload.TotalTDS = totalTDS
	payload.TotalPaid = totalPaid
	payload.Errors = validationErrors

	return payload, nil
}

func GenerateForm27QFVU(payload *Form27QPayload) string {
	var b strings.Builder

	ay := assessmentYear(payload.FinancialYear)

	b.WriteString(fmt.Sprintf("^FH^27Q^1^%s^%s^%s^%s^^%s^^%d^^%s^%s^\n",
		payload.Deductor.TAN,
		payload.Deductor.DeductorPAN,
		ay,
		payload.Quarter,
		payload.Deductor.DeductorName,
		len(payload.Remittances),
		payload.Deductor.Address,
		payload.Deductor.Pincode,
	))

	for _, ch := range challanSummary27Q(payload) {
		b.WriteString(fmt.Sprintf("^BH^%s^%s^%s^%s^%s^\n",
			ch.ChallanNumber, ch.BSRCode, ch.DepositDate,
			ch.TotalTDS, ch.DeducteeCount,
		))
	}

	for i, rem := range payload.Remittances {
		b.WriteString(fmt.Sprintf("^NR^%d^%s^%s^%s^%s^%s^%s^%s^%s^%s^%s^%s^%s^%s^\n",
			i+1,
			rem.PAN,
			rem.Name,
			rem.CountryCode,
			rem.DTAAArticle,
			rem.DTAARate.StringFixed(4),
			rem.GrossAmount.StringFixed(2),
			rem.GrossAmountFC.StringFixed(2),
			rem.CurrencyCode,
			rem.ExchangeRate.StringFixed(4),
			rem.TDSAmount.StringFixed(2),
			rem.Surcharge.StringFixed(2),
			rem.Cess.StringFixed(2),
			rem.TotalTax.StringFixed(2),
		))
	}

	return b.String()
}

func challanSummary27Q(payload *Form27QPayload) []challanLine {
	challanMap := make(map[string]*challanLine)
	for _, rem := range payload.Remittances {
		key := rem.ChallanNumber
		if key == "" {
			key = "PENDING"
		}
		if existing, ok := challanMap[key]; ok {
			amt, _ := decimal.NewFromString(existing.TotalTDS)
			cnt, _ := decimal.NewFromString(existing.DeducteeCount)
			existing.TotalTDS = amt.Add(rem.TotalTax).StringFixed(2)
			existing.DeducteeCount = cnt.Add(decimal.NewFromInt(1)).String()
		} else {
			challanMap[key] = &challanLine{
				ChallanNumber: key,
				BSRCode:       rem.BSRCode,
				DepositDate:   rem.DateOfPayment.Format("02012006"),
				TotalTDS:      rem.TotalTax.StringFixed(2),
				DeducteeCount: "1",
			}
		}
	}
	var lines []challanLine
	for _, v := range challanMap {
		lines = append(lines, *v)
	}
	return lines
}

func lookupDefault(m map[string]string, key, def string) string {
	if m == nil {
		return def
	}
	if v, ok := m[key]; ok {
		return v
	}
	return def
}

func lookupDecimal(m map[string]decimal.Decimal, key string, def decimal.Decimal) decimal.Decimal {
	if m == nil {
		return def
	}
	if v, ok := m[key]; ok {
		return v
	}
	return def
}
