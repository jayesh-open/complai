package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testDeductor() DeductorDetails {
	return DeductorDetails{
		TAN:          "MUMA12345A",
		DeductorName: "Acme Corp Pvt Ltd",
		DeductorPAN:  "AABCA1234F",
		Address:      "100 MG Road",
		City:         "Mumbai",
		State:        "Maharashtra",
		Pincode:      "400001",
		ContactEmail: "finance@acmecorp.in",
		ContactPhone: "9876543210",
	}
}

func makeNonSalaryDeductees() []Deductee {
	return []Deductee{
		{ID: uuid.MustParse("bbbb0001-0001-0001-0001-000000000001"), TenantID: testTenant, VendorID: uuid.New(), Name: "BuildRight Contractors", PAN: "AABCB1234F", PANVerified: true, DeducteeType: DeducteeCompany, ResidentStatus: Resident},
		{ID: uuid.MustParse("bbbb0001-0001-0001-0001-000000000002"), TenantID: testTenant, VendorID: uuid.New(), Name: "OfficeSpace Realty", PAN: "CCDEF5678G", PANVerified: true, DeducteeType: DeducteeCompany, ResidentStatus: Resident},
		{ID: uuid.MustParse("bbbb0001-0001-0001-0001-000000000003"), TenantID: testTenant, VendorID: uuid.New(), Name: "Ravi Auditor", PAN: "EEFGH9012H", PANVerified: true, DeducteeType: DeducteeIndividual, ResidentStatus: Resident},
		{ID: uuid.MustParse("bbbb0001-0001-0001-0001-000000000004"), TenantID: testTenant, VendorID: uuid.New(), Name: "Supreme Supplies", PAN: "GGHIJ3456K", PANVerified: true, DeducteeType: DeducteeCompany, ResidentStatus: Resident},
		{ID: uuid.MustParse("bbbb0001-0001-0001-0001-000000000005"), TenantID: testTenant, VendorID: uuid.New(), Name: "Dev Plumbing", PAN: "", PANVerified: false, DeducteeType: DeducteeIndividual, ResidentStatus: Resident},
	}
}

func makeNonSalaryEntries(deductees []Deductee) []TDSEntry {
	txDate := time.Date(2026, 5, 20, 0, 0, 0, 0, time.UTC)
	return []TDSEntry{
		{
			ID: uuid.New(), TenantID: testTenant, DeducteeID: deductees[0].ID,
			Section: Section393_1, PaymentCode: CodeContractorOther, SubClause: "Sl.6(i).D(b)",
			FinancialYear: "2026-27", TaxYear: "2026-27", Quarter: "Q1",
			TransactionDate: txDate, GrossAmount: decimal.NewFromInt(500000),
			TDSRate: decimal.NewFromFloat(0.02), TDSAmount: decimal.NewFromInt(10000),
			Surcharge: decimal.Zero, Cess: decimal.Zero, TotalTax: decimal.NewFromInt(10000),
			NatureOfPayment: "Civil construction", PANAtDeduction: deductees[0].PAN,
			Status: StatusPending, ChallanNumber: "CHN-001", BSRCode: "BSR001",
		},
		{
			ID: uuid.New(), TenantID: testTenant, DeducteeID: deductees[1].ID,
			Section: Section393_1, PaymentCode: CodeRentLand, SubClause: "Sl.2(ii).D(b)",
			FinancialYear: "2026-27", TaxYear: "2026-27", Quarter: "Q1",
			TransactionDate: txDate, GrossAmount: decimal.NewFromInt(300000),
			TDSRate: decimal.NewFromFloat(0.10), TDSAmount: decimal.NewFromInt(30000),
			Surcharge: decimal.Zero, Cess: decimal.Zero, TotalTax: decimal.NewFromInt(30000),
			NatureOfPayment: "Office rent", PANAtDeduction: deductees[1].PAN,
			Status: StatusDeposited, ChallanNumber: "CHN-002", BSRCode: "BSR001",
		},
		{
			ID: uuid.New(), TenantID: testTenant, DeducteeID: deductees[2].ID,
			Section: Section393_1, PaymentCode: CodeProfessional, SubClause: "Sl.6(iii).D(b)",
			FinancialYear: "2026-27", TaxYear: "2026-27", Quarter: "Q1",
			TransactionDate: txDate, GrossAmount: decimal.NewFromInt(100000),
			TDSRate: decimal.NewFromFloat(0.10), TDSAmount: decimal.NewFromInt(10000),
			Surcharge: decimal.Zero, Cess: decimal.Zero, TotalTax: decimal.NewFromInt(10000),
			NatureOfPayment: "Audit fees", PANAtDeduction: deductees[2].PAN,
			Status: StatusPending, ChallanNumber: "CHN-003", BSRCode: "BSR002",
		},
		{
			ID: uuid.New(), TenantID: testTenant, DeducteeID: deductees[3].ID,
			Section: Section393_1, PaymentCode: CodePurchaseGoods, SubClause: "Sl.8(ii)",
			FinancialYear: "2026-27", TaxYear: "2026-27", Quarter: "Q1",
			TransactionDate: txDate, GrossAmount: decimal.NewFromInt(6000000),
			TDSRate: decimal.NewFromFloat(0.001), TDSAmount: decimal.NewFromInt(6000),
			Surcharge: decimal.Zero, Cess: decimal.Zero, TotalTax: decimal.NewFromInt(6000),
			NatureOfPayment: "Purchase of goods", PANAtDeduction: deductees[3].PAN,
			Status: StatusPending, ChallanNumber: "CHN-004", BSRCode: "BSR001",
		},
		{
			ID: uuid.New(), TenantID: testTenant, DeducteeID: deductees[4].ID,
			Section: Section393_1, PaymentCode: CodeContractorIndiv, SubClause: "Sl.6(i).D(a)",
			FinancialYear: "2026-27", TaxYear: "2026-27", Quarter: "Q1",
			TransactionDate: txDate, GrossAmount: decimal.NewFromInt(200000),
			TDSRate: decimal.NewFromFloat(0.20), TDSAmount: decimal.NewFromInt(40000),
			Surcharge: decimal.Zero, Cess: decimal.Zero, TotalTax: decimal.NewFromInt(40000),
			NatureOfPayment: "Plumbing work", NoPANDeduction: true,
			Status: StatusPending, ChallanNumber: "CHN-005", BSRCode: "BSR001",
		},
	}
}

