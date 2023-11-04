package errors

type CustomError interface {
	GetStatusCode() int
	GetErrorMessage() string
}
