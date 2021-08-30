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
	r.router.HandlerFunc(http.MethodGet, "/v1/healthcheck", healthCheckHandler)

	middlewares = append(middlewares, r.buildIns()...)
	h := http.Handler(r.router)
	for _, m := range middlewares {
		h = m(h)
	}
	return h
}

// Group add a prefix to all path
func (g *RouterGroup) Group(path string) *RouterGroup {
	gp := &RouterGroup{
		r:      g.r,
		prefix: g.prefix + path,
	}
	return gp
}

// NewPath handle new http request
func (r *Router) NewPath(method, path string, handler http.HandlerFunc) {
	r.router.HandlerFunc(method, path, handler)
}

// Get ...
func (g *RouterGroup) Get(path string, handler http.HandlerFunc) {
	g.r.NewPath(http.MethodGet, path, handler)
}

// Post ...
func (g *RouterGroup) Post(path string, handler http.HandlerFunc) {
	g.r.NewPath(http.MethodPost, path, handler)
}

// Put ...
func (g *RouterGroup) Put(path string, handler http.HandlerFunc) {
	g.r.NewPath(http.MethodPut, path, handler)
}

// Patch ...
func (g *RouterGroup) Patch(path string, handler http.HandlerFunc) {
	g.r.NewPath(http.MethodPatch, path, handler)
}

// Delete ...
func (g *RouterGroup) Delete(path string, handler http.HandlerFunc) {
	g.r.NewPath(http.MethodDelete, path, handler)
}
