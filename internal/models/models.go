package models

import "time"

type Outage struct {
	MessageID        string     `json:"message_id"`
	IncidentID       string     `json:"incident_id"`
	Service          string     `json:"service"`
	Organization     string     `json:"organization"`
	ShortDescription string     `json:"short_description"`
	Event            string     `json:"event"`
	EventStart       time.Time  `json:"-"`
	EventStartStr    string     `json:"event_start"`
	EventStop        *time.Time `json:"-"`
	EventStopStr     *string    `json:"event_stop,omitempty"`

	RegionKladr *string `json:"region_kladr,omitempty"`
	RegionType  *string `json:"region_type,omitempty"`
	RegionName  *string `json:"region_name,omitempty"`

	CityKladr *string `json:"city_kladr,omitempty"`
	CityName  *string `json:"city_name"`
	CityType  *string `json:"city_type,omitempty"`

	StreetKladr  *string  `json:"street_kladr,omitempty"`
	StreetName   *string  `json:"street_name"`
	StreetType   *string  `json:"street_type"`
	HouseNumbers []string `json:"house_numbers,omitempty"`
	HouseRanges  []string `json:"house_ranges,omitempty"`
}

type Source struct {
	Channel    *string `json:"channel"`
	SenderURI  *string `json:"sender_uri"`
	SenderName *string `json:"sender_name"`
	SourceURI  *string `json:"source_uri"`
}

type SourceResponse struct {
	CreatedAt  string `json:"created_at"`
	RawMessage string `json:"raw_message"`
	Source     Source `json:"source"`
}
