package main

import (
    "net/http"
    "log"
    "flag"
    "os"
)

func main() {
    // Parse value from command line (value must be dereferenced when used)
    ip := flag.String("ip", ":8080", "Ip address and port of the application")
    flag.Parse()

    // Our custom loggers, one for infos and one for errors.
    // We can also define here to which file we should output the logs (eg.
    // instead of os.Stdout).
    infoLog := log.New(os.Stdout, "INFO:\t", log.Ldate | log.Ltime)
    errorLog := log.New(os.Stderr, "ERROR:\t", log.Ldate | log.Ltime | log.Lshortfile)

    mux := http.NewServeMux()
    mux.HandleFunc("/", home)
    mux.HandleFunc("/snippet", showSnippet)
    mux.HandleFunc("/snippet/create", createSnippet)

    // Serve static files under ./ui/static
    fs := http.FileServer(http.Dir("./ui/static"))
    // Add endpoint for the static files
    mux.Handle("/static/", http.StripPrefix("/static", fs))

    // We need to create our own `Server` struct, so we can use our own error
    // logger.
    srv := &http.Server{
        Addr: *ip,
        ErrorLog: errorLog,
        Handler: mux,
    }

    // run server from console as: w:\snippetbox>go run .\cmd\web\
    infoLog.Printf("Starting server on port %s \n", *ip)
    err := srv.ListenAndServe(*ip, mux)
    errorLog.Fatal(err)
}
