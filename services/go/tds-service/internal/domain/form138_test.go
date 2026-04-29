package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testTenant = uuid.MustParse("11111111-1111-1111-1111-111111111111")
	testDate   = time.Date(2026, 6, 15, 0, 0, 0, 0, time.UTC)
)

func testEmployer() EmployerDetails {
	return EmployerDetails{
		TAN:          "MUMA12345A",
		EmployerName: "Acme Corp Pvt Ltd",
		EmployerPAN:  "AABCA1234F",
		Address:      "100 MG Road",
		City:         "Mumbai",
		State:        "Maharashtra",
		Pincode:      "400001",
		ContactEmail: "hr@acmecorp.in",
		ContactPhone: "9876543210",
	}
}

func makeSalaryDeductees() []Deductee {
	return []Deductee{
		{ID: uuid.MustParse("aaaa0001-0001-0001-0001-000000000001"), TenantID: testTenant, VendorID: uuid.New(), Name: "Rajesh Kumar", PAN: "ABCPK1234E", PANVerified: true, DeducteeType: DeducteeIndividual, ResidentStatus: Resident},
		{ID: uuid.MustParse("aaaa0001-0001-0001-0001-000000000002"), TenantID: testTenant, VendorID: uuid.New(), Name: "Priya Sharma", PAN: "DEFPS5678F", PANVerified: true, DeducteeType: DeducteeIndividual, ResidentStatus: Resident},
		{ID: uuid.MustParse("aaaa0001-0001-0001-0001-000000000003"), TenantID: testTenant, VendorID: uuid.New(), Name: "Amit Patel", PAN: "GHIAP9012G", PANVerified: true, DeducteeType: DeducteeIndividual, ResidentStatus: Resident},
		{ID: uuid.MustParse("aaaa0001-0001-0001-0001-000000000004"), TenantID: testTenant, VendorID: uuid.New(), Name: "Sunita Rao", PAN: "JKLSR3456H", PANVerified: true, DeducteeType: DeducteeIndividual, ResidentStatus: Resident},
		{ID: uuid.MustParse("aaaa0001-0001-0001-0001-000000000005"), TenantID: testTenant, VendorID: uuid.New(), Name: "Vikram Singh", PAN: "MNOVS7890J", PANVerified: true, DeducteeType: DeducteeIndividual, ResidentStatus: Resident},
	}
}

func makeSalaryEntries(deductees []Deductee) []TDSEntry {
	salaries := []int64{2000000, 1500000, 1200000, 800000, 3000000}
	monthlyTDS := []int64{16033, 8775, 0, 0, 39650}

	var entries []TDSEntry
	for i, d := range deductees {
		monthlySalary := decimal.NewFromInt(salaries[i]).Div(decimal.NewFromInt(12)).Round(0)
		tds := decimal.NewFromInt(monthlyTDS[i])
		for month := 0; month < 3; month++ {
			date := time.Date(2026, time.Month(4+month), 28, 0, 0, 0, 0, time.UTC)
			entries = append(entries, TDSEntry{
				ID: uuid.New(), TenantID: testTenant, DeducteeID: d.ID,
				Section: Section392, PaymentCode: CodeSalaryPrivate,
				FinancialYear: "2026-27", TaxYear: "2026-27", Quarter: "Q1",
				TransactionDate: date, GrossAmount: monthlySalary,
				TDSRate: tds.Div(monthlySalary).Round(4), TDSAmount: tds,
				Surcharge: decimal.Zero, Cess: decimal.Zero, TotalTax: tds,
				NatureOfPayment: "Salary", PANAtDeduction: d.PAN,
				Status:        StatusPending,
				ChallanNumber: "CHN-SAL-001", BSRCode: "BSR001",
			})
		}
	}
	return entries
}

