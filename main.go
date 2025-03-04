package main

import (
        "fmt"
        "net/http"
)

func main() {
        fmt.Println("Starting hello-world server...")
        http.HandleFunc("/health", healthServer)
        http.HandleFunc("/hello", helloServer)
        http.HandleFunc("/bye", byeServer)
        if err := http.ListenAndServe(":8080", nil); err != nil {
                panic(err)
        }
}

func healthServer(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "Healthy")
}

func helloServer(w http.ResponseWriter, r *http.Request) {
        fmt.Fprint(w, "Hello Cindy!")
}

func byeServer(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "Bye Cindy!")
}

