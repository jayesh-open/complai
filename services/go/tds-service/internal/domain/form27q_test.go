package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeNonResidentDeductees() []Deductee {
	return []Deductee{
		{ID: uuid.MustParse("cccc0001-0001-0001-0001-000000000001"), TenantID: testTenant, VendorID: uuid.New(), Name: "TechGlobal Inc", PAN: "AABFT1234G", PANVerified: true, DeducteeType: DeducteeCompany, ResidentStatus: NonResident},
		{ID: uuid.MustParse("cccc0001-0001-0001-0001-000000000002"), TenantID: testTenant, VendorID: uuid.New(), Name: "John Smith Consulting", PAN: "BCAJS5678H", PANVerified: true, DeducteeType: DeducteeIndividual, ResidentStatus: NonResident},
		{ID: uuid.MustParse("cccc0001-0001-0001-0001-000000000003"), TenantID: testTenant, VendorID: uuid.New(), Name: "EuroDesign GmbH", PAN: "CDBED9012J", PANVerified: true, DeducteeType: DeducteeCompany, ResidentStatus: NonResident},
		{ID: uuid.MustParse("cccc0001-0001-0001-0001-000000000004"), TenantID: testTenant, VendorID: uuid.New(), Name: "Singapore Consulting Pte", PAN: "DEFSCP3456K", PANVerified: true, DeducteeType: DeducteeCompany, ResidentStatus: NonResident},
		{ID: uuid.MustParse("cccc0001-0001-0001-0001-000000000005"), TenantID: testTenant, VendorID: uuid.New(), Name: "Dubai Trading LLC", PAN: "EFGDT7890L", PANVerified: true, DeducteeType: DeducteeCompany, ResidentStatus: NonResident},
	}
}

func makeSection195Entries(deductees []Deductee) []TDSEntry {
	txDate := time.Date(2025, 5, 15, 0, 0, 0, 0, time.UTC)
	amounts := []int64{5000000, 2000000, 3000000, 1500000, 4000000}
	tdsAmounts := []int64{1000000, 200000, 600000, 150000, 800000}

	var entries []TDSEntry
	for i, d := range deductees {
		entries = append(entries, TDSEntry{
			ID: uuid.New(), TenantID: testTenant, DeducteeID: d.ID,
			Section: Section195, FinancialYear: "2025-26", Quarter: "Q1",
			TransactionDate: txDate, GrossAmount: decimal.NewFromInt(amounts[i]),
			TDSRate: decimal.NewFromFloat(0.20), TDSAmount: decimal.NewFromInt(tdsAmounts[i]),
			Surcharge: decimal.Zero, Cess: decimal.NewFromInt(tdsAmounts[i]).Mul(decimal.NewFromFloat(0.04)).Round(0),
			TotalTax: decimal.NewFromInt(tdsAmounts[i]).Add(decimal.NewFromInt(tdsAmounts[i]).Mul(decimal.NewFromFloat(0.04)).Round(0)),
			NatureOfPayment: "Software services", PANAtDeduction: d.PAN,
			Status: StatusPending, ChallanNumber: "CHN-NR-001", BSRCode: "BSR003",
		})
	}
	return entries
}

