package postgres

import (
	"context"
	"fmt"
	"time"

	"cm_open_api/internal/heatmap"

	"github.com/jackc/pgx/v4/pgxpool"
)

func GetHeatmapSourceRows(ctx context.Context, pool *pgxpool.Pool, from, toExclusive time.Time) ([]heatmap.SourceRow, error) {
	query := `
select
    date_trunc('month', t.event_start) as month,
    coalesce(m.context->'supplier', '') as service,
    coalesce(m.incident_id::text, m.id::text || ':' || m.chat_id::text) as incident_id,
    a.street_kladr,
    coalesce(a.city_kladr, ''),
    coalesce(nullif(coalesce(a.city_name, a.city_original), ''), '') as city_name,
    coalesce(nullif(coalesce(a.street_name, a.street_original), ''), '') as street_name,
    coalesce(nullif(coalesce(a.street_type, a.street_type_raw), ''), '') as street_type,
    coalesce(a.house_numbers, ''),
    coalesce(a.house_ranges, '')
from incident_address a
    join telegram_message_transcribes t on a.message_id = t.id and a.chat_id = t.chat_id
    join telegram_messages m on t.id = m.id and t.chat_id = m.chat_id
where t.event = 'shutdown'
  and t.event_start >= $1
  and t.event_start < $2
  and nullif(a.street_kladr, '') is not null
`

	rows, err := pool.Query(ctx, query, from, toExclusive)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %v", err)
	}
	defer rows.Close()

	var result []heatmap.SourceRow
	for rows.Next() {
		var r heatmap.SourceRow
		if err := rows.Scan(
			&r.Month, &r.Service, &r.IncidentID, &r.StreetKladr,
			&r.CityKladr, &r.CityName, &r.StreetName, &r.StreetType,
			&r.HouseNumbers, &r.HouseRanges,
		); err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		result = append(result, r)
	}
	return result, nil
}
