package main

import (
    "net/http"
    "github.com/bmizerany/pat"
)

func (app *application) routes() http.Handler {
    mux := pat.New()
    mux.Get("/", app.session.Enable(noSurf(app.authenticate(http.HandlerFunc(app.home)))))
    mux.Get("/snippet/create", app.session.Enable(noSurf(app.authenticate(app.requireAuthenticatedUser(http.HandlerFunc(app.createSnippetForm))))))
    mux.Post("/snippet/create", app.session.Enable(noSurf(app.authenticate(app.requireAuthenticatedUser(http.HandlerFunc(app.createSnippet))))))
    mux.Get("/snippet/:id", app.session.Enable(noSurf(app.authenticate(http.HandlerFunc(app.showSnippet)))))

    mux.Get("/user/signup", app.session.Enable(noSurf(app.authenticate(http.HandlerFunc(app.signupUserForm)))))
    mux.Post("/user/signup", app.session.Enable(noSurf(app.authenticate(http.HandlerFunc(app.signupUser)))))
    mux.Get("/user/login", app.session.Enable(noSurf(app.authenticate(http.HandlerFunc(app.loginUserForm)))))
    mux.Post("/user/login", app.session.Enable(noSurf(app.authenticate(http.HandlerFunc(app.loginUser)))))
    mux.Post("/user/logout", app.session.Enable(noSurf(app.authenticate(app.requireAuthenticatedUser(http.HandlerFunc(app.logoutUser))))))

    mux.Get("/ping", http.HandlerFunc(ping))
    mux.Get("/about", app.session.Enable(noSurf(app.authenticate(http.HandlerFunc(app.about)))))
    mux.Get("/user/profile", app.session.Enable(noSurf(app.authenticate(app.requireAuthenticatedUser(http.HandlerFunc(app.profile))))))

    mux.Get("/user/change-password", app.session.Enable(noSurf(app.authenticate(app.requireAuthenticatedUser(http.HandlerFunc(app.changePasswordForm))))))
    mux.Post("/user/change-password", app.session.Enable(noSurf(app.authenticate(app.requireAuthenticatedUser(http.HandlerFunc(app.changePassword))))))

    // Serve static files under ./ui/static
    fs := http.FileServer(http.Dir("./ui/static"))
    // Add endpoint for the static files
    mux.Get("/static/", http.StripPrefix("/static", fs))

    return app.recoverPanic(app.logRequest(secureHeaders(mux)))
}
