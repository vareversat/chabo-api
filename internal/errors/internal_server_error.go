package errors

import "fmt"

type internalServerError struct {
	Info string
}

func NewInternalServerError() CustomError {
	return &internalServerError{}
}

func (e *internalServerError) GetStatusCode() int {
	return 404
}

func (e *internalServerError) GetErrorMessage() string {
	return fmt.Sprintf("Internal server error. Information : %s", e.Info)
}
