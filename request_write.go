package toushi

import (
	"context"
	"encoding/json"
	"net/http"
)

// Envelope is a JSON envelope for better client response
type Envelope map[string]interface{}

// SetValue sets a value on the request context.
func SetValue(r *http.Request, key, value interface{}) {
	ctx := context.WithValue(r.Context(), key, value)
	_ = r.WithContext(ctx)
}

// GetValue gets a value from the request context.
func GetValue(r *http.Request, key interface{}) interface{} {
	return r.Context().Value(key)
}

// WriteStr writes a string to the response.
func WriteStr(w http.ResponseWriter, status int, msg string, headers http.Header) {
	for k, v := range headers {
		w.Header()[k] = v
	}
	w.WriteHeader(status)
	w.Write([]byte(msg))
}

// writeJSON writes a JSON object to the response.
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
