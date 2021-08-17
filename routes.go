package toushi

import (
	"net/http"
)

// Routes return http.Handler for http server
// should be called after all custum path handler
// be registed
func (r *Router) Routes(middlewares ...func(http.Handler) http.Handler) http.Handler {
	r.router.NotFound = &NotFountResponse
	r.router.MethodNotAllowed = MethodNotAllowResponse
	r.Get("/v1/healthcheck", healthCheckHandler)

	middlewares = append(middlewares, r.buildIns()...)
	h := http.Handler(r.router)
	for _, m := range middlewares {
		h = m(h)
	}
	return h
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
