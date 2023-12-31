package main

import (
    "net/http"
    "fmt"
    "github.com/justinas/nosurf"
    "context"
    "spbear/snippetbox/pkg/models"
)

// Classic middleware pattern, we do our logic and then return next handler.
func secureHeaders(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("X-XSS-Protection", "1; mode=block")
        w.Header().Set("X-Frame-Options", "deny")
        next.ServeHTTP(w, r)
    })
}

func (app *application) logRequest(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL)
        next.ServeHTTP(w, r)
    })
}

// This only recover from panics in this go routine. If we spawn another go
// routines we need to check if they panic in different way (those go routines
// failiures can bring down the server). This different way is shown bellow:
// func myHandler(w http.ResponseWriter, r *http.Request) {
//     ...
//     // Spin up a new goroutine to do some background processing.
//     go func() {
//         defer func() {
//             if err := recover(); err != nil {
//             log.Println(fmt.Errorf("%s\n%s", err, debug.Stack()))
//         }
//         }()
//         doSomeBackgroundProcessing()
//     }()
//     w.Write([]byte("OK"))
// }
func (app *application) recoverPanic(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                w.Header().Set("Connection", "close")
                app.serverError(w, fmt.Errorf("%s", err))
            }
        }()

        next.ServeHTTP(w, r)
    })
}

// Midlleware that checks if the user is authorized
func (app *application) requireAuthenticatedUser(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if app.authenticatedUser(r) == nil {
            app.session.Put(r, "wantedURL", r.URL.Path)
            http.Redirect(w, r, "/user/login", 302)
            return
        }

        next.ServeHTTP(w, r)
    })
}

func noSurf(next http.Handler) http.Handler {
    csrfHandler := nosurf.New(next)
    csrfHandler.SetBaseCookie(http.Cookie{
        HttpOnly: true,
        Path: "/",
        Secure: true,
    })
    return csrfHandler
}

func (app *application) authenticate(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        exists := app.session.Exists(r, "userID")
        if !exists {
            next.ServeHTTP(w, r)
            return
        }

        // Check if the user exists - we could delete the user for various reasons.
        user, err := app.users.Get(app.session.GetInt(r, "userID"))
        if err == models.ErrNoRecord {
            app.session.Remove(r, "userID")
            next.ServeHTTP(w, r)
            return
        } else if err != nil {
            app.serverError(w, err)
            return
        }

        ctx := context.WithValue(r.Context(), contextKeyUser, user)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
