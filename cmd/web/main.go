package main

import (
    "net/http"
    "log"
    "flag"
    "os"
)

type application struct {
    infoLog *log.Logger
    errorLog *log.Logger
}

func main() {
    // Parse value from command line (value must be dereferenced when used)
    ip := flag.String("ip", ":8080", "Ip address and port of the application")
    flag.Parse()

    // Our custom loggers, one for infos and one for errors.
    // We can also define here to which file we should output the logs (eg.
    // instead of os.Stdout).
    infoLog := log.New(os.Stdout, "INFO:\t", log.Ldate | log.Ltime)
    errorLog := log.New(os.Stderr, "ERROR:\t", log.Ldate | log.Ltime | log.Lshortfile)

    // Our application struct so we can use our custom loggers from handlers.go
    // file. If our handlers are scattered across multiple files, we can use
    // function closures (first function accepts app and it returns another
    // function with ResponseWriter and Request) and pass application struct as
    // argument.
    app := &application {
        infoLog: infoLog,
        errorLog: errorLog,
    }

    // We need to create our own `Server` struct, so we can use our own error
    // logger.
    srv := &http.Server {
        Addr: *ip,
        ErrorLog: errorLog,
        Handler: app.routes(), 
    }

    // run server from console as: w:\snippetbox>go run .\cmd\web\
    infoLog.Printf("Starting server on port %s \n", *ip)
    err := srv.ListenAndServe()
    errorLog.Fatal(err)
}
