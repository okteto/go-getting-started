package main

import (
        "fmt"
        "net/http"
        "os"
        "time"

        "github.com/prometheus/client_golang/prometheus"
        "github.com/prometheus/client_golang/prometheus/promauto"
        "github.com/prometheus/client_golang/prometheus/promhttp"
)

// For simulating startup time
var startTime = time.Now()

// Define Prometheus metrics
var (
        // Counter for healthz endpoint calls
        healthzCalls = promauto.NewCounter(prometheus.CounterOpts{
                Name: "hello_world_healthz_calls_total",
                Help: "The total number of calls to the healthz endpoint",
        })

        // Counter for homepage calls
        homepageCalls = promauto.NewCounter(prometheus.CounterOpts{
                Name: "hello_world_homepage_calls_total",
                Help: "The total number of calls to the homepage",
        })

        // Histogram for response time
        responseTime = promauto.NewHistogram(prometheus.HistogramOpts{
                Name:    "hello_world_response_time_seconds",
                Help:    "Response time of requests in seconds",
                Buckets: prometheus.DefBuckets,
        })
)

func main() {
        fmt.Println("Starting hello-world server with Prometheus metrics...")
        
        // Register handlers
        http.HandleFunc("/", instrumentHandler(helloServer, homepageCalls))
        http.HandleFunc("/healthz", instrumentHandler(healthzHandler, healthzCalls))
        
        // Register Prometheus metrics endpoint
        http.Handle("/metrics", promhttp.Handler())
        
        // Start server
        if err := http.ListenAndServe(":8080", nil); err != nil {
                panic(err)
        }
}

// instrumentHandler wraps an HTTP handler with metrics instrumentation
func instrumentHandler(next http.HandlerFunc, counter prometheus.Counter) http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
                // Start timer
                start := time.Now()
                
                // Increment counter
                counter.Inc()
                
                // Call the original handler
                next(w, r)
                
                // Record response time
                duration := time.Since(start).Seconds()
                responseTime.Observe(duration)
        }
}

// healthzHandler responds to health check requests
func healthzHandler(w http.ResponseWriter, r *http.Request) {
        // Simple health check - you could add more sophisticated checks here
        // For example, checking database connectivity, external services, etc.
        
        // Simulate startup delay (10 seconds) for readiness
        if time.Since(startTime) < 10*time.Second {
                w.WriteHeader(http.StatusServiceUnavailable)
                fmt.Fprintf(w, "Service is starting up...")
                return
        }
        
        // Return hostname for debugging rolling updates
        hostname, err := os.Hostname()
        if err != nil {
                hostname = "unknown"
        }
        
        w.WriteHeader(http.StatusOK)
        fmt.Fprintf(w, "OK - Pod: %s", hostname)
}

func helloServer(w http.ResponseWriter, r *http.Request) {
        message := `
        <html>
                <head>
                        <title>Okteto Cloud Native App</title>
                        <style>
                                body {
                                        font-family: Arial, sans-serif;
                                        max-width: 800px;
                                        margin: 0 auto;
                                        padding: 20px;
                                        text-align: center;
                                        background-color: #f5f7fa;
                                }
                                h1 {
                                        color: #223c7a;
                                }
                                .container {
                                        background-color: white;
                                        border-radius: 8px;
                                        padding: 20px;
                                        box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
                                        margin-top: 20px;
                                }
                                .highlight {
                                        color: #223c7a;
                                        font-weight: bold;
                                }
                        </style>
                </head>
                <body>
                        <h1>Welcome to Okteto!</h1>
                        <div class="container">
                                <p>Okteto is the <span class="highlight">ultimate platform</span> for Cloud Native development on Kubernetes.</p>
                                <p>With Okteto, you can:</p>
                                <ul style="text-align: left;">
                                        <li>Develop directly in your Kubernetes cluster</li>
                                        <li>Eliminate environment inconsistencies</li>
                                        <li>Accelerate your development feedback loop</li>
                                        <li>Collaborate seamlessly with your team</li>
                                </ul>
                                <p>Experience the future of Cloud Native development today!</p>
                        </div>
                </body>
        </html>
        `
        fmt.Fprint(w, message)
}
