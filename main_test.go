package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNowShowing(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/movie", nil)
	rec := httptest.NewRecorder()

	nowShowing(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	var got movie
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("invalid response: %v", err)
	}
	if got.Studio != studio {
		t.Errorf("studio = %q, want %q", got.Studio, studio)
	}
	if got.Title != title {
		t.Errorf("title = %q, want %q", got.Title, title)
	}
	if !strings.HasPrefix(got.Starring, "Go ") {
		t.Errorf("starring = %q, want the Go version", got.Starring)
	}
}

func TestPoster(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	poster(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if !strings.Contains(rec.Body.String(), `id="title"`) {
		t.Error("the page doesn't contain the poster")
	}
}
