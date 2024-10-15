package util

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"net/http"
	"net/url"
	"sort"
	"strings"
)

// PayLib struct to handle payment data
type PayLib struct {
	RequestData  map[string]string
	ResponseData map[string]string
}

// NewPayLib creates a new PayLib instance
func NewPayLib() *PayLib {
	return &PayLib{
		RequestData:  make(map[string]string),
		ResponseData: make(map[string]string),
	}
}

// AddRequestData adds key-value pair to request data map
func (p *PayLib) AddRequestData(key, value string) {
	if value != "" {
		p.RequestData[key] = value
	}
}

// AddResponseData adds key-value pair to response data map
func (p *PayLib) AddResponseData(key, value string) {
	if value != "" {
		p.ResponseData[key] = value
	}
}

// GetResponseData retrieves value from response data
func (p *PayLib) GetResponseData(key string) string {
	if value, exists := p.ResponseData[key]; exists {
		return value
	}
	return ""
}

// CreateRequestUrl creates the payment request URL with a secure hash
func (p *PayLib) CreateRequestUrl(baseUrl, hashSecret string) string {
	// Sort request data by key
	keys := make([]string, 0, len(p.RequestData))
	for k := range p.RequestData {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Build the query string
	var queryString strings.Builder
	for _, k := range keys {
		v := url.QueryEscape(p.RequestData[k])
		queryString.WriteString(url.QueryEscape(k) + "=" + v + "&")
	}

	// Remove the last '&'
	qs := queryString.String()
	if len(qs) > 0 {
		qs = qs[:len(qs)-1]
	}

	// Generate the secure hash
	vnpSecureHash := HmacSHA512(hashSecret, qs)

	// Append secure hash to query string
	fullUrl := baseUrl + "?" + qs + "&vnp_SecureHash=" + url.QueryEscape(vnpSecureHash)
	return fullUrl
}

// ValidateSignature validates the response signature
func (p *PayLib) ValidateSignature(inputHash, secretKey string) bool {
	rawData := p.getRawResponseData()
	calculatedHash := HmacSHA512(secretKey, rawData)
	return strings.EqualFold(calculatedHash, inputHash)
}

// getRawResponseData creates a query string from response data (excluding vnp_SecureHash and vnp_SecureHashType)
func (p *PayLib) getRawResponseData() string {
	var data strings.Builder

	// Sort response data by key
	keys := make([]string, 0, len(p.ResponseData))
	for k := range p.ResponseData {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Build the raw data string (exclude "vnp_SecureHash" and "vnp_SecureHashType")
	for _, k := range keys {
		if k == "vnp_SecureHash" || k == "vnp_SecureHashType" {
			continue
		}
		v := url.QueryEscape(p.ResponseData[k])
		data.WriteString(url.QueryEscape(k) + "=" + v + "&")
	}

	// Remove the last '&'
	raw := data.String()
	if len(raw) > 0 {
		raw = raw[:len(raw)-1]
	}
	return raw
}

func HmacSHA512(key, data string) string {
	h := hmac.New(sha512.New, []byte(key))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

// GetIpAddress retrieves the client's IP address from the request
func GetIpAddress(r *http.Request) string {
	// First check if X-Forwarded-For is set
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" || strings.ToLower(ip) == "unknown" {
		// Otherwise, use the RemoteAddr from the request
		ip = r.RemoteAddr
	}

	return ip
}
