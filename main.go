package main

import (
        "fmt"
        "net/http"

        "github.com/prometheus/client_golang/prometheus"
        "github.com/prometheus/client_golang/prometheus/promauto"
        "github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
        // Create a counter for the /health endpoint
        healthChecks = promauto.NewCounter(prometheus.CounterOpts{
                Name: "hello_world_health_checks_total",
                Help: "The total number of health check requests",
        })
)

func main() {
        fmt.Println("Starting hello-world server...")
        
        // Register handlers
        http.HandleFunc("/", helloServer)
        http.HandleFunc("/health", healthCheckHandler)
        
        // Expose Prometheus metrics endpoint
        http.Handle("/metrics", promhttp.Handler())
        
        if err := http.ListenAndServe(":8080", nil); err != nil {
                panic(err)
        }
}

func helloServer(w http.ResponseWriter, r *http.Request) {
        fmt.Fprint(w, "Hello from Okteto! Accelerating cloud-native development with real-time Kubernetes deployments!")
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
        // Increment the counter each time this endpoint is called
        healthChecks.Inc()
        
        // Return a simple health status
        w.Header().Set("Content-Type", "application/json")
        fmt.Fprint(w, `{"status":"healthy"}`)
}
