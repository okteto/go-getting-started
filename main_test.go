package main

import (
        "encoding/json"
        "net/http"
        "net/http/httptest"
        "testing"
)

func TestHelloServer(t *testing.T) {
        // Create a request to pass to our handler
        req, err := http.NewRequest("GET", "/", nil)
        if err != nil {
                t.Fatal(err)
        }

        // Create a ResponseRecorder to record the response
        rr := httptest.NewRecorder()
        handler := http.HandlerFunc(helloServer)

        // Serve the HTTP request to our handler
        handler.ServeHTTP(rr, req)

        // Check the status code
        if status := rr.Code; status != http.StatusOK {
                t.Errorf("handler returned wrong status code: got %v want %v",
                        status, http.StatusOK)
        }

        // Check the response body
        expected := "Hello from Okteto! Accelerating cloud-native development with real-time Kubernetes deployments!"
        if rr.Body.String() != expected {
                t.Errorf("handler returned unexpected body: got %v want %v",
                        rr.Body.String(), expected)
        }
}

func TestHelloServerWithDifferentMethod(t *testing.T) {
        // Test with different HTTP methods
        methods := []string{"GET", "POST", "PUT", "DELETE"}
        
        for _, method := range methods {
                // Create a request with the current method
                req, err := http.NewRequest(method, "/", nil)
                if err != nil {
                        t.Fatal(err)
                }

                // Create a ResponseRecorder to record the response
                rr := httptest.NewRecorder()
                handler := http.HandlerFunc(helloServer)

                // Serve the HTTP request to our handler
                handler.ServeHTTP(rr, req)

                // Check the status code is always OK regardless of method
                if status := rr.Code; status != http.StatusOK {
                        t.Errorf("%s: handler returned wrong status code: got %v want %v",
                                method, status, http.StatusOK)
                }

                // Check the response body is consistent across methods
                expected := "Hello from Okteto! Accelerating cloud-native development with real-time Kubernetes deployments!"
                if rr.Body.String() != expected {
                        t.Errorf("%s: handler returned unexpected body: got %v want %v",
                                method, rr.Body.String(), expected)
                }
        }
}

func TestHelloServerWithDifferentPaths(t *testing.T) {
        // Test with different paths
        paths := []string{"/", "/hello", "/okteto", "/test/path"}
        
        for _, path := range paths {
                // Create a request with the current path
                req, err := http.NewRequest("GET", path, nil)
                if err != nil {
                        t.Fatal(err)
                }

                // Create a ResponseRecorder to record the response
                rr := httptest.NewRecorder()
                handler := http.HandlerFunc(helloServer)

                // Serve the HTTP request to our handler
                handler.ServeHTTP(rr, req)

                // Check the status code is always OK regardless of path
                if status := rr.Code; status != http.StatusOK {
                        t.Errorf("Path %s: handler returned wrong status code: got %v want %v",
                                path, status, http.StatusOK)
                }

                // Check the response body is consistent across paths
                expected := "Hello from Okteto! Accelerating cloud-native development with real-time Kubernetes deployments!"
                if rr.Body.String() != expected {
                        t.Errorf("Path %s: handler returned unexpected body: got %v want %v",
                                path, rr.Body.String(), expected)
                }
        }
}

func TestHealthCheckHandler(t *testing.T) {
        // Create a request to pass to our handler
        req, err := http.NewRequest("GET", "/health", nil)
        if err != nil {
                t.Fatal(err)
        }

        // Create a ResponseRecorder to record the response
        rr := httptest.NewRecorder()
        handler := http.HandlerFunc(healthCheckHandler)

        // Serve the HTTP request to our handler
        handler.ServeHTTP(rr, req)

        // Check the status code
        if status := rr.Code; status != http.StatusOK {
                t.Errorf("handler returned wrong status code: got %v want %v",
                        status, http.StatusOK)
        }

        // Check the content type
        contentType := rr.Header().Get("Content-Type")
        if contentType != "application/json" {
                t.Errorf("handler returned wrong content type: got %v want %v",
                        contentType, "application/json")
        }

        // Parse the JSON response
        var response map[string]string
        if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
                t.Errorf("Failed to parse JSON response: %v", err)
        }

        // Check the status field
        if status, exists := response["status"]; !exists || status != "healthy" {
                t.Errorf("handler returned unexpected status: got %v want %v",
                        status, "healthy")
        }
}