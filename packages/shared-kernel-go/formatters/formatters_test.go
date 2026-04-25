package formatters_test

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/complai/complai/packages/shared-kernel-go/formatters"
)

func TestFormatINR(t *testing.T) {
	tests := []struct {
		name     string
		amount   string
		expected string
	}{
		{"zero", "0", "₹0.00"},
		{"small amount", "500", "₹500.00"},
		{"thousands", "1234.56", "₹1,234.56"},
		{"lakhs", "100000", "₹1,00,000.00"},
		{"ten lakhs", "1000000", "₹10,00,000.00"},
		{"crore", "10000000", "₹1,00,00,000.00"},
		{"large crore", "1234567890.12", "₹1,23,45,67,890.12"},
		{"negative", "-5000", "-₹5,000.00"},
		{"single digit", "5", "₹5.00"},
		{"two digits", "42", "₹42.00"},
		{"three digits", "999", "₹999.00"},
		{"four digits", "1000", "₹1,000.00"},
		{"five digits", "12345", "₹12,345.00"},
		{"six digits", "123456", "₹1,23,456.00"},
		{"seven digits", "1234567", "₹12,34,567.00"},
		{"fractional paise", "0.01", "₹0.01"},
		{"99 paise", "0.99", "₹0.99"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			amount, _ := decimal.NewFromString(tt.amount)
			assert.Equal(t, tt.expected, formatters.FormatINR(amount))
		})
	}
}

func TestFormatDate(t *testing.T) {
	dt := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	assert.Equal(t, "15/01/2024", formatters.FormatDate(dt))

	dt2 := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)
	assert.Equal(t, "31/12/2024", formatters.FormatDate(dt2))
}

func TestParseDate(t *testing.T) {
	dt, err := formatters.ParseDate("15/01/2024")
	require.NoError(t, err)
	assert.Equal(t, 15, dt.Day())
	assert.Equal(t, time.January, dt.Month())
	assert.Equal(t, 2024, dt.Year())

	_, err = formatters.ParseDate("2024-01-15")
	assert.Error(t, err)

	_, err = formatters.ParseDate("invalid")
	assert.Error(t, err)
}

func TestValidateGSTIN(t *testing.T) {
	tests := []struct {
		name  string
		gstin string
		valid bool
	}{
		{"valid Maharashtra", "27AAPFU0939F1ZV", true},
		{"valid Delhi", "07AAACH7409R1ZZ", true},
		{"too short", "27AAPFU0939F1Z", false},
		{"too long", "27AAPFU0939F1ZVV", false},
		{"lowercase", "27aapfu0939f1zv", false},
		{"invalid state code", "XX AAPFU0939F1ZV", false},
		{"empty", "", false},
		{"missing Z", "27AAPFU0939F1AV", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.valid, formatters.ValidateGSTIN(tt.gstin))
		})
	}
}

func TestValidatePAN(t *testing.T) {
	tests := []struct {
		name  string
		pan   string
		valid bool
	}{
		{"valid individual", "ABCDE1234F", true},
		{"valid company", "AABCU9603R", true},
		{"too short", "ABCDE123", false},
		{"too long", "ABCDE12345F", false},
		{"lowercase", "abcde1234f", false},
		{"all numbers", "1234567890", false},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.valid, formatters.ValidatePAN(tt.pan))
		})
	}
}

func TestValidateTAN(t *testing.T) {
	tests := []struct {
		name  string
		tan   string
		valid bool
	}{
		{"valid", "BLRA12345F", true},
		{"valid DEL", "DELB12345G", true},
		{"too short", "BLRA1234", false},
		{"too long", "BLRA123456F", false},
		{"lowercase", "blra12345f", false},
		{"wrong format", "1234ABCDE5", false},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.valid, formatters.ValidateTAN(tt.tan))
		})
	}
}
