package Utils

var (
	ErrInvalidInput        = NewError("Invalid input")
	ErrUnauthorized        = NewError("Unauthorized")
	ErrNotFound            = NewError("Not found")
	ErrInternalServerError = NewError("Internal server error")
	// Add more error messages here...
)

type Error struct {
	message string
}

func NewError(message string) *Error {
	return &Error{message: message}
}

func (e *Error) Error() string {
	return e.message
}
