package main

import (
    "net/http"
    "fmt"
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
