package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHello(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/info", nil)
	rec := httptest.NewRecorder()

	hello(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	var got info
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("invalid response: %v", err)
	}
	if got.Message != message {
		t.Errorf("message = %q, want %q", got.Message, message)
	}
	if !strings.HasPrefix(got.Runtime, "Go ") {
		t.Errorf("runtime = %q, want the Go version", got.Runtime)
	}
}

func TestPage(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	page(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if !strings.Contains(rec.Body.String(), `id="message"`) {
		t.Error("the page doesn't show the message")
	}
}
