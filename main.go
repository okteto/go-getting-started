package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	fmt.Println("Starting hello-world server...")
	http.HandleFunc("/", helloServer)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

func helloServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world from %s!", os.Getenv("OKTETO_NAMESPACE"))
}
