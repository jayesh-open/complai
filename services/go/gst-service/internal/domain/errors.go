package domain

import "errors"

var (
	ErrNotFound          = errors.New("not found")
	ErrInvalidState      = errors.New("invalid filing state for this operation")
	ErrDuplicateEntry    = errors.New("duplicate sales register entry")
	ErrValidationFailed  = errors.New("validation errors exist")
	ErrNotApproved       = errors.New("filing not approved by checker")
	ErrStepUpRequired    = errors.New("step-up authentication required")
	ErrInvalidGSTIN      = errors.New("invalid GSTIN format")
	ErrInvalidPeriod     = errors.New("invalid return period format")
	ErrAlreadyFiled      = errors.New("return already filed for this period")
)
