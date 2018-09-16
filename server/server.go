package main

import(
    "os"
    "fmt"
    "net/http"
    "github.com/BurntSushi/toml"
)

type serverConf struct {
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
    fmt.Printf("HTTPS Address: %s\n", conf.HttpsAddr)
    fmt.Printf("Cert File: %s\n", conf.CertFile)
    fmt.Printf("Key File: %s\n", conf.KeyFile)

    http.HandleFunc("/", httpHandler)

    err := http.ListenAndServeTLS(conf.HttpsAddr, conf.CertFile, conf.KeyFile, nil)

    if err != nil {
        fmt.Println(err)
    }

    fmt.Println("Closing Server...")
}

func readConf(path string) (conf serverConf,err error) {
    _, err = toml.DecodeFile(path, &conf)
    return conf, err
}

func httpHandler(res http.ResponseWriter, req *http.Request) {
    res.Write([]byte("<h1>Hi There!</h1>"))
}

