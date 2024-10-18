package logging

import (
	"log"
	"net/http"
	"time"
)

func LogRequests(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a response recorder to capture the status code
		recorder := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}

		// Forward the request
		next.ServeHTTP(recorder, r)

		// Log the request details
		log.Printf("Method: %s, Path: %s, Status: %d, Duration: %v",
			r.Method, r.URL.Path, recorder.statusCode, time.Since(start))
	}
}

// responseRecorder is used to capture the status code
type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rec *responseRecorder) WriteHeader(code int) {
	rec.statusCode = code
	rec.ResponseWriter.WriteHeader(code)
}
