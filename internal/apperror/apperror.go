package apperror

import "fmt"

type ErrorCode string

const (
	CodeInternal     ErrorCode = "INTERNAL_SERVER_ERROR"
	CodeInvalidInput ErrorCode = "INVALID_INPUT"
	CodeUnauthorized ErrorCode = "UNAUTHORIZED"
	CodeForbidden    ErrorCode = "FORBIDDEN"
	CodeNotFound     ErrorCode = "NOT_FOUND"
	CodeConflict     ErrorCode = "CONFLICT"
	CodeEmailTaken   ErrorCode = "EMAIL_TAKEN"
)

type AppError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Details any       `json:"-"`
	Cause   error     `json:"-"`
}

func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

func New(code ErrorCode, message string) *AppError {
	return &AppError{Code: code, Message: message}
}

func Wrap(err error, code ErrorCode, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Cause:   err,
	}
}

func (e *AppError) WithMessagef(format string, args ...any) *AppError {
	return &AppError{
		Code:    e.Code,
		Message: fmt.Sprintf(format, args...),
		Details: e.Details,
		Cause:   e.Cause,
	}
}

// WithDetails attaches arbitrary details to the error
func (e *AppError) WithDetails(details any) *AppError {
	return &AppError{
		Code:    e.Code,
		Message: e.Message,
		Details: details,
		Cause:   e.Cause,
	}
}
