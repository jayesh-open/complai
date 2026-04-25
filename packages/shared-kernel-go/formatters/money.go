package formatters

import (
	"strings"

	"github.com/shopspring/decimal"
)

// FormatINR formats a decimal amount as Indian Rupees with the Indian numbering system.
// Example: 1234567.89 -> "₹12,34,567.89"
func FormatINR(amount decimal.Decimal) string {
	negative := amount.IsNegative()
	if negative {
		amount = amount.Neg()
	}

	str := amount.StringFixed(2)

	parts := strings.SplitN(str, ".", 2)
	intPart := parts[0]
	decPart := parts[1]

	formatted := formatIndianGrouping(intPart)

	result := "₹" + formatted + "." + decPart
	if negative {
		result = "-" + result
	}
	return result
}

// formatIndianGrouping applies the Indian numbering convention:
// last 3 digits as one group, then groups of 2 from right to left.
func formatIndianGrouping(s string) string {
	n := len(s)
	if n <= 3 {
		return s
	}

	// Last 3 digits
	lastThree := s[n-3:]
	remaining := s[:n-3]

	var groups []string
	for len(remaining) > 2 {
		groups = append([]string{remaining[len(remaining)-2:]}, groups...)
		remaining = remaining[:len(remaining)-2]
	}
	if len(remaining) > 0 {
		groups = append([]string{remaining}, groups...)
	}
	groups = append(groups, lastThree)

	return strings.Join(groups, ",")
}