func TestGenerateForm138_Success(t *testing.T) {
	deductees := makeSalaryDeductees()
	entries := makeSalaryEntries(deductees)

	payload, err := GenerateForm138(Form138Input{
		Employer:      testEmployer(),
		FinancialYear: "2026-27",
		Quarter:       "Q1",
		Deductees:     deductees,
		Entries:       entries,
	})

	require.NoError(t, err)
	assert.Equal(t, FormType138, payload.FormType)
	assert.Equal(t, "2026-27", payload.FinancialYear)
	assert.Equal(t, "Q1", payload.Quarter)
	assert.Len(t, payload.Employees, 5)
	assert.True(t, payload.TotalSalary.IsPositive())
	assert.Empty(t, payload.Errors)

	found := map[string]bool{}
	for _, emp := range payload.Employees {
		found[emp.Name] = true
		assert.NotEmpty(t, emp.PAN)
		assert.True(t, emp.GrossSalary.IsPositive())
	}
	assert.True(t, found["Rajesh Kumar"])
	assert.True(t, found["Vikram Singh"])
}

func TestGenerateForm138_AggregatesPerDeductee(t *testing.T) {
	deductees := makeSalaryDeductees()[:1]
	entries := makeSalaryEntries(deductees)

	payload, err := GenerateForm138(Form138Input{
		Employer:      testEmployer(),
		FinancialYear: "2026-27",
		Quarter:       "Q1",
		Deductees:     deductees,
		Entries:       entries,
	})

	require.NoError(t, err)
	assert.Len(t, payload.Employees, 1)
	emp := payload.Employees[0]
	assert.Equal(t, "Rajesh Kumar", emp.Name)
	assert.True(t, emp.TDSDeducted.Equal(decimal.NewFromInt(16033*3)), "should aggregate 3 months of TDS")
}

