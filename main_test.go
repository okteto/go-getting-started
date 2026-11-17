package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHelloServer(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	helloServer(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	if got, want := rec.Body.String(), "Hello world!"; got != want {
		t.Fatalf("expected body %q, got %q", want, got)
	}
}
