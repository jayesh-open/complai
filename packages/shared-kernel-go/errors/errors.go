package errors

import "fmt"

type NotFoundError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func NewNotFound(code, message string, details ...any) *NotFoundError {
	var d any
	if len(details) > 0 {
		d = details[0]
	}
	return &NotFoundError{Code: code, Message: message, Details: d}
}

type ConflictError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

func (e *ConflictError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func NewConflict(code, message string, details ...any) *ConflictError {
	var d any
	if len(details) > 0 {
		d = details[0]
	}
	return &ConflictError{Code: code, Message: message, Details: d}
}

type ValidationError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func NewValidation(code, message string, details ...any) *ValidationError {
	var d any
	if len(details) > 0 {
		d = details[0]
	}
	return &ValidationError{Code: code, Message: message, Details: d}
}

type AuthorizationError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

func (e *AuthorizationError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func NewAuthorization(code, message string, details ...any) *AuthorizationError {
	var d any
	if len(details) > 0 {
		d = details[0]
	}
	return &AuthorizationError{Code: code, Message: message, Details: d}
}

type ProviderError struct {
	Code     string `json:"code"`
	Message  string `json:"message"`
	Details  any    `json:"details,omitempty"`
	Provider string `json:"provider"`
}

func (e *ProviderError) Error() string {
	return fmt.Sprintf("[%s] %s (provider: %s)", e.Code, e.Message, e.Provider)
}

func NewProvider(code, message, provider string, details ...any) *ProviderError {
	var d any
	if len(details) > 0 {
		d = details[0]
	}
	return &ProviderError{Code: code, Message: message, Provider: provider, Details: d}
}

type InternalError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

func (e *InternalError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func NewInternal(code, message string, details ...any) *InternalError {
	var d any
	if len(details) > 0 {
		d = details[0]
	}
	return &InternalError{Code: code, Message: message, Details: d}
}

func IsNotFound(err error) bool {
	_, ok := err.(*NotFoundError)
	return ok
}

func IsConflict(err error) bool {
	_, ok := err.(*ConflictError)
	return ok
}

func IsValidation(err error) bool {
	_, ok := err.(*ValidationError)
	return ok
}

func IsAuthorization(err error) bool {
	_, ok := err.(*AuthorizationError)
	return ok
}

func IsProvider(err error) bool {
	_, ok := err.(*ProviderError)
	return ok
}

func IsInternal(err error) bool {
	_, ok := err.(*InternalError)
	return ok
}
