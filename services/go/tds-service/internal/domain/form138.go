package domain

import (
	"fmt"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

type EmployerDetails struct {
	TAN            string `json:"tan"`
	EmployerName   string `json:"employer_name"`
	EmployerPAN    string `json:"employer_pan"`
	Address        string `json:"address"`
	City           string `json:"city"`
	State          string `json:"state"`
	Pincode        string `json:"pincode"`
	ContactEmail   string `json:"contact_email"`
	ContactPhone   string `json:"contact_phone"`
	ResponsiblePAN string `json:"responsible_pan"`
}

type SalaryDetail struct {
	DeducteeID       string          `json:"deductee_id"`
	PAN              string          `json:"pan"`
	Name             string          `json:"name"`
	PaymentCode      PaymentCode     `json:"payment_code"`
	Designation      string          `json:"designation"`
	GrossSalary      decimal.Decimal `json:"gross_salary"`
	ExemptAllowances decimal.Decimal `json:"exempt_allowances"`
	NetSalary        decimal.Decimal `json:"net_salary"`
	StdDeduction     decimal.Decimal `json:"std_deduction"`
	TaxableIncome    decimal.Decimal `json:"taxable_income"`
	TDSDeducted      decimal.Decimal `json:"tds_deducted"`
	TDSDeposited     decimal.Decimal `json:"tds_deposited"`
	Surcharge        decimal.Decimal `json:"surcharge"`
	Cess             decimal.Decimal `json:"cess"`
	TotalTax         decimal.Decimal `json:"total_tax"`
	DateOfPayment    time.Time       `json:"date_of_payment"`
	DateOfDeduction  time.Time       `json:"date_of_deduction"`
	ChallanNumber    string          `json:"challan_number"`
	BSRCode          string          `json:"bsr_code"`
}

type Form138Payload struct {
	FormType      FormType        `json:"form_type"`
	FinancialYear string          `json:"financial_year"`
	Quarter       string          `json:"quarter"`
	Employer      EmployerDetails `json:"employer"`
	Employees     []SalaryDetail  `json:"employees"`
	TotalTDS      decimal.Decimal `json:"total_tds"`
	TotalSalary   decimal.Decimal `json:"total_salary"`
	CreatedAt     time.Time       `json:"created_at"`
	Errors        []string        `json:"errors,omitempty"`
}

type Form138Input struct {
	Employer      EmployerDetails
	FinancialYear string
	Quarter       string
	Deductees     []Deductee
	Entries       []TDSEntry
}

func GenerateForm138(input Form138Input) (*Form138Payload, error) {
	if input.Employer.TAN == "" {
		return nil, fmt.Errorf("employer TAN is required")
	}
	if input.FinancialYear == "" || input.Quarter == "" {
		return nil, fmt.Errorf("financial_year and quarter are required")
	}

	salaryEntries := filterBySection(input.Entries, Section392)
	if len(salaryEntries) == 0 {
		return nil, fmt.Errorf("no salary entries (s.392) found for %s %s", input.FinancialYear, input.Quarter)
	}

	deducteeMap := make(map[string]*Deductee)
	for i := range input.Deductees {
		deducteeMap[input.Deductees[i].ID.String()] = &input.Deductees[i]
	}

	payload := &Form138Payload{
		FormType:      FormType138,
		FinancialYear: input.FinancialYear,
		Quarter:       input.Quarter,
		Employer:      input.Employer,
		CreatedAt:     time.Now(),
	}

	totalTDS := decimal.Zero
	totalSalary := decimal.Zero
	var validationErrors []string

	byDeductee := groupEntriesByDeductee(salaryEntries)

	for deducteeID, entries := range byDeductee {
		d, ok := deducteeMap[deducteeID]
		if !ok {
			validationErrors = append(validationErrors, fmt.Sprintf("deductee %s not found in provided deductee list", deducteeID))
			continue
		}
		if d.PAN == "" {
			validationErrors = append(validationErrors, fmt.Sprintf("deductee %s (%s) missing PAN", d.Name, deducteeID))
		}

		grossSalary := decimal.Zero
		tdsDeducted := decimal.Zero
		tdsDeposited := decimal.Zero
		surcharge := decimal.Zero
		cess := decimal.Zero
		totalTax := decimal.Zero
		var latestEntry TDSEntry
		for _, e := range entries {
			grossSalary = grossSalary.Add(e.GrossAmount)
			tdsDeducted = tdsDeducted.Add(e.TDSAmount)
			surcharge = surcharge.Add(e.Surcharge)
			cess = cess.Add(e.Cess)
			totalTax = totalTax.Add(e.TotalTax)
			if e.Status == StatusDeposited {
				tdsDeposited = tdsDeposited.Add(e.TotalTax)
			}
			latestEntry = e
		}

		netSalary := grossSalary
		stdDed := decimal.NewFromInt(75000).Div(decimal.NewFromInt(4))
		taxableIncome := netSalary.Sub(stdDed)
		if taxableIncome.IsNegative() {
			taxableIncome = decimal.Zero
		}

		detail := SalaryDetail{
			DeducteeID:       deducteeID,
			PAN:              d.PAN,
			Name:             d.Name,
			PaymentCode:      latestEntry.PaymentCode,
			GrossSalary:      grossSalary,
			ExemptAllowances: decimal.Zero,
			NetSalary:        netSalary,
			StdDeduction:     stdDed,
			TaxableIncome:    taxableIncome,
			TDSDeducted:      tdsDeducted,
			TDSDeposited:     tdsDeposited,
			Surcharge:        surcharge,
			Cess:             cess,
			TotalTax:         totalTax,
			DateOfPayment:    latestEntry.TransactionDate,
			DateOfDeduction:  latestEntry.TransactionDate,
			ChallanNumber:    latestEntry.ChallanNumber,
			BSRCode:          latestEntry.BSRCode,
		}
		payload.Employees = append(payload.Employees, detail)
		totalTDS = totalTDS.Add(totalTax)
		totalSalary = totalSalary.Add(grossSalary)
	}

	payload.TotalTDS = totalTDS
	payload.TotalSalary = totalSalary
	payload.Errors = validationErrors

	return payload, nil
}

func GenerateForm138FVU(payload *Form138Payload) string {
	var b strings.Builder

	ty := TaxYearFromFY(payload.FinancialYear)

	b.WriteString(fmt.Sprintf("^FH^138^1^%s^%s^%s^%s^^%s^^%d^^%s^%s^\n",
		payload.Employer.TAN,
		payload.Employer.EmployerPAN,
		ty,
		payload.Quarter,
		payload.Employer.EmployerName,
		len(payload.Employees),
		payload.Employer.Address,
		payload.Employer.Pincode,
	))

	for _, ch := range challanSummary138(payload) {
		b.WriteString(fmt.Sprintf("^BH^%s^%s^%s^%s^%s^\n",
			ch.ChallanNumber, ch.BSRCode, ch.DepositDate,
			ch.TotalTDS, ch.DeducteeCount,
		))
	}

	for i, emp := range payload.Employees {
		b.WriteString(fmt.Sprintf("^SD^%d^%s^%s^%s^%s^%s^%s^%s^%s^%s^\n",
			i+1,
			emp.PAN,
			emp.Name,
			string(emp.PaymentCode),
			emp.GrossSalary.StringFixed(2),
			emp.TaxableIncome.StringFixed(2),
			emp.TDSDeducted.StringFixed(2),
			emp.Surcharge.StringFixed(2),
			emp.Cess.StringFixed(2),
			emp.TotalTax.StringFixed(2),
		))
	}

	return b.String()
}

type challanLine struct {
	ChallanNumber string
	BSRCode       string
	DepositDate   string
	TotalTDS      string
	DeducteeCount string
}

func challanSummary138(payload *Form138Payload) []challanLine {
	challanMap := make(map[string]*challanLine)
	for _, emp := range payload.Employees {
		key := emp.ChallanNumber
		if key == "" {
			key = "PENDING"
		}
		if existing, ok := challanMap[key]; ok {
			amt, _ := decimal.NewFromString(existing.TotalTDS)
			cnt, _ := decimal.NewFromString(existing.DeducteeCount)
			existing.TotalTDS = amt.Add(emp.TotalTax).StringFixed(2)
			existing.DeducteeCount = cnt.Add(decimal.NewFromInt(1)).String()
		} else {
			challanMap[key] = &challanLine{
				ChallanNumber: key,
				BSRCode:       emp.BSRCode,
				DepositDate:   emp.DateOfPayment.Format("02012006"),
				TotalTDS:      emp.TotalTax.StringFixed(2),
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

func filterBySection(entries []TDSEntry, section Section) []TDSEntry {
	var filtered []TDSEntry
	for _, e := range entries {
		if e.Section == section {
			filtered = append(filtered, e)
		}
	}
	return filtered
}

func groupEntriesByDeductee(entries []TDSEntry) map[string][]TDSEntry {
	grouped := make(map[string][]TDSEntry)
	for _, e := range entries {
		grouped[e.DeducteeID.String()] = append(grouped[e.DeducteeID.String()], e)
	}
	return grouped
}