func TestGenerateForm138_MissingTAN(t *testing.T) {
	_, err := GenerateForm138(Form138Input{
		Employer:      EmployerDetails{},
		FinancialYear: "2026-27",
		Quarter:       "Q1",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "TAN")
}

func TestGenerateForm138_MissingFYQuarter(t *testing.T) {
	_, err := GenerateForm138(Form138Input{
		Employer: testEmployer(),
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "financial_year")
}

func TestGenerateForm138_NoSalaryEntries(t *testing.T) {
	deductees := makeSalaryDeductees()
	entries := []TDSEntry{
		{
			ID: uuid.New(), TenantID: testTenant, DeducteeID: deductees[0].ID,
			Section: Section393_1, PaymentCode: CodeContractorOther,
			FinancialYear: "2026-27", TaxYear: "2026-27", Quarter: "Q1",
			TransactionDate: testDate, GrossAmount: decimal.NewFromInt(100000),
			TDSRate: decimal.NewFromFloat(0.02), TDSAmount: decimal.NewFromInt(2000),
			TotalTax: decimal.NewFromInt(2000), Status: StatusPending,
		},
	}

	_, err := GenerateForm138(Form138Input{
		Employer:      testEmployer(),
		FinancialYear: "2026-27",
		Quarter:       "Q1",
		Deductees:     deductees,
		Entries:       entries,
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no salary entries")
}

func TestGenerateForm138_DeducteeNotFound(t *testing.T) {
	orphanEntry := TDSEntry{
		ID: uuid.New(), TenantID: testTenant, DeducteeID: uuid.New(),
		Section: Section392, PaymentCode: CodeSalaryPrivate,
		FinancialYear: "2026-27", TaxYear: "2026-27", Quarter: "Q1",
		TransactionDate: testDate, GrossAmount: decimal.NewFromInt(150000),
		TDSRate: decimal.NewFromFloat(0.10), TDSAmount: decimal.NewFromInt(15000),
		TotalTax: decimal.NewFromInt(15000), Status: StatusPending,
	}

	payload, err := GenerateForm138(Form138Input{
		Employer:      testEmployer(),
		FinancialYear: "2026-27",
		Quarter:       "Q1",
		Deductees:     []Deductee{},
		Entries:       []TDSEntry{orphanEntry},
	})
	require.NoError(t, err)
	assert.Len(t, payload.Errors, 1)
	assert.Contains(t, payload.Errors[0], "not found")
}

func TestGenerateForm138_MissingPAN(t *testing.T) {
	d := Deductee{
		ID: uuid.MustParse("aaaa0001-0001-0001-0001-000000000001"),
		TenantID: testTenant, VendorID: uuid.New(), Name: "No PAN Person",
		PAN: "", DeducteeType: DeducteeIndividual, ResidentStatus: Resident,
	}
	entry := TDSEntry{
		ID: uuid.New(), TenantID: testTenant, DeducteeID: d.ID,
		Section: Section392, PaymentCode: CodeSalaryPrivate,
		FinancialYear: "2026-27", TaxYear: "2026-27", Quarter: "Q1",
		TransactionDate: testDate, GrossAmount: decimal.NewFromInt(100000),
		TDSRate: decimal.NewFromFloat(0.20), TDSAmount: decimal.NewFromInt(20000),
		TotalTax: decimal.NewFromInt(20000), Status: StatusPending,
	}

	payload, err := GenerateForm138(Form138Input{
		Employer:      testEmployer(),
		FinancialYear: "2026-27",
		Quarter:       "Q1",
		Deductees:     []Deductee{d},
		Entries:       []TDSEntry{entry},
	})
	require.NoError(t, err)
	assert.Len(t, payload.Errors, 1)
	assert.Contains(t, payload.Errors[0], "missing PAN")
}

func TestGenerateForm138FVU_Format(t *testing.T) {
	deductees := makeSalaryDeductees()[:2]
	entries := makeSalaryEntries(deductees)

	payload, err := GenerateForm138(Form138Input{
		Employer:      testEmployer(),
		FinancialYear: "2026-27",
		Quarter:       "Q1",
		Deductees:     deductees,
		Entries:       entries,
	})
	require.NoError(t, err)

	fvu := GenerateForm138FVU(payload)
	assert.Contains(t, fvu, "^FH^138^")
	assert.Contains(t, fvu, "MUMA12345A")
	assert.Contains(t, fvu, "2026-27")
	assert.Contains(t, fvu, "^SD^")
	assert.Contains(t, fvu, "^BH^")
}

func TestTaxYearFromFY(t *testing.T) {
	assert.Equal(t, "2026-27", TaxYearFromFY("2026-27"))
	assert.Equal(t, "2027-28", TaxYearFromFY("2027-28"))
	assert.Equal(t, "abc", TaxYearFromFY("abc"))
}

func TestFilterBySection(t *testing.T) {
	entries := []TDSEntry{
		{Section: Section392}, {Section: Section393_1}, {Section: Section392}, {Section: Section393_2},
	}
	assert.Len(t, filterBySection(entries, Section392), 2)
	assert.Len(t, filterBySection(entries, Section393_1), 1)
	assert.Len(t, filterBySection(entries, Section393_3), 0)
}

func TestGroupEntriesByDeductee(t *testing.T) {
	id1 := uuid.New()
	id2 := uuid.New()
	entries := []TDSEntry{
		{DeducteeID: id1}, {DeducteeID: id2}, {DeducteeID: id1},
	}
	grouped := groupEntriesByDeductee(entries)
	assert.Len(t, grouped[id1.String()], 2)
	assert.Len(t, grouped[id2.String()], 1)
}

func TestGenerateForm138_DepositedEntry(t *testing.T) {
	d := Deductee{
		ID: uuid.MustParse("aaaa0001-0001-0001-0001-000000000001"),
		TenantID: testTenant, VendorID: uuid.New(), Name: "Deposited Test",
		PAN: "ABCPD1234E", DeducteeType: DeducteeIndividual, ResidentStatus: Resident,
	}
	entries := []TDSEntry{
		{
			ID: uuid.New(), TenantID: testTenant, DeducteeID: d.ID,
			Section: Section392, PaymentCode: CodeSalaryPrivate,
			FinancialYear: "2026-27", TaxYear: "2026-27", Quarter: "Q1",
			TransactionDate: testDate, GrossAmount: decimal.NewFromInt(200000),
			TDSRate: decimal.NewFromFloat(0.10), TDSAmount: decimal.NewFromInt(20000),
			TotalTax: decimal.NewFromInt(20000), Status: StatusDeposited,
			PANAtDeduction: d.PAN, ChallanNumber: "CHN-001", BSRCode: "BSR001",
		},
	}

	payload, err := GenerateForm138(Form138Input{
		Employer:      testEmployer(),
		FinancialYear: "2026-27",
		Quarter:       "Q1",
		Deductees:     []Deductee{d},
		Entries:       entries,
	})
	require.NoError(t, err)
	assert.True(t, payload.Employees[0].TDSDeposited.Equal(decimal.NewFromInt(20000)))
}
