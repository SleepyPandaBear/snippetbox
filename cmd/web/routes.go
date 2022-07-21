package main

import (
    "net/http"
    "github.com/bmizerany/pat"
)

func (app *application) routes() http.Handler {
    mux := pat.New()
    mux.Get("/", http.HandlerFunc(app.home))
    mux.Get("/snippet/create", http.HandlerFunc(app.createSnippetForm))
    mux.Post("/snippet/create", http.HandlerFunc(app.createSnippet))
    mux.Get("/snippet/:id", http.HandlerFunc(app.showSnippet))

    // Serve static files under ./ui/static
    fs := http.FileServer(http.Dir("./ui/static"))
    // Add endpoint for the static files
    mux.Get("/static/", http.StripPrefix("/static", fs))

    return app.recoverPanic(app.logRequest(secureHeaders(mux)))
}
