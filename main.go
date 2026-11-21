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

// The line at the top of the poster. Change it and save the file to see hot reload.
const studio = "A Kubernetes Production"

// The title of the movie.
const title = "Captain Okteto"

//go:embed static/index.html
var static embed.FS

// movie is what the poster shows.
type movie struct {
	Studio   string `json:"studio"`
	Title    string `json:"title"`
	Starring string `json:"starring"`
	Pod      string `json:"pod"`
}

func main() {
	fmt.Println("Starting hello-world server...")
	http.HandleFunc("GET /{$}", poster)
	http.HandleFunc("GET /api/movie", nowShowing)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

// poster serves the page that draws the poster.
func poster(w http.ResponseWriter, r *http.Request) {
	http.ServeFileFS(w, r, static, "static/index.html")
}

// nowShowing returns the movie: its studio, its title, the Go version, and the pod that answers.
func nowShowing(w http.ResponseWriter, r *http.Request) {
	pod, _ := os.Hostname()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movie{
		Studio:   studio,
		Title:    title,
		Starring: "Go " + strings.TrimPrefix(runtime.Version(), "go"),
		Pod:      pod,
	})
}
