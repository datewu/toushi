package toushi

import (
	"errors"

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

// NewRouterGroup return a new routergroup
func NewRouterGroup(conf *Config) (*RouterGroup, error) {
	if conf == nil {
		return nil, errors.New("no router config provided")
	}
	r := router{
		router: httprouter.New(),
		config: conf,
	}
	return &RouterGroup{r: &r}, nil
}

// DefaultRouterGroup return a new routergroup with default config
func DefaultRouterGroup() *RouterGroup {
	r := router{
		router: httprouter.New(),
		config: DefaultConf(),
	}
	return &RouterGroup{r: &r}
}