func TestGenerateForm140_Success(t *testing.T) {
	deductees := makeNonSalaryDeductees()
	entries := makeNonSalaryEntries(deductees)

	payload, err := GenerateForm140(Form140Input{
		Deductor:      testDeductor(),
		FinancialYear: "2026-27",
		Quarter:       "Q1",
		Deductees:     deductees,
		Entries:       entries,
	})

	require.NoError(t, err)
	assert.Equal(t, FormType140, payload.FormType)
	assert.Len(t, payload.Deductions, 5)
	assert.True(t, payload.TotalTDS.Equal(decimal.NewFromInt(96000)))
	assert.True(t, payload.TotalPaid.Equal(decimal.NewFromInt(7100000)))
}

func TestGenerateForm140_MultiplePaymentCodes(t *testing.T) {
	deductees := makeNonSalaryDeductees()
	entries := makeNonSalaryEntries(deductees)

	payload, err := GenerateForm140(Form140Input{
		Deductor:      testDeductor(),
		FinancialYear: "2026-27",
		Quarter:       "Q1",
		Deductees:     deductees,
		Entries:       entries,
	})

	require.NoError(t, err)
	codeMap := make(map[PaymentCode]int)
	for _, d := range payload.Deductions {
		codeMap[d.PaymentCode]++
	}
	assert.Equal(t, 1, codeMap[CodeContractorOther], "one contractor-other entry")
	assert.Equal(t, 1, codeMap[CodeContractorIndiv], "one contractor-individual entry")
	assert.Equal(t, 1, codeMap[CodeRentLand], "one rent entry")
	assert.Equal(t, 1, codeMap[CodeProfessional], "one professional entry")
	assert.Equal(t, 1, codeMap[CodePurchaseGoods], "one purchase entry")
}

func TestGenerateForm140_NoPANFlagging(t *testing.T) {
	deductees := makeNonSalaryDeductees()
	entries := makeNonSalaryEntries(deductees)

	payload, err := GenerateForm140(Form140Input{
		Deductor:      testDeductor(),
		FinancialYear: "2026-27",
		Quarter:       "Q1",
		Deductees:     deductees,
		Entries:       entries,
	})

	require.NoError(t, err)
	var noPANCount int
	for _, d := range payload.Deductions {
		if d.NoPAN {
			noPANCount++
		}
	}
	assert.Equal(t, 1, noPANCount)
}

