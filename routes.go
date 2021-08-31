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

// Group add a prefix to all path, for each Gropu call
// prefix will accumulate while middleware don't
func (g *RouterGroup) Group(path string, mds ...Middleware) *RouterGroup {
	gp := &RouterGroup{
		r:      g.r,
		prefix: g.prefix + path,
	}
	if mds != nil {
		gp.middlewares = append(gp.middlewares, mds...)
	}
	return gp
}

// NewHandler handle new http request
func (g *RouterGroup) NewHandler(method, path string, handler http.HandlerFunc) {
	for _, v := range g.middlewares {
		handler = v(handler)
	}
	g.r.router.HandlerFunc(method, path, handler)
}

// Get ...
func (g *RouterGroup) Get(path string, handler http.HandlerFunc) {
	g.NewHandler(http.MethodGet, path, handler)
}

// Post ...
func (g *RouterGroup) Post(path string, handler http.HandlerFunc) {
	g.NewHandler(http.MethodPost, path, handler)
}

// Put ...
func (g *RouterGroup) Put(path string, handler http.HandlerFunc) {
	g.NewHandler(http.MethodPut, path, handler)
}

// Patch ...
func (g *RouterGroup) Patch(path string, handler http.HandlerFunc) {
	g.NewHandler(http.MethodPatch, path, handler)
}

// Delete ...
func (g *RouterGroup) Delete(path string, handler http.HandlerFunc) {
	g.NewHandler(http.MethodDelete, path, handler)
}
