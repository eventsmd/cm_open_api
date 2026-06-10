package handlers

import (
	"cm_open_api/internal/config"
	"cm_open_api/internal/heatmap"
	"cm_open_api/internal/models"
	"cm_open_api/internal/postgres"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"
)

const aggregatesCacheTTL = time.Hour

type aggregatesCache struct {
	mu   sync.Mutex
	key  string
	body []byte
	at   time.Time
}

func getHeatmapAggregates(cfg *config.Config) http.HandlerFunc {
	cache := &aggregatesCache{}
	return func(w http.ResponseWriter, r *http.Request) {
		from, toExcl, err := parseMonthRange(r.URL.Query().Get("from"), r.URL.Query().Get("to"), time.Now())
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fromStr := from.Format("2006-01")
		toStr := toExcl.AddDate(0, -1, 0).Format("2006-01")
		cacheKey := fromStr + ":" + toStr

		cache.mu.Lock()
		if cache.key == cacheKey && time.Since(cache.at) < aggregatesCacheTTL {
			body := cache.body
			cache.mu.Unlock()
			writeAggregates(w, body)
			return
		}
		cache.mu.Unlock()

		rows, err := postgres.GetHeatmapSourceRows(r.Context(), cfg.DB, from, toExcl)
		if err != nil {
			log.Printf("Error getting heatmap rows: %v", err)
			http.Error(w, "Error getting heatmap aggregates", http.StatusInternalServerError)
			return
		}

		aggRows, total := heatmap.Aggregate(rows)
		body, err := json.Marshal(models.HeatmapAggregatesResponse{
			From: fromStr, To: toStr, TotalIncidents: total, Rows: aggRows,
		})
		if err != nil {
			log.Printf("Error marshaling heatmap response: %v", err)
			http.Error(w, "Error getting heatmap aggregates", http.StatusInternalServerError)
			return
		}

		cache.mu.Lock()
		cache.key, cache.body, cache.at = cacheKey, body, time.Now()
		cache.mu.Unlock()

		writeAggregates(w, body)
	}
}

func writeAggregates(w http.ResponseWriter, body []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=3600")
	w.Write(body)
}
