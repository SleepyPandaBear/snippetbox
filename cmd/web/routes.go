package main

import (
    "net/http"
)

func (app *application) routes() http.Handler {
    mux := http.NewServeMux()
    mux.HandleFunc("/", app.home)
    mux.HandleFunc("/snippet", app.showSnippet)
    mux.HandleFunc("/snippet/create", app.createSnippet)

    // Serve static files under ./ui/static
    fs := http.FileServer(http.Dir("./ui/static"))
    // Add endpoint for the static files
    mux.Handle("/static/", http.StripPrefix("/static", fs))

    return app.logRequest(secureHeaders(mux))
}
