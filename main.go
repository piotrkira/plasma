package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Config struct {
	TargetServerUrl   string
	MaxRequestsPerMin int
}

type Plasma struct {
	proxy       *httputil.ReverseProxy
	rateLimiter *RateLimiter
}

func NewPlasma(config *Config) *Plasma {
	targetServerUrl, err := url.Parse(config.TargetServerUrl)
	if err != nil {
		log.Fatalf("Target Server URL is incorrect %v", err)
	}
	return &Plasma{
		proxy:       httputil.NewSingleHostReverseProxy(targetServerUrl),
		rateLimiter: NewRateLimiter(config.MaxRequestsPerMin),
	}
}

func (p *Plasma) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("Got request: [%s] %s", r.Method, r.URL.Path)

	if !p.rateLimiter.IsAllowed(r.RemoteAddr) {
		http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
		log.Printf("Rate limit exceeded for %s", r.RemoteAddr)
		return
	}

	// Handle FAQ endpoint
	if r.URL.Path == "/faq" {
		HandleFAQ(w, r)
		return
	}

	buffer, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error reading request body: %v", err)
		return
	}
	bodyString := string(buffer)
	threat := GetThreat(r.URL.RawQuery, bodyString)
	if threat != "" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		log.Printf("Blocked request due to detected threat: %s", threat)
		return
	}
	r.Body = io.NopCloser(bytes.NewBuffer(buffer))

	p.proxy.ServeHTTP(w, r)
}

func main() {
	config := &Config{TargetServerUrl: "http://localhost:8000", MaxRequestsPerMin: 10}
	plasma := NewPlasma(config)
	server := http.Server{
		Addr:    ":8080",
		Handler: plasma,
	}

	log.Printf("Starting Plasma server at %s", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
