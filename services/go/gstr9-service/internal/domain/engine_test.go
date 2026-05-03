package domain

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestCheckThreshold_BelowTwoCrore(t *testing.T) {
	turnover := decimal.NewFromInt(15000000)
	r := CheckThreshold(turnover)
	assert.False(t, r.GSTR9Mandatory)
	assert.False(t, r.GSTR9CRequired)
	assert.Contains(t, r.Reason, "optional")
}

func TestCheckThreshold_ExactlyTwoCrore(t *testing.T) {
	turnover := decimal.NewFromInt(20000000)
	r := CheckThreshold(turnover)
	assert.False(t, r.GSTR9Mandatory, "exactly ₹2Cr should not be mandatory (threshold is >)")
	assert.False(t, r.GSTR9CRequired)
	assert.Contains(t, r.Reason, "optional")
}

func TestCheckThreshold_AboveTwoCrore(t *testing.T) {
	turnover := decimal.NewFromInt(20000001)
	r := CheckThreshold(turnover)
	assert.True(t, r.GSTR9Mandatory)
	assert.False(t, r.GSTR9CRequired)
	assert.Contains(t, r.Reason, "exceeds")
}

func TestCheckThreshold_BelowFiveCrore(t *testing.T) {
	turnover := decimal.NewFromInt(40000000)
	r := CheckThreshold(turnover)
	assert.True(t, r.GSTR9Mandatory)
	assert.False(t, r.GSTR9CRequired)
}

func TestCheckThreshold_ExactlyFiveCrore(t *testing.T) {
	turnover := decimal.NewFromInt(50000000)
	r := CheckThreshold(turnover)
	assert.True(t, r.GSTR9Mandatory)
	assert.False(t, r.GSTR9CRequired, "exactly ₹5Cr should not require 9C (threshold is >)")
}

func TestCheckThreshold_AboveFiveCrore(t *testing.T) {
	turnover := decimal.NewFromInt(50000001)
	r := CheckThreshold(turnover)
	assert.True(t, r.GSTR9Mandatory)
	assert.True(t, r.GSTR9CRequired)
	assert.Contains(t, r.Reason, "GSTR-9C reconciliation mandatory")
}

func TestCheckThreshold_ZeroTurnover(t *testing.T) {
	r := CheckThreshold(decimal.Zero)
	assert.False(t, r.GSTR9Mandatory)
	assert.False(t, r.GSTR9CRequired)
	assert.True(t, r.AggregateTurnover.IsZero())
}

func TestCheckThreshold_LargeTurnover(t *testing.T) {
	turnover := decimal.NewFromInt(5000000000)
	r := CheckThreshold(turnover)
	assert.True(t, r.GSTR9Mandatory)
	assert.True(t, r.GSTR9CRequired)
}