func TestGenerateForm140_FiltersSalaryEntries(t *testing.T) {
	deductees := makeNonSalaryDeductees()
	mixed := append(makeNonSalaryEntries(deductees), TDSEntry{
		ID: uuid.New(), TenantID: testTenant, DeducteeID: deductees[0].ID,
		Section: Section392, PaymentCode: CodeSalaryPrivate,
		FinancialYear: "2026-27", TaxYear: "2026-27", Quarter: "Q1",
		TransactionDate: testDate, GrossAmount: decimal.NewFromInt(200000),
		TDSRate: decimal.NewFromFloat(0.10), TDSAmount: decimal.NewFromInt(20000),
		TotalTax: decimal.NewFromInt(20000), Status: StatusPending,
	})

	payload, err := GenerateForm140(Form140Input{
		Deductor:      testDeductor(),
		FinancialYear: "2026-27",
		Quarter:       "Q1",
		Deductees:     deductees,
		Entries:       mixed,
	})

	require.NoError(t, err)
	assert.Len(t, payload.Deductions, 5, "should exclude Section 392 entry")
}

func TestGenerateForm140_MissingTAN(t *testing.T) {
	_, err := GenerateForm140(Form140Input{
		Deductor:      DeductorDetails{},
		FinancialYear: "2026-27",
		Quarter:       "Q1",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "TAN")
}

func TestGenerateForm140_MissingFYQuarter(t *testing.T) {
	_, err := GenerateForm140(Form140Input{
		Deductor: testDeductor(),
	})
	assert.Error(t, err)
}

func TestGenerateForm140_NoNonSalaryEntries(t *testing.T) {
	salaryOnly := []TDSEntry{
		{Section: Section392, FinancialYear: "2026-27", Quarter: "Q1"},
		{Section: Section393_2, FinancialYear: "2026-27", Quarter: "Q1"},
	}
	_, err := GenerateForm140(Form140Input{
		Deductor:      testDeductor(),
		FinancialYear: "2026-27",
		Quarter:       "Q1",
		Entries:       salaryOnly,
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no non-salary")
}

func TestGenerateForm140_DeducteeNotFound(t *testing.T) {
	orphanEntry := TDSEntry{
		ID: uuid.New(), TenantID: testTenant, DeducteeID: uuid.New(),
		Section: Section393_1, PaymentCode: CodeContractorOther,
		FinancialYear: "2026-27", TaxYear: "2026-27", Quarter: "Q1",
		TransactionDate: testDate, GrossAmount: decimal.NewFromInt(100000),
		TDSRate: decimal.NewFromFloat(0.02), TDSAmount: decimal.NewFromInt(2000),
		TotalTax: decimal.NewFromInt(2000), Status: StatusPending,
	}

	payload, err := GenerateForm140(Form140Input{
		Deductor:      testDeductor(),
		FinancialYear: "2026-27",
		Quarter:       "Q1",
		Entries:       []TDSEntry{orphanEntry},
	})
	require.NoError(t, err)
	assert.Len(t, payload.Errors, 1)
}

func TestGenerateForm140FVU_Format(t *testing.T) {
	deductees := makeNonSalaryDeductees()
	entries := makeNonSalaryEntries(deductees)

	payload, err := GenerateForm140(Form140Input{
		Deductor:      testDeductor(),
		FinancialYear: "2026-27",
		Quarter:       "Q1",
		Deductees:     deductees[:2],
		Entries:       entries[:2],
	})
	require.NoError(t, err)

	fvu := GenerateForm140FVU(payload)
	assert.Contains(t, fvu, "^FH^140^")
	assert.Contains(t, fvu, "MUMA12345A")
	assert.Contains(t, fvu, "^DD^")
	assert.Contains(t, fvu, "^BH^")
}

func TestGenerateForm140_MissingPANValidation(t *testing.T) {
	d := Deductee{
		ID: uuid.MustParse("bbbb0001-0001-0001-0001-000000000006"),
		TenantID: testTenant, VendorID: uuid.New(), Name: "Missing PAN Vendor",
		PAN: "", PANVerified: false, DeducteeType: DeducteeCompany, ResidentStatus: Resident,
	}
	entry := TDSEntry{
		ID: uuid.New(), TenantID: testTenant, DeducteeID: d.ID,
		Section: Section393_1, PaymentCode: CodeContractorOther, SubClause: "Sl.6(i).D(b)",
		FinancialYear: "2026-27", TaxYear: "2026-27", Quarter: "Q1",
		TransactionDate: testDate, GrossAmount: decimal.NewFromInt(100000),
		TDSRate: decimal.NewFromFloat(0.02), TDSAmount: decimal.NewFromInt(2000),
		TotalTax: decimal.NewFromInt(2000), Status: StatusPending,
		NoPANDeduction: false,
	}

	payload, err := GenerateForm140(Form140Input{
		Deductor:      testDeductor(),
		FinancialYear: "2026-27",
		Quarter:       "Q1",
		Deductees:     []Deductee{d},
		Entries:       []TDSEntry{entry},
	})
	require.NoError(t, err)
	assert.Len(t, payload.Errors, 1)
	assert.Contains(t, payload.Errors[0], "missing PAN")
}
