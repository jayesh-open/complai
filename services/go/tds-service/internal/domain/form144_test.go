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
		{ID: uuid.MustParse("cccc0001-0001-0001-0001-000000000002"), TenantID: testTenant, VendorID: uuid.New(), Name: "CloudServe LLC", PAN: "BBCGS5678H", PANVerified: true, DeducteeType: DeducteeCompany, ResidentStatus: NonResident},
		{ID: uuid.MustParse("cccc0001-0001-0001-0001-000000000003"), TenantID: testTenant, VendorID: uuid.New(), Name: "John Consultant", PAN: "CCDHT9012J", PANVerified: true, DeducteeType: DeducteeIndividual, ResidentStatus: NonResident},
	}
}

func makeSection393_2Entries(deductees []Deductee) []TDSEntry {
	txDate := time.Date(2026, 5, 20, 0, 0, 0, 0, time.UTC)
	return []TDSEntry{
		{
			ID: uuid.New(), TenantID: testTenant, DeducteeID: deductees[0].ID,
			Section: Section393_2, PaymentCode: CodeNonResident, SubClause: "Sl.17",
			FinancialYear: "2026-27", TaxYear: "2026-27", Quarter: "Q1",
			TransactionDate: txDate, GrossAmount: decimal.NewFromInt(5000000),
			TDSRate: decimal.NewFromFloat(0.20), TDSAmount: decimal.NewFromInt(1000000),
			Surcharge: decimal.Zero, Cess: decimal.NewFromInt(40000), TotalTax: decimal.NewFromInt(1040000),
			NatureOfPayment: "Software services", PANAtDeduction: deductees[0].PAN,
			Status: StatusPending, ChallanNumber: "CHN-NR-001", BSRCode: "BSR003",
		},
		{
			ID: uuid.New(), TenantID: testTenant, DeducteeID: deductees[1].ID,
			Section: Section393_2, PaymentCode: CodeNonResident, SubClause: "Sl.17",
			FinancialYear: "2026-27", TaxYear: "2026-27", Quarter: "Q1",
			TransactionDate: txDate, GrossAmount: decimal.NewFromInt(2000000),
			TDSRate: decimal.NewFromFloat(0.15), TDSAmount: decimal.NewFromInt(300000),
			Surcharge: decimal.Zero, Cess: decimal.NewFromInt(12000), TotalTax: decimal.NewFromInt(312000),
			NatureOfPayment: "Cloud hosting", PANAtDeduction: deductees[1].PAN,
			Status: StatusDeposited, ChallanNumber: "CHN-NR-002", BSRCode: "BSR003",
		},
		{
			ID: uuid.New(), TenantID: testTenant, DeducteeID: deductees[2].ID,
			Section: Section393_2, PaymentCode: CodeNonResident, SubClause: "Sl.17",
			FinancialYear: "2026-27", TaxYear: "2026-27", Quarter: "Q1",
			TransactionDate: txDate, GrossAmount: decimal.NewFromInt(1000000),
			TDSRate: decimal.NewFromFloat(0.20), TDSAmount: decimal.NewFromInt(200000),
			Surcharge: decimal.Zero, Cess: decimal.NewFromInt(8000), TotalTax: decimal.NewFromInt(208000),
			NatureOfPayment: "Consulting fees", PANAtDeduction: deductees[2].PAN,
			Status: StatusPending, ChallanNumber: "CHN-NR-003", BSRCode: "BSR003",
		},
	}
}

func makeForm144Input(deductees []Deductee, entries []TDSEntry) Form144Input {
	countryCodes := map[string]string{}
	currencyCodes := map[string]string{}
	exchangeRates := map[string]decimal.Decimal{}
	foreignAmounts := map[string]decimal.Decimal{}
	dtaaArticles := map[string]string{}
	dtaaRates := map[string]decimal.Decimal{}

	countries := []string{"US", "SG", "GB"}
	currencies := []string{"USD", "SGD", "GBP"}
	rates := []float64{83.50, 62.10, 105.20}
	fcAmounts := []float64{59880.24, 32206.12, 9505.70}

	for i, d := range deductees {
		id := d.ID.String()
		countryCodes[id] = countries[i]
		currencyCodes[id] = currencies[i]
		exchangeRates[id] = decimal.NewFromFloat(rates[i])
		foreignAmounts[id] = decimal.NewFromFloat(fcAmounts[i])
		dtaaArticles[id] = "12"
		dtaaRates[id] = decimal.NewFromFloat(0.15)
	}

	return Form144Input{
		Deductor:       testDeductor(),
		FinancialYear:  "2026-27",
		Quarter:        "Q1",
		Deductees:      deductees,
		Entries:        entries,
		CountryCodes:   countryCodes,
		DTAAArticles:   dtaaArticles,
		DTAARates:      dtaaRates,
		CurrencyCodes:  currencyCodes,
		ExchangeRates:  exchangeRates,
		ForeignAmounts: foreignAmounts,
	}
}

