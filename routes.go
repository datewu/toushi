package toushi

import (
	"net/http"
)

// Routes return http.Handler for http server
// should be called after all custum path handler
// be registed
func (g *RouterGroup) Routes(middlewares ...Middleware) http.Handler {
	g.r.router.NotFound = &NotFountResponse
	g.r.router.MethodNotAllowed = MethodNotAllowResponse
	g.Get("/v1/healthcheck", healthCheckHandler)

	middlewares = append(middlewares, g.r.buildIns()...)
	h := http.Handler(g.r.router)
	mm := h.ServeHTTP
	for _, m := range middlewares {
		mm = m(mm)
	}
	return http.HandlerFunc(mm)
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
func (g *RouterGroup) NewPath(method, path string, handler http.HandlerFunc) {
	for _, v := range g.middlewares {
		handler = v(handler)
	}
	g.r.router.HandlerFunc(method, path, handler)
}

// Get ...
func (g *RouterGroup) Get(path string, handler http.HandlerFunc) {
	g.NewPath(http.MethodGet, path, handler)
}

// Post ...
func (g *RouterGroup) Post(path string, handler http.HandlerFunc) {
	g.NewPath(http.MethodPost, path, handler)
}

// Put ...
func (g *RouterGroup) Put(path string, handler http.HandlerFunc) {
	g.NewPath(http.MethodPut, path, handler)
}

// Patch ...
func (g *RouterGroup) Patch(path string, handler http.HandlerFunc) {
	g.NewPath(http.MethodPatch, path, handler)
}

// Delete ...
func (g *RouterGroup) Delete(path string, handler http.HandlerFunc) {
	g.NewPath(http.MethodDelete, path, handler)
}
