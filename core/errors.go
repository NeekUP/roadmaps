package core

type ErrorCode string

const (
	None                  ErrorCode = "NONE"
	BadPassword           ErrorCode = "INVALID_PASSWORD"
	BadEmail              ErrorCode = "INVALID_EMAIL"
	NotExistentEmail      ErrorCode = "NONEXISTENT_EMAIL"
	BadUserName           ErrorCode = "INVALID_USERNAME"
	NameAlreadyExists     ErrorCode = "USERNAME_ALREADY_EXISTS"
	EmailAlreadyExists    ErrorCode = "EMAIL_ALREADY_EXISTS"
	InternalError         ErrorCode = "INTERNAL_ERROR"
	AuthenticationError   ErrorCode = "AUTHENTICATION_ERROR"
	AuthenticationExpired ErrorCode = "AUTHENTICATION_EXPIRED"
	InvalidRequest        ErrorCode = "INVALID_REQUEST"
)

func (e ErrorCode) String() string {
	return string(e)
}

type AppError struct {
	message string
}

func (e *AppError) Error() string {
	return e.message
}

func NewError(msg ErrorCode) error {
	return &AppError{message: string(msg)}
}
