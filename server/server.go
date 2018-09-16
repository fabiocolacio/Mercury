package main

import(
    "fmt"
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
    confPath := "server.toml"
    conf, err := readConf(confPath)

    if err != nil {
        fmt.Printf("Failed to load configuration file '%s': %s", confPath, err)
    }

    fmt.Println("Starting server with the following configuration:")
    fmt.Printf("HTTP Address: %s\n", conf.HttpAddr)
    fmt.Printf("HTTPS Address: %s\n", conf.HttpsAddr)
    fmt.Printf("Cert File: %s\n", conf.CertFile)
    fmt.Printf("Key File: %s\n", conf.KeyFile)

    http.HandleFunc("/", httpHandler)

    ch := make(chan error)

    go func() {
        ch <- http.ListenAndServe(conf.HttpAddr, nil)
    }()

    go func() {
        ch <- http.ListenAndServeTLS(conf.HttpsAddr, conf.CertFile, conf.KeyFile, nil)
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
    if req.TLS == nil {
        // TODO: Redirect to a secure connection
    }

    res.Write([]byte("<h1>Hi There!</h1>"))
}

