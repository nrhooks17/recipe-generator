package handlers

import (
	"encoding/json"
	"net/http"
)

// HealthHandler struct for health check handlers
type HealthHandler struct {
}

// HealthCheck returns an ok string for health checks.
func (h *HealthHandler) HealthCheck() http.HandlerFunc {
	// return an anonymous function that returns a json with the status of my app

	return func(w http.ResponseWriter, r *http.Request) {
		err := json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
		if err != nil {
			return
		}
	}
}
