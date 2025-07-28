package main

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	healthCheckCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "health_endpoint_requests_total",
		Help: "Total number of requests to the health endpoint",
	})
)

func init() {
	prometheus.MustRegister(healthCheckCounter)
}

func main() {
	fmt.Println("Starting hello-world server...")
	http.HandleFunc("/", helloServer)
	http.HandleFunc("/health", healthHandler)
	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

func helloServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "🚀 Okteto makes cloud-native development absolutely amazing! 🌊")
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	healthCheckCounter.Inc()
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, `{"status":"healthy","message":"Service is running"}`)
}
