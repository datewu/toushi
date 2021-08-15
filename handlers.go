package toushi

import (
	"net/http"
)

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := envelope{
		"status": "available",
	}
	WriteJSON(w, http.StatusOK, data, nil)
}
