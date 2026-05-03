package domain

import "errors"

var (
	ErrFilingNotFound   = errors.New("gstr9 filing not found")
	ErrDuplicateFiling  = errors.New("gstr9 filing already exists for this GSTIN and FY")
	ErrInvalidStatus    = errors.New("invalid filing status transition")
	ErrNoMonthlyData    = errors.New("no monthly GSTR-1/3B data found for the financial year")
	ErrIncompleteMonths = errors.New("less than 12 months of data found")
)
