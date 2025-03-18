package main

import (
        "net/http"
        "net/http/httptest"
        "strings"
        "testing"
        "time"
        
        "github.com/prometheus/client_golang/prometheus/promhttp"
)

func TestHelloServer(t *testing.T) {
        // Create a request to pass to our handler
        req, err := http.NewRequest("GET", "/", nil)
        if err != nil {
                t.Fatal(err)
        }

        // Create a ResponseRecorder to record the response
        rr := httptest.NewRecorder()
        
        // Use the instrumented handler for testing
        handler := instrumentHandler(helloServer, homepageCalls)

        // Call the handler with our request and response recorder
        handler.ServeHTTP(rr, req)

        // Check the status code
        if status := rr.Code; status != http.StatusOK {
                t.Errorf("handler returned wrong status code: got %v want %v",
                        status, http.StatusOK)
        }

        // Check the response body contains expected content
        if !strings.Contains(rr.Body.String(), "Welcome to Okteto!") {
                t.Errorf("handler returned unexpected body: got %v, expected to contain 'Welcome to Okteto!'",
                        rr.Body.String())
        }

        // Check that the response contains HTML content
        if !strings.Contains(rr.Body.String(), "<html>") {
                t.Errorf("handler returned non-HTML content")
        }

        // Check that the response mentions Cloud Native development
        if !strings.Contains(rr.Body.String(), "Cloud Native") {
                t.Errorf("handler returned body without 'Cloud Native' mention")
        }
}

func TestHealthzHandler(t *testing.T) {
        // Test cases for the healthz endpoint
        tests := []struct {
                name           string
                timeSinceStart time.Duration
                wantStatus     int
                wantBodyContains string
        }{
                {
                        name:           "Service starting up",
                        timeSinceStart: 5 * time.Second,
                        wantStatus:     http.StatusServiceUnavailable,
                        wantBodyContains: "Service is starting up",
                },
                {
                        name:           "Service ready",
                        timeSinceStart: 15 * time.Second,
                        wantStatus:     http.StatusOK,
                        wantBodyContains: "OK - Pod:",
                },
        }

        // Save the original startTime
        originalStartTime := startTime
        
        for _, tt := range tests {
                t.Run(tt.name, func(t *testing.T) {
                        // Set the startTime for this test case
                        startTime = time.Now().Add(-tt.timeSinceStart)
                        
                        // Create a request to pass to our handler
                        req, err := http.NewRequest("GET", "/healthz", nil)
                        if err != nil {
                                t.Fatal(err)
                        }

                        // Create a ResponseRecorder to record the response
                        rr := httptest.NewRecorder()
                        
                        // Use the instrumented handler for testing
                        handler := instrumentHandler(healthzHandler, healthzCalls)

                        // Call the handler with our request and response recorder
                        handler.ServeHTTP(rr, req)

                        // Check the status code
                        if status := rr.Code; status != tt.wantStatus {
                                t.Errorf("handler returned wrong status code: got %v want %v",
                                        status, tt.wantStatus)
                        }

                        // Check the response body contains expected content
                        if !strings.Contains(rr.Body.String(), tt.wantBodyContains) {
                                t.Errorf("handler returned unexpected body: got %v, expected to contain '%s'",
                                        rr.Body.String(), tt.wantBodyContains)
                        }
                })
        }
        
        // Restore the original startTime
        startTime = originalStartTime
}

func TestMetricsEndpoint(t *testing.T) {
        // Increment the counters directly
        healthzCalls.Inc()
        homepageCalls.Inc()
        responseTime.Observe(0.1)
        
        // Create a request to the metrics endpoint
        req, err := http.NewRequest("GET", "/metrics", nil)
        if err != nil {
                t.Fatal(err)
        }
        
        // Create a ResponseRecorder to record the response
        rr := httptest.NewRecorder()
        
        // Serve the metrics
        promhttp.Handler().ServeHTTP(rr, req)
        
        // Check the status code
        if status := rr.Code; status != http.StatusOK {
                t.Errorf("metrics handler returned wrong status code: got %v want %v",
                        status, http.StatusOK)
        }
        
        // Get the response body
        bodyStr := rr.Body.String()
        
        // Verify metrics are present
        t.Logf("Metrics response: %s", bodyStr)
        
        // Check for healthz calls counter
        if !strings.Contains(bodyStr, "hello_world_healthz_calls_total") {
                t.Errorf("Metrics endpoint missing healthz calls counter")
        }
        
        // Check for homepage calls counter
        if !strings.Contains(bodyStr, "hello_world_homepage_calls_total") {
                t.Errorf("Metrics endpoint missing homepage calls counter")
        }
        
        // Check for response time histogram
        if !strings.Contains(bodyStr, "hello_world_response_time_seconds") {
                t.Errorf("Metrics endpoint missing response time histogram")
        }
}