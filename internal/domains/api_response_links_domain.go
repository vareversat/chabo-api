package domains

type APIResponseSelfLink struct {
	Self APIResponseLink `json:"self"`
}

type APIResponsePreviousLink struct {
	Self APIResponseLink `json:"previous"`
}

type APIResponseNextLink struct {
	Self APIResponseLink `json:"next"`
}

type APIResponseLink struct {
	Link string `json:"href"`
}
