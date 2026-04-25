package provider

import (
	"context"
	"fmt"
	"sync"

	"github.com/complai/complai/services/go/apex-gateway-service/internal/domain"
)

var _ ApexProvider = (*MockProvider)(nil)

// vendorProfile describes the compliance quality tier for mock data generation.
type vendorProfile int

const (
	profileExemplary vendorProfile = iota // Cat A: all compliant
	profileGood                           // Cat B: mostly compliant
	profileAverage                        // Cat C: some issues
	profilePoor                           // Cat D: frequent issues
)

type vendorSeed struct {
	name               string
	legalName          string
	tradeName          string
	pan                string
	stateCode          string
	state              string
	category           string
	msme               bool
	registrationStatus string
	profile            vendorProfile
}

type MockProvider struct {
	mu       sync.RWMutex
	vendors  []domain.Vendor
	invoices []domain.APInvoice
}

func NewMockProvider() *MockProvider {
	p := &MockProvider{}
	p.vendors = p.buildVendors()
	p.invoices = p.buildInvoices()
	return p
}

func (m *MockProvider) FetchVendors(_ context.Context, req *domain.FetchVendorsRequest) (*domain.FetchVendorsResponse, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	vendors := make([]domain.Vendor, len(m.vendors))
	copy(vendors, m.vendors)

	// Set tenant_id on all returned vendors
	for i := range vendors {
		vendors[i].TenantID = req.TenantID
	}

	total := len(vendors)

	// Apply offset
	if req.Offset > 0 {
		if req.Offset >= len(vendors) {
			return &domain.FetchVendorsResponse{
				Vendors:   []domain.Vendor{},
				Total:     total,
				RequestID: req.RequestID,
			}, nil
		}
		vendors = vendors[req.Offset:]
	}

	// Apply limit
	if req.Limit > 0 && req.Limit < len(vendors) {
		vendors = vendors[:req.Limit]
	}

	return &domain.FetchVendorsResponse{
		Vendors:   vendors,
		Total:     total,
		RequestID: req.RequestID,
	}, nil
}

func (m *MockProvider) FetchAPInvoices(_ context.Context, req *domain.FetchAPInvoicesRequest) (*domain.FetchAPInvoicesResponse, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []domain.APInvoice
	for _, inv := range m.invoices {
		inv.TenantID = req.TenantID

		// Filter by vendor_id if specified
		if req.VendorID != "" && inv.VendorID != req.VendorID {
			continue
		}

		// Filter by date range if specified
		if req.FromDate != "" && inv.InvoiceDate < req.FromDate {
			continue
		}
		if req.ToDate != "" && inv.InvoiceDate > req.ToDate {
			continue
		}

		result = append(result, inv)
	}

	if result == nil {
		result = []domain.APInvoice{}
	}

	return &domain.FetchAPInvoicesResponse{
		Invoices:  result,
		Total:     len(result),
		RequestID: req.RequestID,
	}, nil
}

// ---------------------------------------------------------------------------
// Deterministic mock data generation
// ---------------------------------------------------------------------------

func (m *MockProvider) buildVendors() []domain.Vendor {
	seeds := vendorSeeds()
	vendors := make([]domain.Vendor, len(seeds))
	for i, s := range seeds {
		vendors[i] = domain.Vendor{
			ID:                 fmt.Sprintf("VND-%03d", i+1),
			Name:               s.name,
			LegalName:          s.legalName,
			TradeName:          s.tradeName,
			PAN:                s.pan,
			GSTIN:              buildGSTIN(s.stateCode, s.pan, i),
			State:              s.state,
			StateCode:          s.stateCode,
			Category:           s.category,
			RegistrationStatus: s.registrationStatus,
			MSMERegistered:     s.msme,
			Email:              fmt.Sprintf("accounts@%s.example.com", sanitizeEmail(s.tradeName)),
			Phone:              fmt.Sprintf("98%08d", 10000000+i*137),
			Address:            fmt.Sprintf("%d Industrial Area, %s", 100+i*7, s.state),
			CreatedAt:          "2025-01-15T10:00:00Z",
			UpdatedAt:          "2026-03-01T10:00:00Z",
		}
	}
	return vendors
}

