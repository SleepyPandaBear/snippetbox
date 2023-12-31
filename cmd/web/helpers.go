package main

import (
    "fmt"
    "bytes"
    "net/http"
    "runtime/debug"
    "time"

    "github.com/justinas/nosurf"
    "spbear/snippetbox/pkg/models"
)

func (app *application) serverError(w http.ResponseWriter, err error) {
    trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
    app.errorLog.Println(2, trace)

    if app.debugMode {
        http.Error(w, trace, http.StatusInternalServerError)
    } else {
        http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    }
}

func (app *application) clientError(w http.ResponseWriter, status int) {
    http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
    app.clientError(w, http.StatusNotFound)
}

// We get template from our cache. Rendering is seperated into two steps, first
// we try to execute template into a buffer, and if this succedes we can render
// it to the user. If this step fails we have an error and we can show some
// error message to the user.
func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
    ts, ok := app.templateCache[name]

    if !ok {
        app.serverError(w, fmt.Errorf("The template %s does not exist", name))
        return
    }

    buf := &bytes.Buffer{}

    // Execute template into the buffer first
    err := ts.Execute(buf, app.addDefaultData(td, r))
    if err != nil {
        app.serverError(w, err)
        return
    }

    buf.WriteTo(w)
}

func (app *application) addDefaultData(td *templateData, r *http.Request) *templateData {
    if td == nil {
        td = &templateData{}
    }

    td.CurrentYear = time.Now().Year()
    td.AuthenticatedUser = app.authenticatedUser(r)
    td.CSRFToken = nosurf.Token(r)

    // Add flash message to the template data if one exists.
    td.Flash = app.session.PopString(r, "flash")

    return td
}

func (app *application) authenticatedUser(r *http.Request) *models.User {
    user, ok := r.Context().Value(contextKeyUser).(*models.User)
    if !ok {
        return nil
    }

    return user
}
