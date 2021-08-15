package toushi

import (
	"expvar"
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

func (ro *Router) enabledCORS(next http.Handler) http.Handler {
	if ro.config.CORS.TrustedOrigins == nil {
		return next
	}
	middle := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Origin")

		// Add the "Vary: Access-Control-Request-Method" header.
		w.Header().Add("Vary", "Access-Control-Request-Method")

		origin := r.Header.Get("Origin")

		if origin != "" {
			for i := range ro.config.CORS.TrustedOrigins {
				if origin == ro.config.CORS.TrustedOrigins[i] {
					w.Header().Set("Access-Control-Allow-Origin", origin)

					// Check if the request has the HTTP method OPTIONS and contains the
					// "Access-Control-Request-Method" header. If it does, then we treat
					// it as a preflight request.
					if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
						// Set the necessary preflight response headers, as discussed
						// previously.
						w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, PUT, PATCH, DELETE")
						w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")

						// Write the headers along with a 200 OK status and return from
						// the middleware with no further action.
						w.WriteHeader(http.StatusOK)
						return
					}
				}
			}
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(middle)
}
func (ro *Router) rateLimit(next http.Handler) http.Handler {
	if !ro.config.Limiter.Enabled {
		return next
	}
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}
	var (
		clients = make(map[string]*client)
		mu      sync.Mutex
	)
	delOld := func(interval time.Duration) {
		for {
			time.Sleep(interval)
			mu.Lock()
			for k, v := range clients {
				if time.Since(v.lastSeen) > 3*time.Minute {
					delete(clients, k)
				}
			}
			mu.Unlock()
		}
	}
	go delOld(time.Minute)

	middle := func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ServerErrResponse(err).ServeHTTP(w, r)
			return
		}
		mu.Lock()
		if _, existed := clients[ip]; !existed {
			clients[ip] = &client{
				limiter: rate.NewLimiter(rate.Limit(ro.config.Limiter.Rps),
					ro.config.Limiter.Burst),
			}
		}
		clients[ip].lastSeen = time.Now()
		if !clients[ip].limiter.Allow() {
			mu.Unlock()
			RateLimitExceededResponse.ServeHTTP(w, r)
			return
		}
		mu.Unlock()
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(middle)
}

func (ro *Router) metrics(next http.Handler) http.Handler {
	if !ro.config.Metrics {
		return next
	}
	totalRequestReceived := expvar.NewInt("total_requests_received")
	totalResponsesSend := expvar.NewInt("total_responses_send")
	totalProcessingTimeMicroseconds := expvar.NewInt("total_processing_time_us")
	middle := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		totalRequestReceived.Add(1)
		next.ServeHTTP(w, r)
		totalResponsesSend.Add(1)
		duration := time.Since(start).Microseconds()
		totalProcessingTimeMicroseconds.Add(duration)
	}
	return http.HandlerFunc(middle)
}

func recoverPanic(next http.Handler) http.Handler {
	middle := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				ServerErrResponse(err).ServeHTTP(w, r)
				return
			}
		}()
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(middle)
}
