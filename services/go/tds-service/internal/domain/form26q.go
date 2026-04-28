package domain

import (
	"fmt"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

type DeductorDetails struct {
	TAN            string `json:"tan"`
	DeductorName   string `json:"deductor_name"`
	DeductorPAN    string `json:"deductor_pan"`
	Address        string `json:"address"`
	City           string `json:"city"`
	State          string `json:"state"`
	Pincode        string `json:"pincode"`
	ContactEmail   string `json:"contact_email"`
	ContactPhone   string `json:"contact_phone"`
	ResponsiblePAN string `json:"responsible_pan"`
}

type NonSalaryDetail struct {
	DeducteeID      string          `json:"deductee_id"`
	PAN             string          `json:"pan"`
	Name            string          `json:"name"`
	DeducteeType    DeducteeType    `json:"deductee_type"`
	Section         Section         `json:"section"`
	GrossAmount     decimal.Decimal `json:"gross_amount"`
	TDSRate         decimal.Decimal `json:"tds_rate"`
	TDSAmount       decimal.Decimal `json:"tds_amount"`
	Surcharge       decimal.Decimal `json:"surcharge"`
	Cess            decimal.Decimal `json:"cess"`
	TotalTax        decimal.Decimal `json:"total_tax"`
	NatureOfPayment string          `json:"nature_of_payment"`
	DateOfPayment   time.Time       `json:"date_of_payment"`
	DateOfDeduction time.Time       `json:"date_of_deduction"`
	ChallanNumber   string          `json:"challan_number"`
	BSRCode         string          `json:"bsr_code"`
	NoPAN           bool            `json:"no_pan"`
	LowerCert       bool            `json:"lower_cert"`
}

type Form26QPayload struct {
	FormType      FormType          `json:"form_type"`
	FinancialYear string            `json:"financial_year"`
	Quarter       string            `json:"quarter"`
	Deductor      DeductorDetails   `json:"deductor"`
	Deductions    []NonSalaryDetail `json:"deductions"`
	TotalTDS      decimal.Decimal   `json:"total_tds"`
	TotalPaid     decimal.Decimal   `json:"total_paid"`
	CreatedAt     time.Time         `json:"created_at"`
	Errors        []string          `json:"errors,omitempty"`
}

type Form26QInput struct {
	Deductor      DeductorDetails
	FinancialYear string
	Quarter       string
	Deductees     []Deductee
	Entries       []TDSEntry
}

var nonSalarySections = map[Section]bool{
	Section194C: true,
	Section194I: true,
	Section194J: true,
	Section194Q: true,
}

func GenerateForm26Q(input Form26QInput) (*Form26QPayload, error) {
	if input.Deductor.TAN == "" {
		return nil, fmt.Errorf("deductor TAN is required")
	}
	if input.FinancialYear == "" || input.Quarter == "" {
		return nil, fmt.Errorf("financial_year and quarter are required")
	}

	var nonSalaryEntries []TDSEntry
	for _, e := range input.Entries {
		if nonSalarySections[e.Section] {
			nonSalaryEntries = append(nonSalaryEntries, e)
		}
	}
	if len(nonSalaryEntries) == 0 {
		return nil, fmt.Errorf("no non-salary TDS entries found for %s %s", input.FinancialYear, input.Quarter)
	}

	deducteeMap := make(map[string]*Deductee)
	for i := range input.Deductees {
		deducteeMap[input.Deductees[i].ID.String()] = &input.Deductees[i]
	}

	payload := &Form26QPayload{
		FormType:      FormType26Q,
		FinancialYear: input.FinancialYear,
		Quarter:       input.Quarter,
		Deductor:      input.Deductor,
		CreatedAt:     time.Now(),
	}

	totalTDS := decimal.Zero
	totalPaid := decimal.Zero
	var validationErrors []string

	for _, entry := range nonSalaryEntries {
		d, ok := deducteeMap[entry.DeducteeID.String()]
		if !ok {
			validationErrors = append(validationErrors, fmt.Sprintf("deductee %s not found", entry.DeducteeID))
			continue
		}
		if d.PAN == "" && !entry.NoPANDeduction {
			validationErrors = append(validationErrors, fmt.Sprintf("deductee %s (%s) missing PAN", d.Name, entry.DeducteeID))
		}

		detail := NonSalaryDetail{
			DeducteeID:      entry.DeducteeID.String(),
			PAN:             d.PAN,
			Name:            d.Name,
			DeducteeType:    d.DeducteeType,
			Section:         entry.Section,
			GrossAmount:     entry.GrossAmount,
			TDSRate:         entry.TDSRate,
			TDSAmount:       entry.TDSAmount,
			Surcharge:       entry.Surcharge,
			Cess:            entry.Cess,
			TotalTax:        entry.TotalTax,
			NatureOfPayment: entry.NatureOfPayment,
			DateOfPayment:   entry.TransactionDate,
			DateOfDeduction: entry.TransactionDate,
			ChallanNumber:   entry.ChallanNumber,
			BSRCode:         entry.BSRCode,
			NoPAN:           entry.NoPANDeduction,
			LowerCert:       entry.LowerCertApplied,
		}
		payload.Deductions = append(payload.Deductions, detail)
		totalTDS = totalTDS.Add(entry.TotalTax)
		totalPaid = totalPaid.Add(entry.GrossAmount)
	}

	payload.TotalTDS = totalTDS
	payload.TotalPaid = totalPaid
	payload.Errors = validationErrors

	return payload, nil
}

func GenerateForm26QFVU(payload *Form26QPayload) string {
	var b strings.Builder

	ay := assessmentYear(payload.FinancialYear)

	b.WriteString(fmt.Sprintf("^FH^26Q^1^%s^%s^%s^%s^^%s^^%d^^%s^%s^\n",
		payload.Deductor.TAN,
		payload.Deductor.DeductorPAN,
		ay,
		payload.Quarter,
		payload.Deductor.DeductorName,
		len(payload.Deductions),
		payload.Deductor.Address,
		payload.Deductor.Pincode,
	))

	for _, ch := range challanSummary26Q(payload) {
		b.WriteString(fmt.Sprintf("^BH^%s^%s^%s^%s^%s^\n",
			ch.ChallanNumber, ch.BSRCode, ch.DepositDate,
			ch.TotalTDS, ch.DeducteeCount,
		))
	}

	for i, ded := range payload.Deductions {
		noPANFlag := "N"
		if ded.NoPAN {
			noPANFlag = "Y"
		}
		b.WriteString(fmt.Sprintf("^DD^%d^%s^%s^%s^%s^%s^%s^%s^%s^%s^%s^\n",
			i+1,
			ded.PAN,
			ded.Name,
			string(ded.Section),
			ded.GrossAmount.StringFixed(2),
			ded.TDSRate.StringFixed(4),
			ded.TDSAmount.StringFixed(2),
			ded.Surcharge.StringFixed(2),
			ded.Cess.StringFixed(2),
			ded.TotalTax.StringFixed(2),
			noPANFlag,
		))
	}

	return b.String()
}

func challanSummary26Q(payload *Form26QPayload) []challanLine {
	challanMap := make(map[string]*challanLine)
	for _, ded := range payload.Deductions {
		key := ded.ChallanNumber
		if key == "" {
			key = "PENDING"
		}
		if existing, ok := challanMap[key]; ok {
			amt, _ := decimal.NewFromString(existing.TotalTDS)
			cnt, _ := decimal.NewFromString(existing.DeducteeCount)
			existing.TotalTDS = amt.Add(ded.TotalTax).StringFixed(2)
			existing.DeducteeCount = cnt.Add(decimal.NewFromInt(1)).String()
		} else {
			challanMap[key] = &challanLine{
				ChallanNumber: key,
				BSRCode:       ded.BSRCode,
				DepositDate:   ded.DateOfPayment.Format("02012006"),
				TotalTDS:      ded.TotalTax.StringFixed(2),
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
