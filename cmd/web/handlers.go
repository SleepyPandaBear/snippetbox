package main

import (
    "net/http"
    "fmt"
    "strconv"
    "html/template"
    "spbear/snippetbox/pkg/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
    // If we are not on the index page, show not found page.
    if r.URL.Path != "/" {
        app.notFound(w)
        return
    }

    s, err := app.snippets.Latest()
    if err != nil {
        app.serverError(w, err)
        return
    }
    for _, snippet := range s {
        fmt.Fprintf(w, "%v\n", snippet)
    }

    /*
    files := []string{
        "./ui/html/home.page.tmpl",
        "./ui/html/base.layout.tmpl",
        "./ui/html/footer.partial.tmpl",
    }

    // Try to parse our html template file
    ts, err := template.ParseFiles(files...)
    if err != nil {
        app.serverError(w, err)
        return
    }

    // Try to execute the template. `nil` argument is for the data we pass to
    // the template, here we don't pass anyting so its nil.
    err = ts.Execute(w, nil)
    if err != nil {
        app.serverError(w, err)
    }
    */
}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
    // r.URL.Quury().Get("id") returns value of id from url - /snippet?id=123
    // returns 123
    // Validate that id is int greater than 0
    id, err := strconv.Atoi(r.URL.Query().Get("id"))
    if err != nil || id < 1 {
        app.notFound(w)
        return
    }

    s, err := app.snippets.Get(id)
    if err == models.ErrNoRecord {
        app.notFound(w)
        return
    } else if err != nil {
        app.serverError(w, err)
        return
    }

    data := &templateData{Snippet: s}

    files := []string{
        "./ui/html/show.page.tmpl",
        "./ui/html/base.layout.tmpl",
        "./ui/html/footer.partial.tmpl",
    }

    ts, err := template.ParseFiles(files...)
    if err != nil {
        app.serverError(w, err)
        return
    }

    err = ts.Execute(w, data)
    if err != nil {
        app.serverError(w, err)
        return
    }
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
    // Check if we are using POST method on this endpoint.
    if r.Method != "POST" {
        w.Header().Set("Allow", "POST")

        // This code is replaced with http.Error(...)
        //w.WriteHeader(405)
        //w.Write([]byte("Method not allowed"))
        app.clientError(w, http.StatusMethodNotAllowed)
        return
    }

    title := "O snail"
    content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi"
    expires := "7"

    id, err := app.snippets.Insert(title, content, expires)
    if err != nil {
        app.serverError(w, err)
        return
    }

    http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)
}

