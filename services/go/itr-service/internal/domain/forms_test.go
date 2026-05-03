package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckITR1Eligibility_Eligible(t *testing.T) {
	result := CheckITR1Eligibility(
		AssesseeIndividual, Resident, d(3000000), 1, false, d(0), false, false, false, false, d(0),
	)
	assert.True(t, result.Eligible)
}

func TestCheckITR1Eligibility_HUFRejected(t *testing.T) {
	result := CheckITR1Eligibility(
		AssesseeHUF, Resident, d(3000000), 0, false, d(0), false, false, false, false, d(0),
	)
	assert.False(t, result.Eligible)
	assert.Contains(t, result.Reason, "individuals")
}

func TestCheckITR1Eligibility_NRIRejected(t *testing.T) {
	result := CheckITR1Eligibility(
		AssesseeIndividual, NonResident, d(3000000), 0, false, d(0), false, false, false, false, d(0),
	)
	assert.False(t, result.Eligible)
	assert.Contains(t, result.Reason, "resident")
}

func TestCheckITR1Eligibility_IncomeExceeds50L(t *testing.T) {
	result := CheckITR1Eligibility(
		AssesseeIndividual, Resident, d(6000000), 0, false, d(0), false, false, false, false, d(0),
	)
	assert.False(t, result.Eligible)
	assert.Contains(t, result.Reason, "exceeds")
}

func TestCheckITR1Eligibility_TwoPropertiesAllowed(t *testing.T) {
	result := CheckITR1Eligibility(
		AssesseeIndividual, Resident, d(3000000), 2, false, d(0), false, false, false, false, d(0),
	)
	assert.True(t, result.Eligible)
}

func TestCheckITR1Eligibility_ThreePropertiesRejected(t *testing.T) {
	result := CheckITR1Eligibility(
		AssesseeIndividual, Resident, d(3000000), 3, false, d(0), false, false, false, false, d(0),
	)
	assert.False(t, result.Eligible)
	assert.Contains(t, result.Reason, "house properties")
}

func TestCheckITR1Eligibility_LTCG112A_125K_Allowed(t *testing.T) {
	result := CheckITR1Eligibility(
		AssesseeIndividual, Resident, d(3000000), 0, false, d(125000), false, false, false, false, d(0),
	)
	assert.True(t, result.Eligible)
}

func TestCheckITR1Eligibility_LTCG112A_Over125K_Rejected(t *testing.T) {
	result := CheckITR1Eligibility(
		AssesseeIndividual, Resident, d(3000000), 0, false, d(200000), false, false, false, false, d(0),
	)
	assert.False(t, result.Eligible)
	assert.Contains(t, result.Reason, "112A")
}

func TestCheckITR1Eligibility_CapitalGainsOver112A_Rejected(t *testing.T) {
	result := CheckITR1Eligibility(
		AssesseeIndividual, Resident, d(3000000), 0, true, d(0), false, false, false, false, d(0),
	)
	assert.False(t, result.Eligible)
	assert.Contains(t, result.Reason, "capital gains")
}

func TestCheckITR1Eligibility_Business_Rejected(t *testing.T) {
	result := CheckITR1Eligibility(
		AssesseeIndividual, Resident, d(3000000), 0, false, d(0), true, false, false, false, d(0),
	)
	assert.False(t, result.Eligible)
	assert.Contains(t, result.Reason, "business")
}

func TestCheckITR1Eligibility_ForeignAssets_Rejected(t *testing.T) {
	result := CheckITR1Eligibility(
		AssesseeIndividual, Resident, d(3000000), 0, false, d(0), false, true, false, false, d(0),
	)
	assert.False(t, result.Eligible)
	assert.Contains(t, result.Reason, "foreign")
}

func TestCheckITR1Eligibility_UnlistedEquity_Rejected(t *testing.T) {
	result := CheckITR1Eligibility(
		AssesseeIndividual, Resident, d(3000000), 0, false, d(0), false, false, true, false, d(0),
	)
	assert.False(t, result.Eligible)
	assert.Contains(t, result.Reason, "unlisted")
}

func TestCheckITR1Eligibility_Director_Rejected(t *testing.T) {
	result := CheckITR1Eligibility(
		AssesseeIndividual, Resident, d(3000000), 0, false, d(0), false, false, false, true, d(0),
	)
	assert.False(t, result.Eligible)
	assert.Contains(t, result.Reason, "director")
}

func TestCheckITR1Eligibility_AgriIncome_Over5K_Rejected(t *testing.T) {
	result := CheckITR1Eligibility(
		AssesseeIndividual, Resident, d(3000000), 0, false, d(0), false, false, false, false, d(10000),
	)
	assert.False(t, result.Eligible)
	assert.Contains(t, result.Reason, "agricultural")
}

func TestCheckITR2Eligibility_Individual_NoBusiness(t *testing.T) {
	result := CheckITR2Eligibility(AssesseeIndividual, false)
	assert.True(t, result.Eligible)
}

func TestCheckITR2Eligibility_HUF_NoBusiness(t *testing.T) {
	result := CheckITR2Eligibility(AssesseeHUF, false)
	assert.True(t, result.Eligible)
}

func TestCheckITR2Eligibility_BusinessRejected(t *testing.T) {
	result := CheckITR2Eligibility(AssesseeIndividual, true)
	assert.False(t, result.Eligible)
	assert.Contains(t, result.Reason, "ITR-3")
}

func TestCheckITR3Eligibility_Individual(t *testing.T) {
	result := CheckITR3Eligibility(AssesseeIndividual)
	assert.True(t, result.Eligible)
}

func TestCheckITR3Eligibility_HUF(t *testing.T) {
	result := CheckITR3Eligibility(AssesseeHUF)
	assert.True(t, result.Eligible)
}
