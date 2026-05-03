package domain

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestCheckITR4Eligibility(t *testing.T) {
	tests := []struct {
		name           string
		assessee       AssesseeType
		residency      ResidencyStatus
		income         decimal.Decimal
		hasPresumptive bool
		foreignAssets  bool
		isDirector     bool
		ltcg           decimal.Decimal
		wantEligible   bool
		wantReason     string
	}{
		{"Eligible_Individual", AssesseeIndividual, Resident, decimal.NewFromInt(3000000), true, false, false, decimal.Zero, true, ""},
		{"Eligible_HUF", AssesseeHUF, Resident, decimal.NewFromInt(2000000), true, false, false, decimal.Zero, true, ""},
		{"Eligible_Firm", AssesseeFirm, Resident, decimal.NewFromInt(1000000), true, false, false, decimal.Zero, true, ""},
		{"LLP_NotAllowed", AssesseeLLP, Resident, decimal.NewFromInt(1000000), true, false, false, decimal.Zero, false, "ITR-4 is for individuals, HUFs, and firms (not LLPs)"},
		{"Company_NotAllowed", AssesseeCompany, Resident, decimal.NewFromInt(1000000), true, false, false, decimal.Zero, false, "ITR-4 is for individuals, HUFs, and firms (not LLPs)"},
		{"NRI_NotAllowed", AssesseeIndividual, NonResident, decimal.NewFromInt(1000000), true, false, false, decimal.Zero, false, "ITR-4 is only for residents"},
		{"NoPresumptive", AssesseeIndividual, Resident, decimal.NewFromInt(1000000), false, false, false, decimal.Zero, false, "ITR-4 requires presumptive income"},
		{"OverIncomeLimit", AssesseeIndividual, Resident, decimal.NewFromInt(6000000), true, false, false, decimal.Zero, false, "total income exceeds"},
		{"ForeignAssets", AssesseeIndividual, Resident, decimal.NewFromInt(1000000), true, true, false, decimal.Zero, false, "foreign assets"},
		{"Director", AssesseeIndividual, Resident, decimal.NewFromInt(1000000), true, false, true, decimal.Zero, false, "director"},
		{"LTCG_Over", AssesseeIndividual, Resident, decimal.NewFromInt(1000000), true, false, false, decimal.NewFromInt(200000), false, "LTCG under Section 112A exceeds"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := CheckITR4Eligibility(tt.assessee, tt.residency, tt.income, tt.hasPresumptive, tt.foreignAssets, tt.isDirector, tt.ltcg)
			assert.Equal(t, tt.wantEligible, r.Eligible)
			if tt.wantReason != "" {
				assert.Contains(t, r.Reason, tt.wantReason)
			}
		})
	}
}

func TestCheckITR5Eligibility(t *testing.T) {
	tests := []struct {
		name         string
		assessee     AssesseeType
		wantEligible bool
	}{
		{"Firm", AssesseeFirm, true},
		{"LLP", AssesseeLLP, true},
		{"AOP", AssesseeAOP, true},
		{"BOI", AssesseeBOI, true},
		{"Individual", AssesseeIndividual, false},
		{"HUF", AssesseeHUF, false},
		{"Company", AssesseeCompany, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := CheckITR5Eligibility(tt.assessee)
			assert.Equal(t, tt.wantEligible, r.Eligible)
		})
	}
}

func TestCheckITR6Eligibility(t *testing.T) {
	tests := []struct {
		name         string
		assessee     AssesseeType
		claimsITR7   bool
		wantEligible bool
		wantReason   string
	}{
		{"Company_NoExemption", AssesseeCompany, false, true, ""},
		{"Company_ClaimsITR7", AssesseeCompany, true, false, "claiming exemption"},
		{"Individual", AssesseeIndividual, false, false, "companies only"},
		{"Firm", AssesseeFirm, false, false, "companies only"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := CheckITR6Eligibility(tt.assessee, tt.claimsITR7)
			assert.Equal(t, tt.wantEligible, r.Eligible)
			if tt.wantReason != "" {
				assert.Contains(t, r.Reason, tt.wantReason)
			}
		})
	}
}

func TestCheckITR7Eligibility(t *testing.T) {
	tests := []struct {
		name         string
		assessee     AssesseeType
		section      string
		wantEligible bool
		wantReason   string
	}{
		{"Trust_139_4A", AssesseeTrust, "139(4A)", true, ""},
		{"Trust_139_4B", AssesseeTrust, "139(4B)", true, ""},
		{"Company_139_4C", AssesseeCompany, "139(4C)", true, ""},
		{"AOP_139_4D", AssesseeAOP, "139(4D)", true, ""},
		{"Trust_BadSection", AssesseeTrust, "139(1)", false, "filing section must be"},
		{"Individual", AssesseeIndividual, "139(4A)", false, "trusts, institutions"},
		{"HUF", AssesseeHUF, "139(4A)", false, "trusts, institutions"},
		{"Firm", AssesseeFirm, "139(4A)", false, "trusts, institutions"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := CheckITR7Eligibility(tt.assessee, tt.section)
			assert.Equal(t, tt.wantEligible, r.Eligible)
			if tt.wantReason != "" {
				assert.Contains(t, r.Reason, tt.wantReason)
			}
		})
	}
}
