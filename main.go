package main

import (
	"fmt"
	"net/http"

	"github.com/okteto/go-getting-started/pkg/render"
)

func main() {
	fmt.Println("Starting hello-world server...")
	r := render.New()
	http.HandleFunc("/", r.SortByName)
	http.HandleFunc("/age", r.SortByAge)
	http.HandleFunc("/restart", r.SortByRestart)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
