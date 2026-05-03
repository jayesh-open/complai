package domain

import (
	"fmt"

	"github.com/shopspring/decimal"
)

type ITR1Form struct {
	FormType       ITRFormType        `json:"form_type"`
	TaxYear        string             `json:"tax_year"`
	PAN            string             `json:"pan"`
	Name           string             `json:"name"`
	Regime         RegimeType         `json:"regime"`
	Salary         SalaryResult       `json:"schedule_salary"`
	HouseProperty  HousePropertyResult `json:"schedule_hp"`
	OtherSources   OtherSourcesResult  `json:"schedule_os"`
	LTCG112A       decimal.Decimal    `json:"ltcg_112a"`
	AgriculturalIncome decimal.Decimal `json:"agricultural_income"`
	TaxComputation TaxComputeResult   `json:"tax_computation"`
	TDSCredits     []TDSCredit        `json:"tds_credits"`
	AISPreFill     *AISPreFillBlock   `json:"ais_prefill,omitempty"`
	Verification   VerificationBlock  `json:"verification"`
}

type AISPreFillBlock struct {
	Form168Ref     string          `json:"form_168_ref"`
	TDSTotal       decimal.Decimal `json:"tds_total"`
	InterestIncome decimal.Decimal `json:"interest_income"`
	DividendIncome decimal.Decimal `json:"dividend_income"`
	SalaryIncome   decimal.Decimal `json:"salary_income"`
}

type VerificationBlock struct {
	Method VerificationMethod `json:"method"`
	Date   string             `json:"date,omitempty"`
	Place  string             `json:"place,omitempty"`
}

type ITR1Eligibility struct {
	Eligible bool   `json:"eligible"`
	Reason   string `json:"reason,omitempty"`
}

var itr1IncomeLimit = decimal.NewFromInt(5000000)
var itr1LTCG112ALimit = decimal.NewFromInt(125000)
var itr1AgriculturalLimit = decimal.NewFromInt(5000)

func CheckITR1Eligibility(
	assesseeType AssesseeType,
	residency ResidencyStatus,
	totalIncome decimal.Decimal,
	hpCount int,
	hasCapitalGainsOver112A bool,
	ltcg112A decimal.Decimal,
	hasBusiness bool,
	hasForeignAssets bool,
	hasUnlistedEquity bool,
	isDirector bool,
	agriculturalIncome decimal.Decimal,
) ITR1Eligibility {
	if assesseeType != AssesseeIndividual {
		return ITR1Eligibility{Reason: "ITR-1 is only for individuals, not HUF"}
	}
	if residency != Resident {
		return ITR1Eligibility{Reason: "ITR-1 is only for resident individuals"}
	}
	if totalIncome.GreaterThan(itr1IncomeLimit) {
		return ITR1Eligibility{Reason: fmt.Sprintf("total income exceeds ₹%s limit", itr1IncomeLimit.StringFixed(0))}
	}
	if hpCount > 2 {
		return ITR1Eligibility{Reason: "ITR-1 supports at most 2 house properties under ITA 2025"}
	}
	if hasCapitalGainsOver112A {
		return ITR1Eligibility{Reason: "capital gains beyond Section 112A LTCG not allowed in ITR-1"}
	}
	if ltcg112A.GreaterThan(itr1LTCG112ALimit) {
		return ITR1Eligibility{Reason: fmt.Sprintf("LTCG under Section 112A exceeds ₹%s", itr1LTCG112ALimit.StringFixed(0))}
	}
	if hasBusiness {
		return ITR1Eligibility{Reason: "business/profession income not allowed in ITR-1"}
	}
	if hasForeignAssets {
		return ITR1Eligibility{Reason: "foreign assets not allowed in ITR-1"}
	}
	if hasUnlistedEquity {
		return ITR1Eligibility{Reason: "holder of unlisted equity not allowed in ITR-1"}
	}
	if isDirector {
		return ITR1Eligibility{Reason: "director of a company cannot use ITR-1"}
	}
	if agriculturalIncome.GreaterThan(itr1AgriculturalLimit) {
		return ITR1Eligibility{Reason: fmt.Sprintf("agricultural income exceeds ₹%s", itr1AgriculturalLimit.StringFixed(0))}
	}
	return ITR1Eligibility{Eligible: true}
}

type ITR2Form struct {
	FormType       ITRFormType         `json:"form_type"`
	TaxYear        string              `json:"tax_year"`
	PAN            string              `json:"pan"`
	Name           string              `json:"name"`
	AssesseeType   AssesseeType        `json:"assessee_type"`
	Regime         RegimeType          `json:"regime"`
	Salary         SalaryResult        `json:"schedule_salary"`
	HouseProperty  []HousePropertyResult `json:"schedule_hp"`
	CapitalGains   CapitalGainsResult  `json:"schedule_cg"`
	OtherSources   OtherSourcesResult  `json:"schedule_os"`
	Schedule112A   *Schedule112A       `json:"schedule_112a,omitempty"`
	ScheduleVDA    ScheduleVDA         `json:"schedule_vda"`
	ScheduleFA     *ScheduleFA         `json:"schedule_fa,omitempty"`
	TaxComputation TaxComputeResult    `json:"tax_computation"`
	TDSCredits     []TDSCredit         `json:"tds_credits"`
	AISPreFill     *AISPreFillBlock    `json:"ais_prefill,omitempty"`
	Verification   VerificationBlock   `json:"verification"`
}

