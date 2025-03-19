package main

import (
        "fmt"
        "net/http"

        "github.com/prometheus/client_golang/prometheus"
        "github.com/prometheus/client_golang/prometheus/promauto"
        "github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
        // Define a counter for health endpoint requests
        healthEndpointHits = promauto.NewCounter(prometheus.CounterOpts{
                Name: "health_endpoint_total_requests",
                Help: "The total number of requests to the /health endpoint",
        })
)

func main() {
        fmt.Println("Starting hello-world server...")
        
        // Register handlers
        http.HandleFunc("/", helloServer)
        http.HandleFunc("/health", healthHandler)
        
        // Expose Prometheus metrics on /metrics endpoint
        http.Handle("/metrics", promhttp.Handler())
        
        fmt.Println("Server is ready to handle requests at :8080")
        if err := http.ListenAndServe(":8080", nil); err != nil {
                panic(err)
        }
}

func helloServer(w http.ResponseWriter, r *http.Request) {
        // Only respond to exact path "/"
        if r.URL.Path != "/" {
                http.NotFound(w, r)
                return
        }
        fmt.Fprint(w, "Hello from Okteto! Supercharging Kubernetes development with lightning-fast deployments and seamless cloud-native workflows.")
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
        // Increment the counter for health endpoint requests
        healthEndpointHits.Inc()
        
        // Return a simple health status
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        fmt.Fprint(w, `{"status":"healthy"}`)
}
