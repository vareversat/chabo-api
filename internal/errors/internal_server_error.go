package errors

import "fmt"

type internalServerError struct {
	Info string
}

func NewInternalServerError(info string) CustomError {
	return &internalServerError{Info: info}
}

func (e *internalServerError) GetStatusCode() int {
	return 500
}

func (e *internalServerError) GetErrorMessage() string {
	return fmt.Sprintf("Internal server error. Information : %s", e.Info)
}
