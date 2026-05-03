package domain

import "errors"

var (
	ErrFilingNotFound   = errors.New("gstr9 filing not found")
	ErrDuplicateFiling  = errors.New("gstr9 filing already exists for this GSTIN and FY")
	ErrInvalidStatus    = errors.New("invalid filing status transition")
	ErrNoMonthlyData    = errors.New("no monthly GSTR-1/3B data found for the financial year")
	ErrIncompleteMonths = errors.New("less than 12 months of data found")

	ErrGSTR9CNotFound       = errors.New("gstr9c filing not found")
	ErrGSTR9CDuplicate      = errors.New("gstr9c reconciliation already exists for this GSTR-9 filing")
	ErrMismatchNotFound     = errors.New("gstr9c mismatch not found")
	ErrGSTR9CNotReconciled  = errors.New("gstr9c must be reconciled before certification")
	ErrGSTR9CAlreadyCertified = errors.New("gstr9c is already certified")
	ErrUnresolvedMismatches = errors.New("unresolved ERROR-severity mismatches block certification")
)
