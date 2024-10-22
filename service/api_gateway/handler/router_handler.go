package handler

import (
	"io"
	"net/http"
	"net/url"
	"strings"
)

// Microservice base URLs (replace with the actual URLs for your services)
const (
	oauthServiceBaseURL           = "http://localhost:8080/auth"
	productServiceBaseURL         = "http://localhost:8081/api/products"
	userServiceBaseURL            = "http://localhost:8082/api/users"
	newsServiceBaseURL            = "http://localhost:8083/api/news"
	cartItemServiceBaseURL        = "http://localhost:8084/api/cartItems"
	cartServiceBaseURL            = "http://localhost:8085/api/carts"
	categoryServiceBaseURL        = "http://localhost:8086/api/categories"
	courierServiceBaseURL         = "http://localhost:8087/api/couriers"
	discountServiceBaseURL        = "http://localhost:8088/api/discounts"
	freightRateServiceBaseURL     = "http://localhost:8089/api/freightRates"
	orderServiceBaseURL           = "http://localhost:8090/api/orders"
	orderDetailsServiceBaseURL    = "http://localhost:8091/api/orderDetails"
	productDiscountServiceBaseURL = "http://localhost:8092/api/productDiscounts"
	reviewServiceBaseURL          = "http://localhost:8093/api/reviews"
	paymentServiceBaseURL         = "http://localhost:8094/api/payments"
	voucherServiceBaseURL         = "http://localhost:8095/api/vouchers"
	mailServiceBaseURL            = "http://localhost:8096/api/mail"
	momoServiceBaseURL            = "http://localhost:8097/api/momo"
	vnpayServiceBaseURL           = "http://localhost:8098/api/vnpay"
	authServiceBaseURL            = "http://localhost:8099/api/authentication"
)

var productServiceURLs = []string{
	"http://localhost:8081/api/products",
	"http://localhost:9081/api/products", // Another instance of product service
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

	case strings.HasPrefix(path, "/api/news"):
		targetURL := newsServiceBaseURL + strings.TrimPrefix(path, "/api/news")
		ForwardRequest(w, r, targetURL)

	case strings.HasPrefix(path, "/api/cartItems"):
		targetURL := cartItemServiceBaseURL + strings.TrimPrefix(path, "/api/cartItems")
		ForwardRequest(w, r, targetURL)

	case strings.HasPrefix(path, "/api/carts"):
		targetURL := cartServiceBaseURL + strings.TrimPrefix(path, "/api/carts")
		ForwardRequest(w, r, targetURL)

	case strings.HasPrefix(path, "/api/categories"):
		targetURL := categoryServiceBaseURL + strings.TrimPrefix(path, "/api/categories")
		ForwardRequest(w, r, targetURL)

	case strings.HasPrefix(path, "/api/couriers"):
		targetURL := courierServiceBaseURL + strings.TrimPrefix(path, "/api/couriers")
		ForwardRequest(w, r, targetURL)

	case strings.HasPrefix(path, "/api/discounts"):
		targetURL := discountServiceBaseURL + strings.TrimPrefix(path, "/api/discounts")
		ForwardRequest(w, r, targetURL)

	case strings.HasPrefix(path, "/api/freightRates"):
		targetURL := freightRateServiceBaseURL + strings.TrimPrefix(path, "/api/freightRates")
		ForwardRequest(w, r, targetURL)

	case strings.HasPrefix(path, "/api/orders"):
		targetURL := orderServiceBaseURL + strings.TrimPrefix(path, "/api/orders")
		ForwardRequest(w, r, targetURL)

	case strings.HasPrefix(path, "/api/orderDetails"):
		targetURL := orderDetailsServiceBaseURL + strings.TrimPrefix(path, "/api/orderDetails")
		ForwardRequest(w, r, targetURL)

	case strings.HasPrefix(path, "/api/productDiscounts"):
		targetURL := productDiscountServiceBaseURL + strings.TrimPrefix(path, "/api/productDiscounts")
		ForwardRequest(w, r, targetURL)

	case strings.HasPrefix(path, "/api/reviews"):
		targetURL := reviewServiceBaseURL + strings.TrimPrefix(path, "/api/reviews")
		ForwardRequest(w, r, targetURL)

	case strings.HasPrefix(path, "/api/payments"):
		targetURL := paymentServiceBaseURL + strings.TrimPrefix(path, "/api/payments")
		ForwardRequest(w, r, targetURL)

	case strings.HasPrefix(path, "/api/vouchers"):
		targetURL := voucherServiceBaseURL + strings.TrimPrefix(path, "/api/vouchers")
		ForwardRequest(w, r, targetURL)

	case strings.HasPrefix(path, "/api/mail"):
		targetURL := mailServiceBaseURL + strings.TrimPrefix(path, "/api/mail")
		ForwardRequest(w, r, targetURL)

	case strings.HasPrefix(path, "/api/momo"):
		targetURL := momoServiceBaseURL + strings.TrimPrefix(path, "/api/momo")
		ForwardRequest(w, r, targetURL)

	case strings.HasPrefix(path, "/api/vnpay"):
		targetURL := vnpayServiceBaseURL + strings.TrimPrefix(path, "/api/vnpay")
		ForwardRequest(w, r, targetURL)

	case strings.HasPrefix(path, "/api/authentication"):
		targetURL := authServiceBaseURL + strings.TrimPrefix(path, "/api/authentication")
		ForwardRequest(w, r, targetURL)

	default:
		http.Error(w, "Not Found", http.StatusNotFound)
	}
}