func (m *MockProvider) buildInvoices() []domain.APInvoice {
	seeds := vendorSeeds()
	var invoices []domain.APInvoice

	// Invoice months: Oct 2025 to Mar 2026 (6 months)
	months := []struct {
		month string
		year  string
	}{
		{"10", "2025"}, {"11", "2025"}, {"12", "2025"},
		{"01", "2026"}, {"02", "2026"}, {"03", "2026"},
	}

	// GST rate options
	gstRates := []float64{0.05, 0.12, 0.18, 0.28}

	// Deterministic base amounts (in rupees)
	baseAmounts := []float64{
		15000, 42000, 125000, 350000, 780000, 1250000, 2500000, 4800000,
		18500, 67000, 230000, 550000, 95000, 310000, 1100000, 3200000,
	}

	invCounter := 0
	for vendorIdx, s := range seeds {
		vendorID := fmt.Sprintf("VND-%03d", vendorIdx+1)
		gstin := buildGSTIN(s.stateCode, s.pan, vendorIdx)

		// Number of invoices varies by profile
		var numInvoices int
		switch s.profile {
		case profileExemplary:
			numInvoices = 6 + (vendorIdx % 4) // 6-9
		case profileGood:
			numInvoices = 5 + (vendorIdx % 3) // 5-7
		case profileAverage:
			numInvoices = 4 + (vendorIdx % 3) // 4-6
		case profilePoor:
			numInvoices = 3 + (vendorIdx % 3) // 3-5
		}

		for j := 0; j < numInvoices; j++ {
			invCounter++
			monthIdx := j % len(months)
			m := months[monthIdx]
			day := 5 + (invCounter*3)%23 // Days 5-27

			invoiceDate := fmt.Sprintf("%s-%s-%02d", m.year, m.month, day)
			dueDay := day + 30
			dueMonth := monthIdx
			dueYear := m.year
			if dueDay > 28 {
				dueDay = dueDay - 28
				dueMonth++
			}
			if dueMonth >= len(months) {
				dueMonth = len(months) - 1
				dueYear = "2026"
			}
			dm := months[dueMonth]
			dueDate := fmt.Sprintf("%s-%s-%02d", dueYear, dm.month, dueDay)

			// Select base amount deterministically
			baseAmt := baseAmounts[(vendorIdx*7+j*3)%len(baseAmounts)]

			// Select GST rate deterministically
			gstRate := gstRates[(vendorIdx+j)%len(gstRates)]

			// Determine if inter-state (IGST) or intra-state (CGST+SGST)
			isInterState := (vendorIdx+j)%3 == 0
			var cgst, sgst, igst float64
			if isInterState {
				igst = baseAmt * gstRate
			} else {
				cgst = baseAmt * gstRate / 2
				sgst = baseAmt * gstRate / 2
			}
			totalAmt := baseAmt + cgst + sgst + igst

			// Payment status based on profile and invoice index
			paymentStatus := determinePaymentStatus(s.profile, vendorIdx, j)
			paymentDate := ""
			if paymentStatus == "paid" {
				paymentDate = fmt.Sprintf("2026-%02d-15", 1+(j%3))
			}

			// IRN generation based on profile
			irnGenerated := determineIRN(s.profile, vendorIdx, j)
			irn := ""
			if irnGenerated {
				irn = fmt.Sprintf("IRN%012d", invCounter)
			}

			// GST filing status based on profile
			gstFilingStatus := determineGSTFilingStatus(s.profile, vendorIdx, j)

			// Mismatch status based on profile
			mismatchStatus := determineMismatchStatus(s.profile, vendorIdx, j)

			invoices = append(invoices, domain.APInvoice{
				ID:              fmt.Sprintf("INV-%05d", invCounter),
				VendorID:        vendorID,
				VendorGSTIN:     gstin,
				InvoiceNumber:   fmt.Sprintf("%s/INV/%s%s/%03d", s.tradeName[:3], m.year[2:], m.month, j+1),
				InvoiceDate:     invoiceDate,
				DueDate:         dueDate,
				TaxableValue:    baseAmt,
				CGSTAmount:      cgst,
				SGSTAmount:      sgst,
				IGSTAmount:      igst,
				TotalAmount:     totalAmt,
				PaymentStatus:   paymentStatus,
				PaymentDate:     paymentDate,
				IRNGenerated:    irnGenerated,
				IRN:             irn,
				GSTFilingStatus: gstFilingStatus,
				MismatchStatus:  mismatchStatus,
				CreatedAt:       invoiceDate + "T10:00:00Z",
			})
		}
	}

	return invoices
}

