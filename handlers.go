package toushi

import (
	"net/http"
)

// HealthCheckHandler is a handler for health check
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	data := Envelope{
		"status": "available",
	}
	WriteJSON(w, http.StatusOK, data, nil)
}
