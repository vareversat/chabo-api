package domains

type ErrorResponse struct {
	Error string `json:"error" example:"error in params"`
}

type OKResponse struct {
	Message string `json:"message" example:"ok"`
}