func determinePaymentStatus(profile vendorProfile, vendorIdx, invoiceIdx int) string {
	seed := vendorIdx*11 + invoiceIdx*7
	switch profile {
	case profileExemplary:
		// 80% paid, 15% unpaid (current), 5% partial
		v := seed % 20
		if v < 16 {
			return "paid"
		} else if v < 19 {
			return "unpaid"
		}
		return "partial"
	case profileGood:
		// 70% paid, 15% unpaid, 10% overdue, 5% partial
		v := seed % 20
		if v < 14 {
			return "paid"
		} else if v < 17 {
			return "unpaid"
		} else if v < 19 {
			return "overdue"
		}
		return "partial"
	case profileAverage:
		// 50% paid, 20% unpaid, 20% overdue, 10% partial
		v := seed % 10
		if v < 5 {
			return "paid"
		} else if v < 7 {
			return "unpaid"
		} else if v < 9 {
			return "overdue"
		}
		return "partial"
	case profilePoor:
		// 30% paid, 20% unpaid, 40% overdue, 10% partial
		v := seed % 10
		if v < 3 {
			return "paid"
		} else if v < 5 {
			return "unpaid"
		} else if v < 9 {
			return "overdue"
		}
		return "partial"
	}
	return "unpaid"
}

func determineIRN(profile vendorProfile, vendorIdx, invoiceIdx int) bool {
	seed := vendorIdx*13 + invoiceIdx*5
	switch profile {
	case profileExemplary:
		return true // 100% IRN
	case profileGood:
		return seed%10 < 9 // 90%
	case profileAverage:
		return seed%10 < 6 // 60%
	case profilePoor:
		return seed%10 < 2 // 20%
	}
	return false
}

func determineGSTFilingStatus(profile vendorProfile, vendorIdx, invoiceIdx int) string {
	seed := vendorIdx*17 + invoiceIdx*3
	switch profile {
	case profileExemplary:
		return "filed" // All filed on time
	case profileGood:
		// 80% filed, 10% pending, 10% late
		v := seed % 10
		if v < 8 {
			return "filed"
		} else if v < 9 {
			return "pending"
		}
		return "late"
	case profileAverage:
		// 50% filed, 20% pending, 20% late, 10% not_filed
		v := seed % 10
		if v < 5 {
			return "filed"
		} else if v < 7 {
			return "pending"
		} else if v < 9 {
			return "late"
		}
		return "not_filed"
	case profilePoor:
		// 20% filed, 10% pending, 30% late, 40% not_filed
		v := seed % 10
		if v < 2 {
			return "filed"
		} else if v < 3 {
			return "pending"
		} else if v < 6 {
			return "late"
		}
		return "not_filed"
	}
	return "pending"
}

