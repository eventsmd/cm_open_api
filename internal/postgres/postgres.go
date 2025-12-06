package postgres

import (
	"cm_open_api/internal/models"
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

func GetOutages(ctx context.Context, pool *pgxpool.Pool) ([]models.Outage, error) {

	query := `
select
    a.chat_id::text || ':' || a.message_id::text || ':' || a.id::text as id,
    m.incident_id,
    m.context->'supplier' as service,
    t.organization,
    t.description,
    t.event,
    t.event_start,
    t.event_stop,
    NULLIF(a.region_kladr, ''),
    NULLIF(a.region_type, ''),
    NULLIF(a. region_name, ''),
    NULLIF(a.city_kladr, ''),
    NULLIF(COALESCE(city_name, city_original), '') as city_name,
    NULLIF(city_type, ''),
    NULLIF(street_kladr, ''),
    NULLIF(COALESCE(street_name, street_original), '') as street_name,
    NULLIF(COALESCE(street_type, street_type_raw), '') as street_type,
    NULLIF(house_numbers, ''),
    NULLIF(house_ranges, '')
from incident_address a
    join telegram_message_transcribes t on a.message_id = t.id and a.chat_id = t.chat_id
    join telegram_messages m on t.id = m.id and t.chat_id = m.chat_id
WHERE event = 'shutdown'
  AND event_start >= date_trunc('day', current_date)
ORDER BY event_start
`

	rows, err := pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %v", err)
	}
	defer rows.Close()

	var outages []models.Outage
	for rows.Next() {
		var outage models.Outage
		var regionKladr, regionType, regionName, streetKladr, houseNumbers, houseRanges *string
		var cityKladr, cityName, cityType, streetName, streetType *string
		var eventStart, eventStop *time.Time

		err := rows.Scan(
			&outage.MessageID, &outage.IncidentID, &outage.Service, &outage.Organization,
			&outage.ShortDescription, &outage.Event, &eventStart, &eventStop,
			&regionKladr, &regionType, &regionName,
			&cityKladr, &cityName, &cityType,
			&streetKladr, &streetName, &streetType,
			&houseNumbers, &houseRanges,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		if eventStart != nil {
			outage.EventStart = *eventStart
		}
		if eventStop != nil {
			outage.EventStop = eventStop
		}
		if regionKladr != nil {
			outage.RegionKladr = regionKladr
		}
		if regionType != nil {
			outage.RegionType = regionType
		}
		if regionName != nil {
			outage.RegionName = regionName
		}
		if cityKladr != nil {
			outage.CityKladr = cityKladr
		}
		if cityName != nil {
			outage.CityName = cityName
		}
		if cityType != nil {
			outage.CityType = cityType
		}
		if streetKladr != nil {
			outage.StreetKladr = streetKladr
		}
		if streetName != nil {
			outage.StreetName = streetName
		}
		if streetType != nil {
			outage.StreetType = streetType
		}
		if houseNumbers != nil {
			outage.HouseNumbers = strings.Split(*houseNumbers, ",")
		}
		if houseRanges != nil {
			ranges := strings.Split(*houseRanges, ";")
			for i, r := range ranges {
				ranges[i] = strings.TrimSuffix(r, "-")
			}
			outage.HouseRanges = ranges
		}

		outages = append(outages, outage)
	}

	return outages, nil
}

func GetSource(ctx context.Context, pool *pgxpool.Pool, messageId, chatId string) (*models.SourceResponse, error) {

	var sourceResponse models.SourceResponse

	query := `
        select 
            date,
            text,
            chat_id::text as chat_id,
            from_id::text as from_id,
            from_name
        from telegram_messages
        where chat_id = $1 and id = $2`

	rows, err := pool.Query(ctx, query, chatId, messageId)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %v", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, fmt.Errorf("no source found for chat_id=%s, message_id=%s", chatId, messageId)
	}

	var (
		createdAt  *time.Time
		rawTextNS  sql.NullString
		chatIDNS   sql.NullString
		fromIDNS   sql.NullString
		fromNameNS sql.NullString
	)

	if err := rows.Scan(&createdAt, &rawTextNS, &chatIDNS, &fromIDNS, &fromNameNS); err != nil {
		return nil, fmt.Errorf("error scanning row: %v", err)
	}

	if createdAt != nil {
		sourceResponse.CreatedAt = createdAt.Format(time.RFC3339)
	}

	if rawTextNS.Valid {
		sourceResponse.RawMessage = rawTextNS.String
	}

	channel := fmt.Sprintf("https://t.me/c/%s", chatIDNS.String)
	sourceResponse.Source.Channel = &channel

	if fromIDNS.Valid {
		senderURI := "tg://user?id=" + fromIDNS.String
		sourceResponse.Source.SenderURI = &senderURI
	}
	if fromNameNS.Valid {
		senderName := fromNameNS.String
		sourceResponse.Source.SenderName = &senderName
	}
	if chatIDNS.Valid {
		sourceURI := fmt.Sprintf("https://t.me/c/%s/%s", chatIDNS.String, messageId)
		sourceResponse.Source.SourceURI = &sourceURI
	}

	return &sourceResponse, nil
}
