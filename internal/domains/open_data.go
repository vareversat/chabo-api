package domains

import "time"

type OpenDataAPIResponse struct {
	Hits       int                           `json:"nhits"`
	Parameters OpenDataAPIResponseParameters `json:"parameters"`
	Records    []OpenDataAPIResponseForecast `json:"records"`
}

type OpenDataAPIResponseParameters struct {
	Dataset  string   `json:"dataset"`
	Row      int      `json:"rows"`
	Start    int      `json:"start"`
	Sort     []string `json:"sort"`
	Format   string   `json:"format"`
	Timezone string   `json:"timezone"`
}

type OpenDataAPIResponseForecast struct {
	DatasetID       string                           `json:"datasetid"`
	RecordID        string                           `json:"recordid"`
	Fields          OpenDataAPIResponseForecastField `json:"fields"`
	RecordTimestamp time.Time                        `json:"record_timestamp"`
}

type OpenDataAPIResponseForecastField struct {
	ClosingDate  string `json:"date_passage"`
	ClosingTime  string `json:"fermeture_a_la_circulation"`
	OpeningTime  string `json:"re_ouverture_a_la_circulation"`
	TotalClosing string `json:"fermeture_totale"`
	Boat         string `json:"bateau"`
	ClosingType  string `json:"type_de_fermeture"`
}
