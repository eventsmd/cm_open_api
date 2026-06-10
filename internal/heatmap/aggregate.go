// Package heatmap turns flat incident rows into per-month/service/street
// aggregates with distinct-incident counts per house.
package heatmap

import (
	"sort"
	"time"

	"cm_open_api/internal/houses"
	"cm_open_api/internal/models"
)

type SourceRow struct {
	Month        time.Time
	Service      string
	IncidentID   string
	StreetKladr  string
	CityKladr    string
	CityName     string
	StreetName   string
	StreetType   string
	HouseNumbers string
	HouseRanges  string
}

type groupKey struct {
	month   string
	service string
	street  string
}

type group struct {
	row                 models.HeatmapRow
	houseIncidents      map[string]map[string]struct{} // house -> set of incident ids
	streetIncidents     map[string]struct{}            // candidate incident ids from rows without resolvable houses
	incidentsWithHouses map[string]struct{}            // incident ids with at least one resolvable house in the group
}

// Aggregate groups rows by (month, service, street) counting distinct
// incidents per house. Returns deterministically sorted rows and the total
// number of distinct incidents across the period.
func Aggregate(rows []SourceRow) ([]models.HeatmapRow, int) {
	groups := map[groupKey]*group{}
	allIncidents := map[string]struct{}{}

	for _, r := range rows {
		allIncidents[r.IncidentID] = struct{}{}
		key := groupKey{r.Month.Format("2006-01"), r.Service, r.StreetKladr}
		g, ok := groups[key]
		if !ok {
			g = &group{
				row: models.HeatmapRow{
					Month: key.month, Service: key.service, StreetKladr: key.street,
					CityKladr: r.CityKladr, CityName: r.CityName,
					StreetName: r.StreetName, StreetType: r.StreetType,
				},
				houseIncidents:      map[string]map[string]struct{}{},
				streetIncidents:     map[string]struct{}{},
				incidentsWithHouses: map[string]struct{}{},
			}
			groups[key] = g
		}

		hs, _ := houses.Expand(r.HouseNumbers, r.HouseRanges)
		if len(hs) == 0 {
			g.streetIncidents[r.IncidentID] = struct{}{}
			continue
		}
		g.incidentsWithHouses[r.IncidentID] = struct{}{}
		for _, h := range hs {
			if g.houseIncidents[h] == nil {
				g.houseIncidents[h] = map[string]struct{}{}
			}
			g.houseIncidents[h][r.IncidentID] = struct{}{}
		}
	}

	result := make([]models.HeatmapRow, 0, len(groups))
	for _, g := range groups {
		if len(g.houseIncidents) > 0 {
			g.row.HouseIncidents = make(map[string]int, len(g.houseIncidents))
			for h, ids := range g.houseIncidents {
				g.row.HouseIncidents[h] = len(ids)
			}
		}
		// Final per-incident classification: an incident with at least one
		// resolvable house in the group must not be counted at street level.
		for id := range g.streetIncidents {
			if _, ok := g.incidentsWithHouses[id]; !ok {
				g.row.StreetIncidents++
			}
		}
		result = append(result, g.row)
	}
	sort.Slice(result, func(i, j int) bool {
		a, b := result[i], result[j]
		if a.Month != b.Month {
			return a.Month < b.Month
		}
		if a.Service != b.Service {
			return a.Service < b.Service
		}
		return a.StreetKladr < b.StreetKladr
	})
	return result, len(allIncidents)
}
