package domain

import (
	"fmt"

	"github.com/shopspring/decimal"
)

var itr4IncomeLimit = decimal.NewFromInt(5000000)
var itr4LTCG112ALimit = decimal.NewFromInt(125000)

type ITR4Form struct {
	FormType           ITRFormType         `json:"form_type"`
	TaxYear            string              `json:"tax_year"`
	PAN                string              `json:"pan"`
	Name               string              `json:"name"`
	AssesseeType       AssesseeType        `json:"assessee_type"`
	Regime             RegimeType          `json:"regime"`
	Salary             SalaryResult        `json:"schedule_salary"`
	HouseProperty      HousePropertyResult `json:"schedule_hp"`
	OtherSources       OtherSourcesResult  `json:"schedule_os"`
	PresumptiveBP      PresumptiveSchedule `json:"schedule_bp_presumptive"`
	LTCG112A           decimal.Decimal     `json:"ltcg_112a"`
	TaxComputation     TaxComputeResult    `json:"tax_computation"`
	TDSCredits         []TDSCredit         `json:"tds_credits"`
	AISPreFill         *AISPreFillBlock    `json:"ais_prefill,omitempty"`
	Verification       VerificationBlock   `json:"verification"`
}

type PresumptiveSchedule struct {
	Section44AD  PresumptiveEntry `json:"section_44ad,omitempty"`
	Section44ADA PresumptiveEntry `json:"section_44ada,omitempty"`
	Section44AE  PresumptiveEntry `json:"section_44ae,omitempty"`
	TotalIncome  decimal.Decimal  `json:"total_income"`
}

type PresumptiveEntry struct {
	GrossTurnover   decimal.Decimal `json:"gross_turnover"`
	PresumptiveRate decimal.Decimal `json:"presumptive_rate"`
	PresumptiveInc  decimal.Decimal `json:"presumptive_income"`
}

type ITR4Eligibility struct {
	Eligible bool   `json:"eligible"`
	Reason   string `json:"reason,omitempty"`
}

func CheckITR4Eligibility(
	assesseeType AssesseeType,
	residency ResidencyStatus,
	totalIncome decimal.Decimal,
	hasPresumptive bool,
	hasForeignAssets bool,
	isDirector bool,
	ltcg112A decimal.Decimal,
) ITR4Eligibility {
	if assesseeType != AssesseeIndividual && assesseeType != AssesseeHUF && assesseeType != AssesseeFirm {
		return ITR4Eligibility{Reason: "ITR-4 is for individuals, HUFs, and firms (not LLPs)"}
	}
	if residency != Resident {
		return ITR4Eligibility{Reason: "ITR-4 is only for residents"}
	}
	if !hasPresumptive {
		return ITR4Eligibility{Reason: "ITR-4 requires presumptive income under Section 44AD/44ADA/44AE"}
	}
	if totalIncome.GreaterThan(itr4IncomeLimit) {
		return ITR4Eligibility{Reason: fmt.Sprintf("total income exceeds ₹%s limit for ITR-4", itr4IncomeLimit.StringFixed(0))}
	}
	if hasForeignAssets {
		return ITR4Eligibility{Reason: "foreign assets not allowed in ITR-4"}
	}
	if isDirector {
		return ITR4Eligibility{Reason: "director of a company cannot use ITR-4"}
	}
	if ltcg112A.GreaterThan(itr4LTCG112ALimit) {
		return ITR4Eligibility{Reason: fmt.Sprintf("LTCG under Section 112A exceeds ₹%s", itr4LTCG112ALimit.StringFixed(0))}
	}
	return ITR4Eligibility{Eligible: true}
}

type ITR5Form struct {
	FormType       ITRFormType           `json:"form_type"`
	TaxYear        string                `json:"tax_year"`
	PAN            string                `json:"pan"`
	Name           string                `json:"name"`
	AssesseeType   AssesseeType          `json:"assessee_type"`
	ScheduleBP     ScheduleBP            `json:"schedule_bp"`
	HouseProperty  []HousePropertyResult `json:"schedule_hp"`
	CapitalGains   CapitalGainsResult    `json:"schedule_cg"`
	OtherSources   OtherSourcesResult    `json:"schedule_os"`
	ScheduleVDA    ScheduleVDA           `json:"schedule_vda"`
	ScheduleFA     *ScheduleFA           `json:"schedule_fa,omitempty"`
	PartnerDetails []PartnerDetail       `json:"partner_details,omitempty"`
	TaxComputation TaxComputeResult      `json:"tax_computation"`
	TDSCredits     []TDSCredit           `json:"tds_credits"`
	AISPreFill     *AISPreFillBlock      `json:"ais_prefill,omitempty"`
	Verification   VerificationBlock     `json:"verification"`
}

type PartnerDetail struct {
	Name         string          `json:"name"`
	PAN          string          `json:"pan"`
	SharePercent decimal.Decimal `json:"share_percent"`
	ShareAmount  decimal.Decimal `json:"share_amount"`
}

