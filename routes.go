package toushi

import (
	"expvar"
	"net/http"
)

// Routes return http.Handler for http server
// should be called after all custum path handler
// be registed
func (r *Router) Routes() http.Handler {
	r.router.NotFound = &NotFountResponse
	notAllow := func(w http.ResponseWriter, r *http.Request) {
		MethodNotAllowResponse(r.Method).ServeHTTP(w, r)
	}
	r.router.MethodNotAllowed = http.HandlerFunc(notAllow)

	r.router.HandlerFunc(
		http.MethodGet,
		"/v1/healthcheck",
		healthCheckHandler)

	if r.config.Metrics {
		r.router.Handler(
			http.MethodGet,
			"/debug/vars",
			expvar.Handler())
	}
	rlMiddle := r.rateLimit(r.router)
	corsMiddle := r.enabledCORS(rlMiddle)
	recoverMiddle := recoverPanic(corsMiddle)
	return r.metrics(recoverMiddle)
}

// NewPath handle new http request
func (r *Router) NewPath(method, path string, handler http.HandlerFunc) {
	r.router.HandlerFunc(method, path, handler)
}

// Get ...
func (r *Router) Get(path string, handler http.HandlerFunc) {
	r.NewPath(http.MethodGet, path, handler)
}

// Post ...
func (r *Router) Post(path string, handler http.HandlerFunc) {
	r.NewPath(http.MethodPost, path, handler)
}

// Put ...
func (r *Router) Put(path string, handler http.HandlerFunc) {
	r.NewPath(http.MethodPut, path, handler)
}

// Patch ...
func (r *Router) Patch(path string, handler http.HandlerFunc) {
	r.NewPath(http.MethodPatch, path, handler)
}

// Delete ...
func (r *Router) Delete(path string, handler http.HandlerFunc) {
	r.NewPath(http.MethodDelete, path, handler)
}