func makeForm27QInput(deductees []Deductee) Form27QInput {
	countryCodes := make(map[string]string)
	dtaaArticles := make(map[string]string)
	dtaaRates := make(map[string]decimal.Decimal)
	currencyCodes := make(map[string]string)
	exchangeRates := make(map[string]decimal.Decimal)
	foreignAmounts := make(map[string]decimal.Decimal)

	countries := []string{"US", "US", "DE", "SG", "AE"}
	currencies := []string{"USD", "USD", "EUR", "SGD", "AED"}
	rates := []float64{83.50, 83.50, 91.20, 62.30, 22.75}
	articles := []string{"12", "15", "12", "", ""}
	dtaaRateValues := []float64{0.15, 0.10, 0.10, 0, 0}
	fcAmounts := []float64{59880.24, 23952.10, 32894.74, 24077.85, 175824.18}

	for i, d := range deductees {
		key := d.ID.String()
		countryCodes[key] = countries[i]
		currencyCodes[key] = currencies[i]
		exchangeRates[key] = decimal.NewFromFloat(rates[i])
		foreignAmounts[key] = decimal.NewFromFloat(fcAmounts[i])
		if articles[i] != "" {
			dtaaArticles[key] = articles[i]
			dtaaRates[key] = decimal.NewFromFloat(dtaaRateValues[i])
		}
	}

	return Form27QInput{
		Deductor:       testDeductor(),
		FinancialYear:  "2025-26",
		Quarter:        "Q1",
		Deductees:      deductees,
		Entries:        makeSection195Entries(deductees),
		CountryCodes:   countryCodes,
		DTAAArticles:   dtaaArticles,
		DTAARates:      dtaaRates,
		CurrencyCodes:  currencyCodes,
		ExchangeRates:  exchangeRates,
		ForeignAmounts: foreignAmounts,
	}
}

func TestGenerateForm27Q_Success(t *testing.T) {
	deductees := makeNonResidentDeductees()
	input := makeForm27QInput(deductees)

	payload, err := GenerateForm27Q(input)

	require.NoError(t, err)
	assert.Equal(t, FormType27Q, payload.FormType)
	assert.Len(t, payload.Remittances, 5)
	assert.True(t, payload.TotalTDS.IsPositive())
	assert.True(t, payload.TotalPaid.IsPositive())
	assert.Empty(t, payload.Errors)
}

func TestGenerateForm27Q_CountryCodesAndDTAA(t *testing.T) {
	deductees := makeNonResidentDeductees()
	input := makeForm27QInput(deductees)

	payload, err := GenerateForm27Q(input)
	require.NoError(t, err)

	usCount := 0
	dtaaCount := 0
	for _, r := range payload.Remittances {
		if r.CountryCode == "US" {
			usCount++
		}
		if r.DTAAArticle != "" {
			dtaaCount++
		}
	}
	assert.Equal(t, 2, usCount, "two US deductees")
	assert.Equal(t, 3, dtaaCount, "three entries with DTAA articles")
}

func TestGenerateForm27Q_ForeignCurrencyTracking(t *testing.T) {
	deductees := makeNonResidentDeductees()
	input := makeForm27QInput(deductees)

	payload, err := GenerateForm27Q(input)
	require.NoError(t, err)

	for _, r := range payload.Remittances {
		assert.NotEmpty(t, r.CurrencyCode)
		assert.True(t, r.ExchangeRate.IsPositive())
		assert.True(t, r.GrossAmountFC.IsPositive())
	}
}

func TestGenerateForm27Q_MissingTAN(t *testing.T) {
	_, err := GenerateForm27Q(Form27QInput{
		Deductor:      DeductorDetails{},
		FinancialYear: "2025-26",
		Quarter:       "Q1",
	})
	assert.Error(t, err)
}

func TestGenerateForm27Q_MissingFYQuarter(t *testing.T) {
	_, err := GenerateForm27Q(Form27QInput{
		Deductor: testDeductor(),
	})
	assert.Error(t, err)
}