func determineMismatchStatus(profile vendorProfile, vendorIdx, invoiceIdx int) string {
	seed := vendorIdx*19 + invoiceIdx*11
	switch profile {
	case profileExemplary:
		return "matched" // 100% matched
	case profileGood:
		// 95% matched, 5% mismatched
		if seed%20 < 19 {
			return "matched"
		}
		return "mismatched"
	case profileAverage:
		// 75% matched, 15% mismatched, 10% pending
		v := seed % 20
		if v < 15 {
			return "matched"
		} else if v < 18 {
			return "mismatched"
		}
		return "pending"
	case profilePoor:
		// 50% matched, 30% mismatched, 20% pending
		v := seed % 10
		if v < 5 {
			return "matched"
		} else if v < 8 {
			return "mismatched"
		}
		return "pending"
	}
	return "pending"
}

func buildGSTIN(stateCode, pan string, idx int) string {
	// GSTIN format: 2-digit state code + PAN (10 chars) + 1 digit + Z + check char
	entityNum := (idx % 9) + 1
	checkChars := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	check := string(checkChars[idx%26])
	return fmt.Sprintf("%s%s%dZ%s", stateCode, pan, entityNum, check)
}

func sanitizeEmail(name string) string {
	result := make([]byte, 0, len(name))
	for _, c := range name {
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') {
			if c >= 'A' && c <= 'Z' {
				c = c + 32 // to lowercase
			}
			result = append(result, byte(c))
		}
	}
	if len(result) == 0 {
		return "vendor"
	}
	return string(result)
}

