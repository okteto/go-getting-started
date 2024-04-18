package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer(t *testing.T) {

	rw := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "", nil)

	helloServer(rw, r)

	if rw.Code != http.StatusOK {
		t.Fatal("expected 200 status code")
	}

	if rw.Body.String() != "Hello world!" {
		t.Fatal(`invalid response. Expected: "Hello world!"`)
	}
}
