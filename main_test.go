package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
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
			name:           "GET request returns hello world",
			method:         "GET",
			path:           "/",
			expectedStatus: http.StatusOK,
			expectedBody:   "Hello world!",
		},
		{
			name:           "POST request returns hello world",
			method:         "POST",
			path:           "/",
			expectedStatus: http.StatusOK,
			expectedBody:   "Hello world!",
		},
		{
			name:           "PUT request returns hello world",
			method:         "PUT",
			path:           "/",
			expectedStatus: http.StatusOK,
			expectedBody:   "Hello world!",
		},
		{
			name:           "DELETE request returns hello world",
			method:         "DELETE",
			path:           "/",
			expectedStatus: http.StatusOK,
			expectedBody:   "Hello world!",
		},
		{
			name:           "HEAD request returns hello world",
			method:         "HEAD",
			path:           "/",
			expectedStatus: http.StatusOK,
			expectedBody:   "Hello world!",
		},
		{
			name:           "Request with query parameters",
			method:         "GET",
			path:           "/?name=test&value=123",
			expectedStatus: http.StatusOK,
			expectedBody:   "Hello world!",
		},
		{
			name:           "Request with different path",
			method:         "GET",
			path:           "/test",
			expectedStatus: http.StatusOK,
			expectedBody:   "Hello world!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a request
			req, err := http.NewRequest(tt.method, tt.path, nil)
			if err != nil {
				t.Fatal(err)
			}

			// Create a ResponseRecorder to record the response
			rr := httptest.NewRecorder()

			// Call the handler function
			handler := http.HandlerFunc(helloServer)
			handler.ServeHTTP(rr, req)

			// Check the status code
			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}

			// Check the response body (HEAD requests don't return body)
			if tt.method != "HEAD" && rr.Body.String() != tt.expectedBody {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rr.Body.String(), tt.expectedBody)
			}
		})
	}
}

func TestHelloServerWithHeaders(t *testing.T) {
	// Test with custom headers
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	
	// Add custom headers
	req.Header.Set("User-Agent", "test-agent")
	req.Header.Set("Accept", "text/plain")
	
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(helloServer)
	handler.ServeHTTP(rr, req)
	
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	
	if rr.Body.String() != "Hello world!" {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), "Hello world!")
	}
}

func TestHelloServerWithBody(t *testing.T) {
	// Test with request body
	body := strings.NewReader("test request body")
	req, err := http.NewRequest("POST", "/", body)
	if err != nil {
		t.Fatal(err)
	}
	
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(helloServer)
	handler.ServeHTTP(rr, req)
	
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	
	if rr.Body.String() != "Hello world!" {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), "Hello world!")
	}
}

func TestMainFunction(t *testing.T) {
	// Since the main function starts an HTTP server that blocks,
	// we can't test it directly. Instead, we'll test that the
	// server can be started on a test port.
	
	// This is a simple smoke test to ensure the handler is registered
	mux := http.NewServeMux()
	mux.HandleFunc("/", helloServer)
	
	// Create a test server
	ts := httptest.NewServer(mux)
	defer ts.Close()
	
	// Make a request to the test server
	resp, err := http.Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", resp.Status)
	}
	
	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	
	if string(body) != "Hello world!" {
		t.Errorf("expected body 'Hello world!'; got %v", string(body))
	}
}

func BenchmarkHelloServer(b *testing.B) {
	// Benchmark the handler performance
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		b.Fatal(err)
	}
	
	handler := http.HandlerFunc(helloServer)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
	}
}