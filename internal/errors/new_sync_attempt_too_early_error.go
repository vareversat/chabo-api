package errors

type syncAttemptTooEarlyError struct {
}

func NewSyncAttemptTooEarlyError() CustomError {
	return &syncAttemptTooEarlyError{}
}

func (e *syncAttemptTooEarlyError) GetStatusCode() int {
	return 429
}

func (e *syncAttemptTooEarlyError) GetErrorMessage() string {
	return "The last sync is too recent. Please try again later"
}
