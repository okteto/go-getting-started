package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

func TestHelloServer(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "GET request to root path",
			method:         "GET",
			path:           "/",
			expectedStatus: http.StatusOK,
			expectedBody:   "🚀 Okteto makes cloud-native development absolutely amazing! 🌊",
		},
		{
			name:           "POST request to root path",
			method:         "POST",
			path:           "/",
			expectedStatus: http.StatusOK,
			expectedBody:   "🚀 Okteto makes cloud-native development absolutely amazing! 🌊",
		},
		{
			name:           "PUT request to root path",
			method:         "PUT",
			path:           "/",
			expectedStatus: http.StatusOK,
			expectedBody:   "🚀 Okteto makes cloud-native development absolutely amazing! 🌊",
		},
		{
			name:           "DELETE request to root path",
			method:         "DELETE",
			path:           "/",
			expectedStatus: http.StatusOK,
			expectedBody:   "🚀 Okteto makes cloud-native development absolutely amazing! 🌊",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, tt.path, nil)
			if err != nil {
				t.Fatalf("Could not create request: %v", err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(helloServer)

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}

			if rr.Body.String() != tt.expectedBody {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rr.Body.String(), tt.expectedBody)
			}
		})
	}
}

func TestHelloServerResponseHeaders(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(helloServer)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
	}

	contentType := rr.Header().Get("Content-Type")
	if contentType != "text/plain; charset=utf-8" {
		t.Errorf("Expected Content-Type 'text/plain; charset=utf-8', got '%s'", contentType)
	}
}

func TestHelloServerWithDifferentPaths(t *testing.T) {
	paths := []string{"/", "/test", "/api"}
	
	for _, path := range paths {
		t.Run("path_"+path, func(t *testing.T) {
			req, err := http.NewRequest("GET", path, nil)
			if err != nil {
				t.Fatalf("Could not create request: %v", err)
			}

			rr := httptest.NewRecorder()
			
			mux := http.NewServeMux()
			mux.HandleFunc("/", helloServer)
			
			mux.ServeHTTP(rr, req)

			if rr.Code != http.StatusOK {
				t.Errorf("Expected status code %d for path %s, got %d", http.StatusOK, path, rr.Code)
			}

			expectedBody := "🚀 Okteto makes cloud-native development absolutely amazing! 🌊"
			if rr.Body.String() != expectedBody {
				t.Errorf("Expected body '%s' for path %s, got '%s'", expectedBody, path, rr.Body.String())
			}
		})
	}
}

func TestHealthHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(healthHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expectedBody := `{"status":"healthy","message":"Service is running"}`
	if rr.Body.String() != expectedBody {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expectedBody)
	}

	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
	}
}

func TestHealthHandlerPrometheusCounter(t *testing.T) {
	// Reset the counter before test
	healthCheckCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "health_endpoint_requests_total_test",
		Help: "Total number of requests to the health endpoint",
	})

	// Make multiple requests to health endpoint
	for i := 0; i < 3; i++ {
		req, err := http.NewRequest("GET", "/health", nil)
		if err != nil {
			t.Fatalf("Could not create request: %v", err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(healthHandler)
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
		}
	}
}

func TestHealthEndpointDifferentMethods(t *testing.T) {
	methods := []string{"GET", "POST", "PUT", "DELETE"}
	
	for _, method := range methods {
		t.Run("method_"+method, func(t *testing.T) {
			req, err := http.NewRequest(method, "/health", nil)
			if err != nil {
				t.Fatalf("Could not create request: %v", err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(healthHandler)
			handler.ServeHTTP(rr, req)

			if rr.Code != http.StatusOK {
				t.Errorf("Expected status code %d for method %s, got %d", http.StatusOK, method, rr.Code)
			}

			expectedBody := `{"status":"healthy","message":"Service is running"}`
			if rr.Body.String() != expectedBody {
				t.Errorf("Expected body '%s' for method %s, got '%s'", expectedBody, method, rr.Body.String())
			}
		})
	}
}

func TestServerRouting(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", helloServer)
	mux.HandleFunc("/health", healthHandler)

	tests := []struct {
		path         string
		expectedBody string
		contentType  string
	}{
		{
			path:         "/",
			expectedBody: "🚀 Okteto makes cloud-native development absolutely amazing! 🌊",
			contentType:  "text/plain; charset=utf-8",
		},
		{
			path:         "/health",
			expectedBody: `{"status":"healthy","message":"Service is running"}`,
			contentType:  "application/json",
		},
	}

	for _, tt := range tests {
		t.Run("route_"+tt.path, func(t *testing.T) {
			req, err := http.NewRequest("GET", tt.path, nil)
			if err != nil {
				t.Fatalf("Could not create request: %v", err)
			}

			rr := httptest.NewRecorder()
			mux.ServeHTTP(rr, req)

			if rr.Code != http.StatusOK {
				t.Errorf("Expected status code %d for path %s, got %d", http.StatusOK, tt.path, rr.Code)
			}

			if rr.Body.String() != tt.expectedBody {
				t.Errorf("Expected body '%s' for path %s, got '%s'", tt.expectedBody, tt.path, rr.Body.String())
			}

			if tt.contentType != "" {
				contentType := rr.Header().Get("Content-Type")
				if contentType != tt.contentType {
					t.Errorf("Expected Content-Type '%s' for path %s, got '%s'", tt.contentType, tt.path, contentType)
				}
			}
		})
	}
}

func BenchmarkHelloServer(b *testing.B) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		b.Fatalf("Could not create request: %v", err)
	}

	handler := http.HandlerFunc(helloServer)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
	}
}

func BenchmarkHealthHandler(b *testing.B) {
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		b.Fatalf("Could not create request: %v", err)
	}

	handler := http.HandlerFunc(healthHandler)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
	}
}