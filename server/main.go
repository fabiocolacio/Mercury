package main

import(
    "os"
    "fmt"
    "strings"
    "net/http"
    "github.com/BurntSushi/toml"
)

type serverConf struct {
    HttpAddr  string
    HttpsAddr string
    CertFile  string
    KeyFile   string
}

func main() {
    var conf serverConf

    if args := os.Args[1:]; len(args) > 0 {
        confPath := args[0]

        var err error
        conf, err = readConf(confPath)

        if err != nil {
            fmt.Printf("Failed to load configuration file '%s': %s\n", confPath, err)
            return
        }
    } else {
        fmt.Println("No configuration file specified.")
        return
    }

    fmt.Println("Starting server with the following configuration:")
    fmt.Printf("HTTP Address: %s\n", conf.HttpAddr)
    fmt.Printf("HTTPS Address: %s\n", conf.HttpsAddr)
    fmt.Printf("Cert File: %s\n", conf.CertFile)
    fmt.Printf("Key File: %s\n", conf.KeyFile)

    ch := make(chan error)

    go func () {
        ch <- http.ListenAndServe(conf.HttpAddr, http.HandlerFunc(httpHandler))
    }()

    go func () {
        ch <- http.ListenAndServeTLS(conf.HttpsAddr, conf.CertFile, conf.KeyFile, http.HandlerFunc(httpsHandler))
    }()

    if err := <-ch; err != nil {
        fmt.Println(err)
    }

    fmt.Println("Closing Server...")
}

func readConf(path string) (conf serverConf,err error) {
    _, err = toml.DecodeFile(path, &conf)
    return conf, err
}

func httpHandler(res http.ResponseWriter, req *http.Request) {
    // TODO: Use the port from the configuration rather than assuming 443

    port := 443
    host := strings.Split(req.Host, ":")[0]
    path := req.URL.Path

    dest := fmt.Sprintf("https://%s:%d%s", host, port, path)

    fmt.Println(dest)

    http.Redirect(res, req, dest, http.StatusTemporaryRedirect)
}

func httpsHandler(res http.ResponseWriter, req *http.Request) {
    res.Write([]byte("<h1>Hi There!</h1>"))
}