func TestGenerateForm27Q_NoSection195Entries(t *testing.T) {
	nonS195 := []TDSEntry{
		{Section: Section194C, FinancialYear: "2025-26", Quarter: "Q1"},
	}
	_, err := GenerateForm27Q(Form27QInput{
		Deductor:      testDeductor(),
		FinancialYear: "2025-26",
		Quarter:       "Q1",
		Entries:       nonS195,
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no section 195")
}

func TestGenerateForm27Q_DeducteeNotFound(t *testing.T) {
	entry := TDSEntry{
		ID: uuid.New(), TenantID: testTenant, DeducteeID: uuid.New(),
		Section: Section195, FinancialYear: "2025-26", Quarter: "Q1",
		TransactionDate: testDate, GrossAmount: decimal.NewFromInt(1000000),
		TDSRate: decimal.NewFromFloat(0.20), TDSAmount: decimal.NewFromInt(200000),
		TotalTax: decimal.NewFromInt(208000), Status: StatusPending,
	}

	payload, err := GenerateForm27Q(Form27QInput{
		Deductor:      testDeductor(),
		FinancialYear: "2025-26",
		Quarter:       "Q1",
		Entries:       []TDSEntry{entry},
	})
	require.NoError(t, err)
	assert.Len(t, payload.Errors, 1)
}

func TestGenerateForm27Q_ResidentStatusWarning(t *testing.T) {
	resident := Deductee{
		ID: uuid.MustParse("cccc0001-0001-0001-0001-000000000006"),
		TenantID: testTenant, VendorID: uuid.New(), Name: "Local Firm",
		PAN: "ABCDE1234F", DeducteeType: DeducteeCompany, ResidentStatus: Resident,
	}
	entry := TDSEntry{
		ID: uuid.New(), TenantID: testTenant, DeducteeID: resident.ID,
		Section: Section195, FinancialYear: "2025-26", Quarter: "Q1",
		TransactionDate: testDate, GrossAmount: decimal.NewFromInt(500000),
		TDSRate: decimal.NewFromFloat(0.20), TDSAmount: decimal.NewFromInt(100000),
		TotalTax: decimal.NewFromInt(104000), Status: StatusPending,
	}

	payload, err := GenerateForm27Q(Form27QInput{
		Deductor:      testDeductor(),
		FinancialYear: "2025-26",
		Quarter:       "Q1",
		Deductees:     []Deductee{resident},
		Entries:       []TDSEntry{entry},
	})
	require.NoError(t, err)
	assert.Len(t, payload.Errors, 1)
	assert.Contains(t, payload.Errors[0], "not non-resident")
}

func TestGenerateForm27QFVU_Format(t *testing.T) {
	deductees := makeNonResidentDeductees()[:2]
	input := makeForm27QInput(deductees)
	input.Deductees = deductees
	input.Entries = makeSection195Entries(deductees)

	payload, err := GenerateForm27Q(input)
	require.NoError(t, err)

	fvu := GenerateForm27QFVU(payload)
	assert.Contains(t, fvu, "^FH^27Q^")
	assert.Contains(t, fvu, "MUMA12345A")
	assert.Contains(t, fvu, "^NR^")
	assert.Contains(t, fvu, "^BH^")
	assert.Contains(t, fvu, "US")
}

func TestLookupDefault(t *testing.T) {
	m := map[string]string{"a": "1"}
	assert.Equal(t, "1", lookupDefault(m, "a", "def"))
	assert.Equal(t, "def", lookupDefault(m, "b", "def"))
	assert.Equal(t, "def", lookupDefault(nil, "a", "def"))
}

func TestLookupDecimal(t *testing.T) {
	m := map[string]decimal.Decimal{"a": decimal.NewFromInt(42)}
	assert.True(t, lookupDecimal(m, "a", decimal.Zero).Equal(decimal.NewFromInt(42)))
	assert.True(t, lookupDecimal(m, "b", decimal.NewFromInt(99)).Equal(decimal.NewFromInt(99)))
	assert.True(t, lookupDecimal(nil, "a", decimal.NewFromInt(99)).Equal(decimal.NewFromInt(99)))
}

func TestValidFormType(t *testing.T) {
	assert.True(t, ValidFormType(FormType24Q))
	assert.True(t, ValidFormType(FormType26Q))
	assert.True(t, ValidFormType(FormType27Q))
	assert.False(t, ValidFormType(FormType("28Q")))
}

func TestFilingIdempotencyKey(t *testing.T) {
	key := FilingIdempotencyKey(testTenant, FormType26Q, "2025-26", "Q1")
	assert.Contains(t, key, "26Q")
	assert.Contains(t, key, "2025-26")
	assert.Contains(t, key, "Q1")
	assert.Contains(t, key, testTenant.String())
}
