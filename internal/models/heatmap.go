package models

type HeatmapAggregatesResponse struct {
	From           string       `json:"from"`
	To             string       `json:"to"`
	TotalIncidents int          `json:"total_incidents"`
	Rows           []HeatmapRow `json:"rows"`
}

type HeatmapRow struct {
	Month           string         `json:"month"`   // "2026-01"
	Service         string         `json:"service"` // water | electricity
	StreetKladr     string         `json:"street_kladr"`
	CityKladr       string         `json:"city_kladr,omitempty"`
	CityName        string         `json:"city_name,omitempty"`
	StreetName      string         `json:"street_name,omitempty"`
	StreetType      string         `json:"street_type,omitempty"`
	HouseIncidents  map[string]int `json:"house_incidents,omitempty"`
	StreetIncidents int            `json:"street_incidents,omitempty"`
}
