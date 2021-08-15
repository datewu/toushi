package toushi

import "github.com/julienschmidt/httprouter"

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

// Router holds all paths relative funcs
type Router struct {
	router *httprouter.Router
	config *Config
}

// New return a new router...
func New(cnf *Config) *Router {
	r := Router{
		router: httprouter.New(),
	}
	if cnf == nil {
		cnf = DefaultConf()
	}
	r.config = cnf
	return &r
}
