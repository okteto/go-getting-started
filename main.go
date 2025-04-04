package main

import (
        "fmt"
        "net/http"
        "sync/atomic"
        "time"
)

var (
        // Simple counters for metrics
        healthzHitsCount uint64
        mainHitsCount    uint64
)

func main() {
        fmt.Println("Starting hello-world server...")
        
        // Register handlers
        http.HandleFunc("/", helloServer)
        http.HandleFunc("/healthz", healthzHandler)
        http.HandleFunc("/metrics", metricsHandler)
        
        // Start server
        fmt.Println("Server listening on :8080")
        if err := http.ListenAndServe(":8080", nil); err != nil {
                panic(err)
        }
}

func helloServer(w http.ResponseWriter, r *http.Request) {
        // Increment the main endpoint counter
        atomic.AddUint64(&mainHitsCount, 1)
        
        fmt.Fprint(w, "Welcome to Okteto - Accelerating Cloud Native Development with Seamless Kubernetes Integration!")
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
        // Increment the health endpoint counter
        atomic.AddUint64(&healthzHitsCount, 1)
        
        // Set response headers
        w.Header().Set("Content-Type", "application/json")
        
        // Return a simple health status
        fmt.Fprintf(w, `{"status":"healthy","timestamp":"%s"}`, time.Now().Format(time.RFC3339))
}

func metricsHandler(w http.ResponseWriter, r *http.Request) {
        // Set response headers for Prometheus format
        w.Header().Set("Content-Type", "text/plain")
        
        // Return metrics in Prometheus format
        fmt.Fprintf(w, "# HELP healthz_endpoint_hits_total Total number of hits to the healthz endpoint\n")
        fmt.Fprintf(w, "# TYPE healthz_endpoint_hits_total counter\n")
        fmt.Fprintf(w, "healthz_endpoint_hits_total %d\n", atomic.LoadUint64(&healthzHitsCount))
        
        fmt.Fprintf(w, "# HELP main_endpoint_hits_total Total number of hits to the main endpoint\n")
        fmt.Fprintf(w, "# TYPE main_endpoint_hits_total counter\n")
        fmt.Fprintf(w, "main_endpoint_hits_total %d\n", atomic.LoadUint64(&mainHitsCount))
}
