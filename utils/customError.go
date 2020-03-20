package utils

type customError struct {
	message string
}

func (e *customError) Error() string {
	return e.message
}

// NewCustomError …
func NewCustomError(message string) error {
	return &customError{message}
}
