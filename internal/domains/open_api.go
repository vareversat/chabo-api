package domains

type OpenAPISelfLink struct {
	Self OpenAPILink `json:"self"`
}

type OpenAPIPreviousLink struct {
	Self OpenAPILink `json:"previous"`
}

type OpenAPINextLink struct {
	Self OpenAPILink `json:"next"`
}

type OpenAPILink struct {
	Link string `json:"href"`
}
