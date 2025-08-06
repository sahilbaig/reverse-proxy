package reverseproxy

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metrics
var (
	requestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "reverse_proxy_requests_total",
			Help: "Total number of requests proxied",
		},
		[]string{"method", "upstream", "status"},
	)

	requestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "reverse_proxy_request_duration_seconds",
			Help:    "Time taken to proxy requests",
			Buckets: []float64{0.01, 0.05, 0.1, 0.5, 1, 2.5}, // Buckets in seconds
		},
		[]string{"method", "upstream", "status"},
	)
)

// responseRecorder captures HTTP status codes
type responseRecorder struct {
	http.ResponseWriter
	status int
}

func (r *responseRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func Reverseproxy(target string) http.HandlerFunc {
	parsedURL, err := url.Parse(target)
	if err != nil {
		log.Fatalf("Invalid proxy target URL: %v", err)
	}

	transport := &http.Transport{
		MaxIdleConns:        1000,
		MaxIdleConnsPerHost: 1000,
		IdleConnTimeout:     90 * time.Second,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2: true,
	}

	proxy := httputil.NewSingleHostReverseProxy(parsedURL)
	proxy.Transport = transport

	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		recorder := &responseRecorder{w, http.StatusOK}

		defer func() {
			duration := time.Since(start).Seconds()
			status := fmt.Sprintf("%d", recorder.status)

			// Record metrics
			requestsTotal.WithLabelValues(r.Method, target, status).Inc()
			requestDuration.WithLabelValues(r.Method, target, status).Observe(duration)
		}()

		log.Printf("→ Proxying %s %s to %s", r.Method, r.URL.Path, target)
		proxy.ServeHTTP(recorder, r)
	}
}