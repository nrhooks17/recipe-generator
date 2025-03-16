// Package handlers provides HTTP request handlers for the recipe generator API.
// It contains handlers for various endpoints including health checks and recipe operations.
package handler

import (
	"encoding/json"
	"net/http"
)

// HealthHandler provides HTTP handlers for health-related endpoints.
// It is used to check the status and availability of the service.
type HealthHandler struct {
	// HealthHandler currently has no fields but can be extended
	// to include dependencies like database connections or loggers.
}

// HealthCheck returns an HTTP handler function that performs a health check.
// The handler responds with a JSON object containing a status field set to "ok"
// to indicate that the service is running properly.
//
// Returns:
//   - http.HandlerFunc: A handler function that writes a JSON response with status information.
func (h *HealthHandler) HealthCheck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set the content type header to application/json
		w.Header().Set("Content-Type", "application/json")
		
		// Encode and send a simple JSON response with status "ok"
		err := json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
		if err != nil {
			// If encoding fails, log the error (could be extended to use a logger)
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}
