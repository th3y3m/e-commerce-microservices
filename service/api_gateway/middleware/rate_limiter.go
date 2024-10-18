package middleware

import (
	"net/http"
	"time"
)

var rateLimiter = make(map[string]time.Time)

// Basic rate limiting logic
func RateLimit(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientIP := r.RemoteAddr
		now := time.Now()

		// Check if the client has made a request recently
		if lastRequestTime, exists := rateLimiter[clientIP]; exists {
			// Allow one request per second
			if now.Sub(lastRequestTime) < time.Second {
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}
		}

		// Record current request time
		rateLimiter[clientIP] = now

		// Forward to the next handler if within the rate limit
		next.ServeHTTP(w, r)
	}
}
