package main

import (
	"fmt"
	"net/http"

	"github.com/common-nighthawk/go-figure"
)

func main() {
	fmt.Println("Starting hello-world server...")
	http.HandleFunc("/", helloServer)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

func helloServer(w http.ResponseWriter, r *http.Request) {
	ascii := figure.NewFigure("Welcome to Civo Navigate SF", "", true)
	fmt.Fprint(w, ascii.String())
}
