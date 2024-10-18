package handler

import (
	"io"
	"net/http"
	"net/url"
	"strings"
)

// Microservice base URLs (replace with the actual URLs for your services)
const (
	productServiceBaseURL = "http://localhost:8080/api/products"
	userServiceBaseURL    = "http://localhost:8081/api/users"
)

var productServiceURLs = []string{
	"http://localhost:8080/api/products",
	"http://localhost:8082/api/products", // Another instance of product service
}

var currentProductIndex = 0

// ForwardRequest forwards the incoming request to the appropriate microservice
func ForwardRequest(w http.ResponseWriter, r *http.Request, targetURL string) {
	// Create a new request based on the incoming request
	proxyURL, err := url.Parse(targetURL)
	if err != nil {
		http.Error(w, "Invalid service URL", http.StatusInternalServerError)
		return
	}

	proxyReq, err := http.NewRequest(r.Method, proxyURL.String(), r.Body)
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	// Copy original headers
	for name, values := range r.Header {
		for _, value := range values {
			proxyReq.Header.Add(name, value)
		}
	}

	// Perform the request
	client := &http.Client{}
	resp, err := client.Do(proxyReq)
	if err != nil {
		http.Error(w, "Error forwarding the request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Copy response headers
	for name, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(name, value)
		}
	}

	// Set the status code and write the response body
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func getNextProductServiceURL() string {
	// Round-robin load balancing
	url := productServiceURLs[currentProductIndex]
	currentProductIndex = (currentProductIndex + 1) % len(productServiceURLs)
	return url
}

func RouteHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	switch {
	case strings.HasPrefix(path, "/api/products"):
		// Load balance between product service instances
		targetURL := getNextProductServiceURL() + strings.TrimPrefix(path, "/api/products")
		ForwardRequest(w, r, targetURL)

	case strings.HasPrefix(path, "/api/users"):
		targetURL := userServiceBaseURL + strings.TrimPrefix(path, "/api/users")
		ForwardRequest(w, r, targetURL)

	default:
		http.Error(w, "Not Found", http.StatusNotFound)
	}
}
