package models

type ForecastsResponse struct {
	Hits      int           `json:"hits"`
	Limit     int           `json:"limit"`
	Offset    int           `json:"offset"`
	Timezone  string        `json:"timezone" example:"UTC"`
	Links     []interface{} `json:"_links,omitempty"`
	Forecasts []Forecast    `json:"forecasts"`
}

type ForecastResponse struct {
	Timezone string   `json:"timezone" example:"UTC"`
	Forecast Forecast `json:"forecast"`
}

type MongoResponse struct {
	Results []Forecast           `json:"results" bson:"results"`
	Count   []MongoCountResponse `json:"count" bson:"count"`
}

type MongoCountResponse struct {
	ItemCount int `json:"itemCount" bson:"itemCount"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"error in params"`
}

type OKResponse struct {
	Message string `json:"message" example:"ok"`
}
