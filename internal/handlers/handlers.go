package handlers

import (
	"cm_open_api/internal/config"
	"cm_open_api/internal/postgres"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func SetupRouter(cfg *config.Config) http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/v1/api/outages", getOutages(cfg)).Methods("GET")
	r.HandleFunc("/v1/api/outages/{message_id}/source", getSource(cfg)).Methods("GET")

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET"},
	})
	return c.Handler(r)
}

func getOutages(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		outages, err := postgres.GetOutages(r.Context(), cfg.DB)
		if err != nil {
			log.Printf("Error getting outages: %v", err)
			http.Error(w, "Error getting outages", http.StatusInternalServerError)
			return
		}

		// Конвертируем EventStart и EventStop в строку
		for i := range outages {
			outages[i].EventStartStr = outages[i].EventStart.Format(time.RFC3339)
			if outages[i].EventStop != nil {
				eventStopStr := outages[i].EventStop.Format(time.RFC3339)
				outages[i].EventStopStr = &eventStopStr
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(outages)
	}
}

func getSource(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		messageID := vars["message_id"]
		parts := strings.Split(messageID, ":")
		if len(parts) < 3 {
			http.Error(w, "Invalid message_id format", http.StatusBadRequest)
			return
		}

		source, err := postgres.GetSource(r.Context(), cfg.DB, parts[1], parts[0])
		if err != nil {
			log.Printf("Error getting source: %v", err)
			http.Error(w, "Error getting source", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(source)
	}
}