type Schedule112A struct {
	Entries        []LTCG112AEntry `json:"entries"`
	TotalLTCG      decimal.Decimal `json:"total_ltcg"`
	ExemptionUsed  decimal.Decimal `json:"exemption_used"`
	TaxableLTCG    decimal.Decimal `json:"taxable_ltcg"`
}

type LTCG112AEntry struct {
	ISIN          string          `json:"isin"`
	SecurityName  string          `json:"security_name"`
	Quantity      int             `json:"quantity"`
	SalePrice     decimal.Decimal `json:"sale_price"`
	CostBasis     decimal.Decimal `json:"cost_basis"`
	Gain          decimal.Decimal `json:"gain"`
}

type ScheduleVDA struct {
	HasTransactions bool            `json:"has_transactions"`
	Entries         []VDAEntry      `json:"entries,omitempty"`
	TotalGain       decimal.Decimal `json:"total_gain"`
}

type VDAEntry struct {
	AssetName      string          `json:"asset_name"`
	DateOfTransfer string          `json:"date_of_transfer"`
	SaleAmount     decimal.Decimal `json:"sale_amount"`
	CostBasis      decimal.Decimal `json:"cost_basis"`
	Gain           decimal.Decimal `json:"gain"`
}

type ScheduleFA struct {
	ForeignBankAccounts  int `json:"foreign_bank_accounts"`
	ForeignEquity        int `json:"foreign_equity"`
	ForeignProperty      int `json:"foreign_property"`
	OtherForeignAssets   int `json:"other_foreign_assets"`
	SigningAuthority     int `json:"signing_authority_accounts"`
}

type ITR2Eligibility struct {
	Eligible bool   `json:"eligible"`
	Reason   string `json:"reason,omitempty"`
}

func CheckITR2Eligibility(
	assesseeType AssesseeType,
	hasBusiness bool,
) ITR2Eligibility {
	if assesseeType != AssesseeIndividual && assesseeType != AssesseeHUF {
		return ITR2Eligibility{Reason: "ITR-2 is for individuals and HUFs only"}
	}
	if hasBusiness {
		return ITR2Eligibility{Reason: "business/profession income requires ITR-3"}
	}
	return ITR2Eligibility{Eligible: true}
}

type ITR3Form struct {
	FormType       ITRFormType           `json:"form_type"`
	TaxYear        string                `json:"tax_year"`
	PAN            string                `json:"pan"`
	Name           string                `json:"name"`
	AssesseeType   AssesseeType          `json:"assessee_type"`
	Regime         RegimeType            `json:"regime"`
	Form10IEARef   string                `json:"form_10iea_ref,omitempty"`
	Salary         SalaryResult          `json:"schedule_salary"`
	HouseProperty  []HousePropertyResult `json:"schedule_hp"`
	Business       BusinessResult        `json:"schedule_bp"`
	CapitalGains   CapitalGainsResult    `json:"schedule_cg"`
	OtherSources   OtherSourcesResult    `json:"schedule_os"`
	Schedule112A   *Schedule112A         `json:"schedule_112a,omitempty"`
	ScheduleVDA    ScheduleVDA           `json:"schedule_vda"`
	ScheduleFA     *ScheduleFA           `json:"schedule_fa,omitempty"`
	ScheduleBP     ScheduleBP            `json:"schedule_bp_detail"`
	ScheduleTDSIT  []ScheduleTDSITEntry  `json:"schedule_tds_it"`
	TaxComputation TaxComputeResult      `json:"tax_computation"`
	TDSCredits     []TDSCredit           `json:"tds_credits"`
	AISPreFill     *AISPreFillBlock      `json:"ais_prefill,omitempty"`
	Verification   VerificationBlock     `json:"verification"`
}

type ScheduleBP struct {
	GrossTurnover     decimal.Decimal `json:"gross_turnover"`
	GrossReceipts     decimal.Decimal `json:"gross_receipts"`
	TotalExpenses     decimal.Decimal `json:"total_expenses"`
	Depreciation      decimal.Decimal `json:"depreciation"`
	NetProfit         decimal.Decimal `json:"net_profit"`
	Section44BBCApply bool            `json:"section_44bbc_apply"`
}

type ScheduleTDSITEntry struct {
	DeductorTAN   string          `json:"deductor_tan"`
	DeductorName  string          `json:"deductor_name"`
	Section       string          `json:"section"`
	TotalIncome   decimal.Decimal `json:"total_income"`
	TDSDeducted   decimal.Decimal `json:"tds_deducted"`
	TDSClaimed    decimal.Decimal `json:"tds_claimed"`
}

type ITR3Eligibility struct {
	Eligible bool   `json:"eligible"`
	Reason   string `json:"reason,omitempty"`
}

func CheckITR3Eligibility(assesseeType AssesseeType) ITR3Eligibility {
	if assesseeType != AssesseeIndividual && assesseeType != AssesseeHUF {
		return ITR3Eligibility{Reason: "ITR-3 is for individuals and HUFs only"}
	}
	return ITR3Eligibility{Eligible: true}
}
