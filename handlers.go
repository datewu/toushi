package toushi

import (
	"net/http"
)

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := Envelope{
		"status": "available",
	}
	WriteJSON(w, http.StatusOK, data, nil)
}
