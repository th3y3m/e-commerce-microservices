package middleware

import (
	"net/http"
	"strings"
)

// Global cache
var cache = make(map[string][]byte)

// Basic caching logic
func CacheMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check if the response is cached
		if cachedResp, found := cache[r.URL.Path]; found {
			w.Header().Set("Content-Type", "application/json")
			w.Write(cachedResp)
			return
		}

		// Capture the response using a buffer
		recorder := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK, body: &strings.Builder{}}
		next.ServeHTTP(recorder, r)

		// Cache the response if it's a success
		if recorder.statusCode == http.StatusOK {
			cache[r.URL.Path] = []byte(recorder.body.String())
		}

		// Write the response
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(recorder.body.String()))
	}
}

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
	body       *strings.Builder
}

func (rec *responseRecorder) WriteHeader(code int) {
	rec.statusCode = code
	rec.ResponseWriter.WriteHeader(code)
}

func (rec *responseRecorder) Write(b []byte) (int, error) {
	rec.body.Write(b)
	return rec.ResponseWriter.Write(b)
}
