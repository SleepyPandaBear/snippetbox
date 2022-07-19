package main

import (
    "net/http"
    "fmt"
    "strconv"
    "html/template"
    "log"
)

func home(w http.ResponseWriter, r *http.Request) {
    // If we are not on the index page, show not found page.
    if r.URL.Path != "/" {
        http.NotFound(w, r)
        return
    }

    files := []string{
        "./ui/html/home.page.tmpl",
        "./ui/html/base.layout.tmpl",
        "./ui/html/footer.partial.tmpl",
    }

    // Try to parse our html template file
    ts, err := template.ParseFiles(files...)
    if err != nil {
        log.Println(err.Error())
        http.Error(w, "Internal server error", 500)
        return
    }

    // Try to execute the template. `nil` argument is for the data we pass to
    // the template, here we don't pass anyting so its nil.
    err = ts.Execute(w, nil)
    if err != nil {
        log.Println(err.Error())
        http.Error(w, "Internal server error", 500)
    }
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

