package core

import "encoding/json"

type ErrorCode string

const (
	None                  ErrorCode = "NONE"
	BadPassword           ErrorCode = "INVALID_PASSWORD"
	BadEmail              ErrorCode = "INVALID_EMAIL"
	NotExistentEmail      ErrorCode = "NONEXISTENT_EMAIL"
	BadUserName           ErrorCode = "INVALID_USERNAME"
	AlreadyExists         ErrorCode = "ALREADY_EXISTS"
	InternalError         ErrorCode = "INTERNAL_ERROR"
	AuthenticationError   ErrorCode = "AUTHENTICATION_ERROR"
	AuthenticationExpired ErrorCode = "AUTHENTICATION_EXPIRED"
	InvalidRequest        ErrorCode = "INVALID_REQUEST"
	InvalidUrl            ErrorCode = "INVALID_URL"
	InvalidISBN           ErrorCode = "INVALID_ISBN"
	InvalidTitle          ErrorCode = "INVALID_TITLE"
	InvalidProperties     ErrorCode = "INVALID_PROPS"
	InvalidSourceType     ErrorCode = "INVALID_SOURCE_TYPE"
	InaccessibleWebPage   ErrorCode = "INACCESSIBLE_WEBPAGE"
	InvalidFormat         ErrorCode = "INVALID_FORMAT"
	SourceNotFound        ErrorCode = "SOURCE_NOT_FOUND"
	InvalidCount          ErrorCode = "INVALID_COUNT"
	NotExists             ErrorCode = "NOT_EXISTS"
	
)

func (e ErrorCode) String() string {
	return string(e)
}

type AppError struct {
	Message    string            `json:"error"`
	Validation map[string]string `json:"validation"`
}

func (e *AppError) Error() string {
	msg, _ := json.Marshal(e)
	return string(msg)
}

func NewError(msg ErrorCode) *AppError {
	return &AppError{Message: string(msg)}
}

func ValidationError(validation map[string]string) *AppError {
	return &AppError{
		Message:    string(InvalidRequest),
		Validation: validation}
}
