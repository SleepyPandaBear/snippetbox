package main

import (
    "net/http"
    "log"
)

func home(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Hello world"))
}

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/", home)

    log.Println("Starting server on port 8080")

    err := http.ListenAndServe(":8080", mux)
    log.Fatal(err)

}
