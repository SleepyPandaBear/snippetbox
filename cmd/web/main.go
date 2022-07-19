package main

import (
    "net/http"
    "log"
)

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/", home)
    mux.HandleFunc("/snippet", showSnippet)
    mux.HandleFunc("/snippet/create", createSnippet)

    // Serve static files under ./ui/static
    fs := http.FileServer(http.Dir("./ui/static"))
    // Add endpoint for the static files
    mux.Handle("/static/", http.StripPrefix("/static", fs))

    // run server from console as: w:\snippetbox>go run .\cmd\web\
    log.Println("Starting server on port 8080")
    err := http.ListenAndServe(":8080", mux)
    log.Fatal(err)
}
