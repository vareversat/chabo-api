package domains

type APIErrorResponse struct {
	Error string `json:"error" example:"error in params"`
}

type APIOKResponse struct {
	Message string `json:"message" example:"ok"`
}
