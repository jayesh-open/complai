package formatters

import "regexp"

// GSTIN: 2-digit state code + 10-char PAN + 1-digit entity number + Z + 1 check char
// Example: 27AAPFU0939F1ZV
var gstinRegex = regexp.MustCompile(`^[0-9]{2}[A-Z]{5}[0-9]{4}[A-Z]{1}[1-9A-Z]{1}Z[0-9A-Z]{1}$`)

// PAN: 5 letters + 4 digits + 1 letter
// Example: ABCDE1234F
var panRegex = regexp.MustCompile(`^[A-Z]{5}[0-9]{4}[A-Z]{1}$`)

// TAN: 4 letters + 5 digits + 1 letter
// Example: BLRA12345F
var tanRegex = regexp.MustCompile(`^[A-Z]{4}[0-9]{5}[A-Z]{1}$`)

func ValidateGSTIN(gstin string) bool {
	if len(gstin) != 15 {
		return false
	}
	return gstinRegex.MatchString(gstin)
}

func ValidatePAN(pan string) bool {
	if len(pan) != 10 {
		return false
	}
	return panRegex.MatchString(pan)
}

func ValidateTAN(tan string) bool {
	if len(tan) != 10 {
		return false
	}
	return tanRegex.MatchString(tan)
}
