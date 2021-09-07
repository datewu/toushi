package toushi

import (
	"encoding/json"
	"net/http"
)

type Envelope map[string]interface{}

func WriteStr(w http.ResponseWriter, status int, msg string, headers http.Header) {
	for k, v := range headers {
		w.Header()[k] = v
	}
	w.WriteHeader(status)
	w.Write([]byte(msg))
}

func WriteJSON(w http.ResponseWriter, status int, data Envelope, headers http.Header) {
	js, err := json.Marshal(data)
	if err != nil {
		msg := Envelope{"error": err}
		WriteJSON(w, http.StatusInternalServerError, msg, nil)
		return
	}
	for k, v := range headers {
		w.Header()[k] = v
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)
}
