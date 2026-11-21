package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strings"
)

// The message under the logo. Change it and save the file to see hot reload.
const message = "Built, deployed, and tested on Kubernetes."

//go:embed static/index.html
var static embed.FS

// info is what the page shows.
type info struct {
	Message string `json:"message"`
	Runtime string `json:"runtime"`
	Pod     string `json:"pod"`
}

func main() {
	fmt.Println("Starting hello-world server...")
	http.HandleFunc("GET /{$}", page)
	http.HandleFunc("GET /api/info", hello)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

// page serves the page with the Okteto logo.
func page(w http.ResponseWriter, r *http.Request) {
	http.ServeFileFS(w, r, static, "static/index.html")
}

// hello returns the message, the Go version, and the pod that answers.
func hello(w http.ResponseWriter, r *http.Request) {
	pod, _ := os.Hostname()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(info{
		Message: message,
		Runtime: "Go " + strings.TrimPrefix(runtime.Version(), "go"),
		Pod:     pod,
	})
}
