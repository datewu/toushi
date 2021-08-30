package toushi

import (
	"github.com/julienschmidt/httprouter"
)

type Config struct {
	Limiter struct {
		Rps     float64
		Burst   int
		Enabled bool
	}
	CORS struct {
		TrustedOrigins []string
	}
	Metrics bool
}

func DefaultConf() *Config {
	cnf := &Config{Metrics: true}
	cnf.Limiter.Enabled = true
	cnf.Limiter.Rps = 200
	cnf.Limiter.Burst = 10
	cnf.CORS.TrustedOrigins = nil
	return cnf
}

// router holds all paths relative funcs
type router struct {
	router *httprouter.Router
	config *Config
}

// RourerGroup is a group of routes
type RouterGroup struct {
	r           *router
	prefix      string
	middlewares []Middleware
}

// NewGroup return a new routergroup
func NewGroup(cnf *Config) *RouterGroup {
	r := router{
		router: httprouter.New(),
	}
	if cnf == nil {
		cnf = DefaultConf()
	}
	r.config = cnf
	return &RouterGroup{
		r: &r,
	}
}