// vendorSeeds returns the deterministic list of 50 vendor seed data.
func vendorSeeds() []vendorSeed {
	return []vendorSeed{
		// --- 10 Exemplary (Cat A) vendors ---
		{name: "Tata Steel Ltd", legalName: "Tata Steel Limited", tradeName: "Tata Steel", pan: "AABCT1234A", stateCode: "29", state: "Karnataka", category: "Manufacturer", msme: false, registrationStatus: "Active", profile: profileExemplary},
		{name: "Infosys Ltd", legalName: "Infosys Limited", tradeName: "Infosys", pan: "AABCI5678B", stateCode: "29", state: "Karnataka", category: "Service Provider", msme: false, registrationStatus: "Active", profile: profileExemplary},
		{name: "Reliance Industries Ltd", legalName: "Reliance Industries Limited", tradeName: "Reliance", pan: "AABCR9012C", stateCode: "27", state: "Maharashtra", category: "Manufacturer", msme: false, registrationStatus: "Active", profile: profileExemplary},
		{name: "Wipro Ltd", legalName: "Wipro Limited", tradeName: "Wipro", pan: "AABCW3456D", stateCode: "29", state: "Karnataka", category: "Service Provider", msme: false, registrationStatus: "Active", profile: profileExemplary},
		{name: "Larsen & Toubro Ltd", legalName: "Larsen and Toubro Limited", tradeName: "LnT", pan: "AABCL7890E", stateCode: "27", state: "Maharashtra", category: "Manufacturer", msme: false, registrationStatus: "Active", profile: profileExemplary},
		{name: "Hindustan Unilever Ltd", legalName: "Hindustan Unilever Limited", tradeName: "HUL", pan: "AABCH2345F", stateCode: "27", state: "Maharashtra", category: "Manufacturer", msme: false, registrationStatus: "Active", profile: profileExemplary},
		{name: "Mahindra & Mahindra Ltd", legalName: "Mahindra and Mahindra Limited", tradeName: "Mahindra", pan: "AABCM6789G", stateCode: "27", state: "Maharashtra", category: "Manufacturer", msme: false, registrationStatus: "Active", profile: profileExemplary},
		{name: "TCS Ltd", legalName: "Tata Consultancy Services Limited", tradeName: "TCS", pan: "AABCT1122H", stateCode: "27", state: "Maharashtra", category: "Service Provider", msme: false, registrationStatus: "Active", profile: profileExemplary},
		{name: "Asian Paints Ltd", legalName: "Asian Paints Limited", tradeName: "AsianPaints", pan: "AABCA3344I", stateCode: "27", state: "Maharashtra", category: "Manufacturer", msme: false, registrationStatus: "Active", profile: profileExemplary},
		{name: "Bosch Ltd", legalName: "Bosch Limited", tradeName: "Bosch", pan: "AABCB5566J", stateCode: "29", state: "Karnataka", category: "Manufacturer", msme: false, registrationStatus: "Active", profile: profileExemplary},

		// --- 15 Good (Cat B) vendors ---
		{name: "Godrej Industries Ltd", legalName: "Godrej Industries Limited", tradeName: "Godrej", pan: "BBDCG1234A", stateCode: "27", state: "Maharashtra", category: "Manufacturer", msme: false, registrationStatus: "Active", profile: profileGood},
		{name: "Bajaj Auto Ltd", legalName: "Bajaj Auto Limited", tradeName: "BajajAuto", pan: "BBDCB5678B", stateCode: "27", state: "Maharashtra", category: "Manufacturer", msme: false, registrationStatus: "Active", profile: profileGood},
		{name: "Havells India Ltd", legalName: "Havells India Limited", tradeName: "Havells", pan: "BBDCH9012C", stateCode: "07", state: "Delhi", category: "Manufacturer", msme: false, registrationStatus: "Active", profile: profileGood},
		{name: "Amara Raja Batteries Ltd", legalName: "Amara Raja Batteries Limited", tradeName: "AmaraRaja", pan: "BBDCA3456D", stateCode: "37", state: "Andhra Pradesh", category: "Manufacturer", msme: false, registrationStatus: "Active", profile: profileGood},
		{name: "Voltas Ltd", legalName: "Voltas Limited", tradeName: "Voltas", pan: "BBDCV7890E", stateCode: "27", state: "Maharashtra", category: "Trader", msme: false, registrationStatus: "Active", profile: profileGood},
		{name: "Mphasis Ltd", legalName: "Mphasis Limited", tradeName: "Mphasis", pan: "BBDCM2345F", stateCode: "29", state: "Karnataka", category: "Service Provider", msme: false, registrationStatus: "Active", profile: profileGood},
		{name: "Blue Star Ltd", legalName: "Blue Star Limited", tradeName: "BlueStar", pan: "BBDCB6789G", stateCode: "27", state: "Maharashtra", category: "Manufacturer", msme: true, registrationStatus: "Active", profile: profileGood},
		{name: "Cummins India Ltd", legalName: "Cummins India Limited", tradeName: "Cummins", pan: "BBDCC1122H", stateCode: "27", state: "Maharashtra", category: "Manufacturer", msme: false, registrationStatus: "Active", profile: profileGood},
		{name: "Thermax Ltd", legalName: "Thermax Limited", tradeName: "Thermax", pan: "BBDCT3344I", stateCode: "27", state: "Maharashtra", category: "Manufacturer", msme: true, registrationStatus: "Active", profile: profileGood},
		{name: "Persistent Systems Ltd", legalName: "Persistent Systems Limited", tradeName: "Persistent", pan: "BBDCP5566J", stateCode: "27", state: "Maharashtra", category: "Service Provider", msme: false, registrationStatus: "Active", profile: profileGood},
		{name: "Zensar Technologies Ltd", legalName: "Zensar Technologies Limited", tradeName: "Zensar", pan: "BBDCZ7788K", stateCode: "27", state: "Maharashtra", category: "Service Provider", msme: false, registrationStatus: "Active", profile: profileGood},
		{name: "Kalpataru Projects Ltd", legalName: "Kalpataru Projects International Limited", tradeName: "Kalpataru", pan: "BBDCK9900L", stateCode: "27", state: "Maharashtra", category: "Manufacturer", msme: false, registrationStatus: "Active", profile: profileGood},
		{name: "Galaxy Surfactants Ltd", legalName: "Galaxy Surfactants Limited", tradeName: "Galaxy", pan: "BBDCG1100M", stateCode: "27", state: "Maharashtra", category: "Manufacturer", msme: true, registrationStatus: "Active", profile: profileGood},
		{name: "KPIT Technologies Ltd", legalName: "KPIT Technologies Limited", tradeName: "KPIT", pan: "BBDCK3300N", stateCode: "27", state: "Maharashtra", category: "Service Provider", msme: false, registrationStatus: "Active", profile: profileGood},
		{name: "Deepak Nitrite Ltd", legalName: "Deepak Nitrite Limited", tradeName: "DeepakNitrite", pan: "BBDCD5500P", stateCode: "24", state: "Gujarat", category: "Manufacturer", msme: false, registrationStatus: "Active", profile: profileGood},

		// --- 15 Average (Cat C) vendors ---
		{name: "Sri Balaji Traders", legalName: "Sri Balaji Traders Pvt Ltd", tradeName: "BalajTraders", pan: "CCDCB1234A", stateCode: "33", state: "Tamil Nadu", category: "Trader", msme: true, registrationStatus: "Active", profile: profileAverage},
		{name: "Priya Enterprises LLP", legalName: "Priya Enterprises LLP", tradeName: "PriyaEnt", pan: "CCDCP5678B", stateCode: "29", state: "Karnataka", category: "Trader", msme: true, registrationStatus: "Active", profile: profileAverage},
		{name: "Gupta Logistics Pvt Ltd", legalName: "Gupta Logistics Private Limited", tradeName: "GuptaLog", pan: "CCDCG9012C", stateCode: "09", state: "Uttar Pradesh", category: "Logistics", msme: false, registrationStatus: "Active", profile: profileAverage},
		{name: "Devi Engineering Works", legalName: "Devi Engineering Works Pvt Ltd", tradeName: "DeviEng", pan: "CCDCD3456D", stateCode: "36", state: "Telangana", category: "Manufacturer", msme: true, registrationStatus: "Active", profile: profileAverage},
		{name: "Sharma IT Solutions LLP", legalName: "Sharma IT Solutions LLP", tradeName: "SharmaIT", pan: "CCDCS7890E", stateCode: "07", state: "Delhi", category: "Service Provider", msme: true, registrationStatus: "Active", profile: profileAverage},
		{name: "Kolkata Paper Mills", legalName: "Kolkata Paper Mills Pvt Ltd", tradeName: "KolPaper", pan: "CCDCK2345F", stateCode: "19", state: "West Bengal", category: "Manufacturer", msme: true, registrationStatus: "Active", profile: profileAverage},
		{name: "Rajasthan Marble Co", legalName: "Rajasthan Marble Company Pvt Ltd", tradeName: "RajMarble", pan: "CCDCR6789G", stateCode: "08", state: "Rajasthan", category: "Manufacturer", msme: false, registrationStatus: "Active", profile: profileAverage},
		{name: "Sunrise Imports OPC", legalName: "Sunrise Imports OPC Pvt Ltd", tradeName: "SunriseImp", pan: "CCDCS1122H", stateCode: "27", state: "Maharashtra", category: "Import Supplier", msme: false, registrationStatus: "Active", profile: profileAverage},
		{name: "Patel Chemicals LLP", legalName: "Patel Chemicals LLP", tradeName: "PatelChem", pan: "CCDCP3344I", stateCode: "24", state: "Gujarat", category: "Manufacturer", msme: true, registrationStatus: "Active", profile: profileAverage},
		{name: "Nandi Transport Services", legalName: "Nandi Transport Services Pvt Ltd", tradeName: "NandiTrans", pan: "CCDCN5566J", stateCode: "29", state: "Karnataka", category: "Logistics", msme: false, registrationStatus: "Active", profile: profileAverage},
		{name: "Bharat Textiles Pvt Ltd", legalName: "Bharat Textiles Private Limited", tradeName: "BharatTex", pan: "CCDCB7788K", stateCode: "33", state: "Tamil Nadu", category: "Manufacturer", msme: true, registrationStatus: "Active", profile: profileAverage},
		{name: "Kumar Electricals", legalName: "Kumar Electricals Pvt Ltd", tradeName: "KumarElec", pan: "CCDCK9900L", stateCode: "09", state: "Uttar Pradesh", category: "Trader", msme: true, registrationStatus: "Active", profile: profileAverage},
		{name: "Chennai Auto Parts", legalName: "Chennai Auto Parts Pvt Ltd", tradeName: "ChennaiAuto", pan: "CCDCC1100M", stateCode: "33", state: "Tamil Nadu", category: "Manufacturer", msme: false, registrationStatus: "Active", profile: profileAverage},
		{name: "Hyderabad Packaging LLP", legalName: "Hyderabad Packaging LLP", tradeName: "HydPack", pan: "CCDCH3300N", stateCode: "36", state: "Telangana", category: "Manufacturer", msme: true, registrationStatus: "Active", profile: profileAverage},
		{name: "Jaipur Handicrafts OPC", legalName: "Jaipur Handicrafts OPC Pvt Ltd", tradeName: "JaipurHC", pan: "CCDCJ5500P", stateCode: "08", state: "Rajasthan", category: "Manufacturer", msme: true, registrationStatus: "Active", profile: profileAverage},

		// --- 10 Poor (Cat D) vendors ---
		{name: "Fly-By-Night Traders", legalName: "Fly By Night Traders Pvt Ltd", tradeName: "FlyByNight", pan: "DDDCF1234A", stateCode: "09", state: "Uttar Pradesh", category: "Trader", msme: false, registrationStatus: "Active", profile: profilePoor},
		{name: "Lucky Imports Co", legalName: "Lucky Imports Company", tradeName: "LuckyImp", pan: "DDDCL5678B", stateCode: "07", state: "Delhi", category: "Import Supplier", msme: false, registrationStatus: "Suspended", profile: profilePoor},
		{name: "Dubious Supplies", legalName: "Dubious Supplies Proprietorship", tradeName: "DubiousSup", pan: "DDDCD9012C", stateCode: "09", state: "Uttar Pradesh", category: "Trader", msme: false, registrationStatus: "Active", profile: profilePoor},
		{name: "Ghost Services Pvt Ltd", legalName: "Ghost Services Private Limited", tradeName: "GhostSvc", pan: "DDDCG3456D", stateCode: "19", state: "West Bengal", category: "Service Provider", msme: false, registrationStatus: "Active", profile: profilePoor},
		{name: "Shady Logistics", legalName: "Shady Logistics Proprietorship", tradeName: "ShadyLog", pan: "DDDCS7890E", stateCode: "33", state: "Tamil Nadu", category: "Logistics", msme: false, registrationStatus: "Active", profile: profilePoor},
		{name: "Phantom Metals LLP", legalName: "Phantom Metals LLP", tradeName: "PhantomMet", pan: "DDDCP2345F", stateCode: "24", state: "Gujarat", category: "Manufacturer", msme: false, registrationStatus: "Cancelled", profile: profilePoor},
		{name: "Rogue Electronics", legalName: "Rogue Electronics Proprietorship", tradeName: "RogueElec", pan: "DDDCR6789G", stateCode: "27", state: "Maharashtra", category: "Trader", msme: false, registrationStatus: "Active", profile: profilePoor},
		{name: "Vanishing Textiles OPC", legalName: "Vanishing Textiles OPC Pvt Ltd", tradeName: "VanishTex", pan: "DDDCV1122H", stateCode: "08", state: "Rajasthan", category: "Manufacturer", msme: false, registrationStatus: "Active", profile: profilePoor},
		{name: "Bogus Chemicals", legalName: "Bogus Chemicals Proprietorship", tradeName: "BogusChem", pan: "DDDCB3344I", stateCode: "36", state: "Telangana", category: "Manufacturer", msme: false, registrationStatus: "Suspended", profile: profilePoor},
		{name: "Shadow IT Solutions", legalName: "Shadow IT Solutions Pvt Ltd", tradeName: "ShadowIT", pan: "DDDCS5566J", stateCode: "07", state: "Delhi", category: "Service Provider", msme: false, registrationStatus: "Active", profile: profilePoor},
	}
}