func TestGenerateForm144_Success(t *testing.T) {
	deductees := makeNonResidentDeductees()
	entries := makeSection393_2Entries(deductees)

	payload, err := GenerateForm144(makeForm144Input(deductees, entries))
	require.NoError(t, err)
	assert.Equal(t, FormType144, payload.FormType)
	assert.Equal(t, "2026-27", payload.FinancialYear)
	assert.Equal(t, "Q1", payload.Quarter)
	assert.Len(t, payload.Remittances, 3)
	assert.True(t, payload.TotalTDS.Equal(decimal.NewFromInt(1560000)))
	assert.True(t, payload.TotalPaid.Equal(decimal.NewFromInt(8000000)))
	assert.Empty(t, payload.Errors)
}

func TestGenerateForm144_CountryCodesAndDTAA(t *testing.T) {
	deductees := makeNonResidentDeductees()
	entries := makeSection393_2Entries(deductees)

	payload, err := GenerateForm144(makeForm144Input(deductees, entries))
	require.NoError(t, err)

	found := map[string]bool{}
	for _, rem := range payload.Remittances {
		found[rem.CountryCode] = true
		assert.NotEmpty(t, rem.DTAAArticle)
		assert.True(t, rem.DTAARate.IsPositive())
	}
	assert.True(t, found["US"])
	assert.True(t, found["SG"])
	assert.True(t, found["GB"])
}

func TestGenerateForm144_ForeignCurrencyTracking(t *testing.T) {
	deductees := makeNonResidentDeductees()
	entries := makeSection393_2Entries(deductees)

	payload, err := GenerateForm144(makeForm144Input(deductees, entries))
	require.NoError(t, err)

	for _, rem := range payload.Remittances {
		assert.NotEmpty(t, rem.CurrencyCode)
		assert.True(t, rem.ExchangeRate.IsPositive())
		assert.True(t, rem.GrossAmountFC.IsPositive())
	}
}

