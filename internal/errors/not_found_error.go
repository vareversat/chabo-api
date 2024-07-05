package errors

type notFoundError struct {
	Info string
}

func NewNotFoundError(info string) CustomError {
	return &notFoundError{Info: info}
}

func (e *notFoundError) GetStatusCode() int {
	return 404
}

func (e *notFoundError) GetErrorMessage() string {
	return e.Info
}