type ITR5Eligibility struct {
	Eligible bool   `json:"eligible"`
	Reason   string `json:"reason,omitempty"`
}

func CheckITR5Eligibility(assesseeType AssesseeType) ITR5Eligibility {
	switch assesseeType {
	case AssesseeFirm, AssesseeLLP, AssesseeAOP, AssesseeBOI:
		return ITR5Eligibility{Eligible: true}
	default:
		return ITR5Eligibility{Reason: "ITR-5 is for firms, LLPs, AOPs, and BOIs only"}
	}
}

type ITR6Form struct {
	FormType        ITRFormType           `json:"form_type"`
	TaxYear         string                `json:"tax_year"`
	PAN             string                `json:"pan"`
	CompanyName     string                `json:"company_name"`
	CIN             string                `json:"cin"`
	ScheduleBP      ScheduleBP            `json:"schedule_bp"`
	HouseProperty   []HousePropertyResult `json:"schedule_hp"`
	CapitalGains    CapitalGainsResult    `json:"schedule_cg"`
	OtherSources    OtherSourcesResult    `json:"schedule_os"`
	ScheduleVDA     ScheduleVDA           `json:"schedule_vda"`
	ScheduleFA      *ScheduleFA           `json:"schedule_fa,omitempty"`
	ScheduleMAT     ScheduleMAT           `json:"schedule_mat"`
	BuybackLoss     decimal.Decimal       `json:"buyback_loss"`
	DeemedDividend  decimal.Decimal       `json:"deemed_dividend_2_22_f"`
	TaxComputation  TaxComputeResult      `json:"tax_computation"`
	TDSCredits      []TDSCredit           `json:"tds_credits"`
	AISPreFill      *AISPreFillBlock      `json:"ais_prefill,omitempty"`
	Verification    VerificationBlock     `json:"verification"`
}

type ScheduleMAT struct {
	BookProfit     decimal.Decimal `json:"book_profit"`
	MATRate        decimal.Decimal `json:"mat_rate"`
	MATPayable     decimal.Decimal `json:"mat_payable"`
	NormalTax      decimal.Decimal `json:"normal_tax"`
	TaxApplicable  string          `json:"tax_applicable"`
}

type ITR6Eligibility struct {
	Eligible bool   `json:"eligible"`
	Reason   string `json:"reason,omitempty"`
}

func CheckITR6Eligibility(assesseeType AssesseeType, claimsITR7Exemption bool) ITR6Eligibility {
	if assesseeType != AssesseeCompany {
		return ITR6Eligibility{Reason: "ITR-6 is for companies only"}
	}
	if claimsITR7Exemption {
		return ITR6Eligibility{Reason: "companies claiming exemption under Section 139(4A)/(4B)/(4C)/(4D) must use ITR-7"}
	}
	return ITR6Eligibility{Eligible: true}
}

type ITR7Form struct {
	FormType           ITRFormType           `json:"form_type"`
	TaxYear            string                `json:"tax_year"`
	PAN                string                `json:"pan"`
	EntityName         string                `json:"entity_name"`
	FilingSection      string                `json:"filing_section"`
	ScheduleBP         ScheduleBP            `json:"schedule_bp"`
	ScheduleIncome     OtherSourcesResult    `json:"schedule_income"`
	CapitalGains       CapitalGainsResult    `json:"schedule_cg"`
	AnonymousDonations decimal.Decimal       `json:"anonymous_donations"`
	VoluntaryContrib   decimal.Decimal       `json:"voluntary_contributions"`
	ApplicationOfInc   decimal.Decimal       `json:"application_of_income"`
	AccumPercent       decimal.Decimal       `json:"accumulation_percent"`
	TaxComputation     TaxComputeResult      `json:"tax_computation"`
	TDSCredits         []TDSCredit           `json:"tds_credits"`
	AISPreFill         *AISPreFillBlock      `json:"ais_prefill,omitempty"`
	Verification       VerificationBlock     `json:"verification"`
}

type ITR7Eligibility struct {
	Eligible bool   `json:"eligible"`
	Reason   string `json:"reason,omitempty"`
}

var validITR7Sections = map[string]bool{
	"139(4A)": true,
	"139(4B)": true,
	"139(4C)": true,
	"139(4D)": true,
}

func CheckITR7Eligibility(assesseeType AssesseeType, filingSection string) ITR7Eligibility {
	if assesseeType != AssesseeTrust && assesseeType != AssesseeCompany && assesseeType != AssesseeAOP {
		return ITR7Eligibility{Reason: "ITR-7 is for trusts, institutions, and entities filing under Section 139(4A)/(4B)/(4C)/(4D)"}
	}
	if !validITR7Sections[filingSection] {
		return ITR7Eligibility{Reason: "filing section must be one of 139(4A), 139(4B), 139(4C), 139(4D)"}
	}
	return ITR7Eligibility{Eligible: true}
}