func TestGenerateForm144_MissingTAN(t *testing.T) {
	_, err := GenerateForm144(Form144Input{
		Deductor:      DeductorDetails{},
		FinancialYear: "2026-27",
		Quarter:       "Q1",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "TAN")
}

func TestGenerateForm144_MissingFYQuarter(t *testing.T) {
	_, err := GenerateForm144(Form144Input{
		Deductor: testDeductor(),
	})
	assert.Error(t, err)
}

func TestGenerateForm144_NoSection393_2Entries(t *testing.T) {
	residentOnly := []TDSEntry{
		{Section: Section393_1, FinancialYear: "2026-27", Quarter: "Q1"},
		{Section: Section392, FinancialYear: "2026-27", Quarter: "Q1"},
	}
	_, err := GenerateForm144(Form144Input{
		Deductor:      testDeductor(),
		FinancialYear: "2026-27",
		Quarter:       "Q1",
		Entries:       residentOnly,
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no non-resident")
}

func TestGenerateForm144_DeducteeNotFound(t *testing.T) {
	orphanEntry := TDSEntry{
		ID: uuid.New(), TenantID: testTenant, DeducteeID: uuid.New(),
		Section: Section393_2, PaymentCode: CodeNonResident,
		FinancialYear: "2026-27", TaxYear: "2026-27", Quarter: "Q1",
		TransactionDate: testDate, GrossAmount: decimal.NewFromInt(1000000),
		TDSRate: decimal.NewFromFloat(0.20), TDSAmount: decimal.NewFromInt(200000),
		Cess: decimal.NewFromInt(8000), TotalTax: decimal.NewFromInt(208000),
		Status: StatusPending,
	}

	payload, err := GenerateForm144(Form144Input{
		Deductor:      testDeductor(),
		FinancialYear: "2026-27",
		Quarter:       "Q1",
		Entries:       []TDSEntry{orphanEntry},
	})
	require.NoError(t, err)
	assert.Len(t, payload.Errors, 1)
	assert.Contains(t, payload.Errors[0], "not found")
}

func TestGenerateForm144_ResidentStatusWarning(t *testing.T) {
	residentDeductee := Deductee{
		ID: uuid.MustParse("cccc0001-0001-0001-0001-000000000099"),
		TenantID: testTenant, VendorID: uuid.New(), Name: "Domestic Corp",
		PAN: "AABCD1234E", DeducteeType: DeducteeCompany, ResidentStatus: Resident,
	}
	entry := TDSEntry{
		ID: uuid.New(), TenantID: testTenant, DeducteeID: residentDeductee.ID,
		Section: Section393_2, PaymentCode: CodeNonResident,
		FinancialYear: "2026-27", TaxYear: "2026-27", Quarter: "Q1",
		TransactionDate: testDate, GrossAmount: decimal.NewFromInt(500000),
		TDSRate: decimal.NewFromFloat(0.20), TDSAmount: decimal.NewFromInt(100000),
		Cess: decimal.NewFromInt(4000), TotalTax: decimal.NewFromInt(104000),
		Status: StatusPending, ChallanNumber: "CHN-NR-X", BSRCode: "BSR003",
	}

	payload, err := GenerateForm144(Form144Input{
		Deductor:      testDeductor(),
		FinancialYear: "2026-27",
		Quarter:       "Q1",
		Deductees:     []Deductee{residentDeductee},
		Entries:       []TDSEntry{entry},
	})
	require.NoError(t, err)
	assert.Len(t, payload.Errors, 1)
	assert.Contains(t, payload.Errors[0], "not non-resident")
}

func TestGenerateForm144FVU_Format(t *testing.T) {
	deductees := makeNonResidentDeductees()
	entries := makeSection393_2Entries(deductees)

	payload, err := GenerateForm144(makeForm144Input(deductees, entries))
	require.NoError(t, err)

	fvu := GenerateForm144FVU(payload)
	assert.Contains(t, fvu, "^FH^144^")
	assert.Contains(t, fvu, "MUMA12345A")
	assert.Contains(t, fvu, "^BH^")
	assert.Contains(t, fvu, "^NR^")
	assert.Contains(t, fvu, "2026-27")
}

func TestLookupDefault(t *testing.T) {
	m := map[string]string{"a": "val-a"}
	assert.Equal(t, "val-a", lookupDefault(m, "a", "def"))
	assert.Equal(t, "def", lookupDefault(m, "missing", "def"))
	assert.Equal(t, "def", lookupDefault(nil, "a", "def"))
}

func TestLookupDecimal(t *testing.T) {
	m := map[string]decimal.Decimal{"a": decimal.NewFromFloat(1.5)}
	assert.True(t, lookupDecimal(m, "a", decimal.Zero).Equal(decimal.NewFromFloat(1.5)))
	assert.True(t, lookupDecimal(m, "missing", decimal.NewFromInt(99)).Equal(decimal.NewFromInt(99)))
	assert.True(t, lookupDecimal(nil, "a", decimal.NewFromInt(42)).Equal(decimal.NewFromInt(42)))
}

func TestValidFormType(t *testing.T) {
	assert.True(t, ValidFormType(FormType138))
	assert.True(t, ValidFormType(FormType140))
	assert.True(t, ValidFormType(FormType144))
	assert.False(t, ValidFormType(FormType("999")))
	assert.False(t, ValidFormType(FormType("")))
}

func TestFilingIdempotencyKey(t *testing.T) {
	key := FilingIdempotencyKey(testTenant, FormType140, "2026-27", "Q1")
	assert.Contains(t, key, "140")
	assert.Contains(t, key, "2026-27")
	assert.Contains(t, key, "Q1")
	assert.Contains(t, key, testTenant.String())
}
