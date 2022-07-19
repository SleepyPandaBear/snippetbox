package main

import (
    "net/http"
    "log"
    "flag"
)

func main() {
    // Parse value from command line (value must be dereferenced when used)
    ip := flag.String("ip", ":8080", "Ip address and port of the application")
    flag.Parse()

    mux := http.NewServeMux()
    mux.HandleFunc("/", home)
    mux.HandleFunc("/snippet", showSnippet)
    mux.HandleFunc("/snippet/create", createSnippet)

    // Serve static files under ./ui/static
    fs := http.FileServer(http.Dir("./ui/static"))
    // Add endpoint for the static files
    mux.Handle("/static/", http.StripPrefix("/static", fs))

    // run server from console as: w:\snippetbox>go run .\cmd\web\
    log.Printf("Starting server on port %s \n", *ip)
    err := http.ListenAndServe(*ip, mux)
    log.Fatal(err)
}
