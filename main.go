package main

import (
    "net/http"
    "log"
    "fmt"
    "strconv"
)

func home(w http.ResponseWriter, r *http.Request) {
    // If we are not on the index page, show not found page.
    if r.URL.Path != "/" {
        http.NotFound(w, r)
        return
    }

    w.Write([]byte("Hello world"))
}

func showSnippet(w http.ResponseWriter, r *http.Request) {

    // r.URL.Quury().Get("id") returns value of id from url - /snippet?id=123
    // returns 123
    // Validate that id is int greater than 0
    id, err := strconv.Atoi(r.URL.Query().Get("id"))
    if err != nil || id < 1 {
        http.NotFound(w, r)
        return
    }

    fmt.Fprintf(w, "Display snippet with id: %d", id)
}

func createSnippet(w http.ResponseWriter, r *http.Request) {
    // Check if we are using POST method on this endpoint.
    if r.Method != "POST" {
        w.Header().Set("Allow", "POST")

        // This code is replaced with http.Error(...)
        //w.WriteHeader(405)
        //w.Write([]byte("Method not allowed"))
        http.Error(w, "Method not allowed", 405)
        return
    }

    w.Write([]byte("Create snippet"))
}

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/", home)
    mux.HandleFunc("/snippet", showSnippet)
    mux.HandleFunc("/snippet/create", createSnippet)

    log.Println("Starting server on port 8080")

    err := http.ListenAndServe(":8080", mux)
    log.Fatal(err)

}
