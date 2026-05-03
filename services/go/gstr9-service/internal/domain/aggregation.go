package domain

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

var (
	GSTR9MandatoryThreshold = decimal.NewFromInt(20000000)
	GSTR9CThreshold         = decimal.NewFromInt(50000000)
)

type ThresholdResult struct {
	GSTR9Mandatory  bool            `json:"gstr9_mandatory"`
	GSTR9CRequired  bool            `json:"gstr9c_required"`
	AggregateTurnover decimal.Decimal `json:"aggregate_turnover"`
	Reason          string          `json:"reason"`
}

func CheckThreshold(turnover decimal.Decimal) ThresholdResult {
	r := ThresholdResult{AggregateTurnover: turnover}
	if turnover.GreaterThan(GSTR9MandatoryThreshold) {
		r.GSTR9Mandatory = true
		r.Reason = fmt.Sprintf(
			"aggregate turnover ₹%s exceeds ₹2 crore threshold (Notification 15/2025-CT)",
			turnover.StringFixed(2),
		)
	} else {
		r.Reason = fmt.Sprintf(
			"aggregate turnover ₹%s is within ₹2 crore — GSTR-9 filing optional",
			turnover.StringFixed(2),
		)
	}
	if turnover.GreaterThan(GSTR9CThreshold) {
		r.GSTR9CRequired = true
		r.Reason += "; GSTR-9C reconciliation mandatory (turnover exceeds ₹5 crore)"
	}
	return r
}

type MonthlyData struct {
	ReturnPeriod   string
	Outward        TaxBreakdown
	Inward         TaxBreakdown
	ITC            ITCBreakdown
	TaxPaid        TaxBreakdown
	LateITCReclaim ITCBreakdown
	Rule37Reclaim  ITCBreakdown
}

type TaxBreakdown struct {
	TaxableValue decimal.Decimal `json:"taxable_value"`
	CGST         decimal.Decimal `json:"cgst"`
	SGST         decimal.Decimal `json:"sgst"`
	IGST         decimal.Decimal `json:"igst"`
	Cess         decimal.Decimal `json:"cess"`
}

func (t TaxBreakdown) Add(o TaxBreakdown) TaxBreakdown {
	return TaxBreakdown{
		TaxableValue: t.TaxableValue.Add(o.TaxableValue),
		CGST:         t.CGST.Add(o.CGST),
		SGST:         t.SGST.Add(o.SGST),
		IGST:         t.IGST.Add(o.IGST),
		Cess:         t.Cess.Add(o.Cess),
	}
}

func (t TaxBreakdown) TotalTax() decimal.Decimal {
	return t.CGST.Add(t.SGST).Add(t.IGST).Add(t.Cess)
}

type ITCBreakdown struct {
	CGST decimal.Decimal `json:"cgst"`
	SGST decimal.Decimal `json:"sgst"`
	IGST decimal.Decimal `json:"igst"`
	Cess decimal.Decimal `json:"cess"`
}

func (i ITCBreakdown) Add(o ITCBreakdown) ITCBreakdown {
	return ITCBreakdown{
		CGST: i.CGST.Add(o.CGST),
		SGST: i.SGST.Add(o.SGST),
		IGST: i.IGST.Add(o.IGST),
		Cess: i.Cess.Add(o.Cess),
	}
}

func (i ITCBreakdown) Total() decimal.Decimal {
	return i.CGST.Add(i.SGST).Add(i.IGST).Add(i.Cess)
}

type AggregationResult struct {
	FilingID      uuid.UUID
	FinancialYear string
	GSTIN         string
	MonthsFound   int
	Tables        []GSTR9TableData
	Turnover      decimal.Decimal
	Threshold     ThresholdResult
}
