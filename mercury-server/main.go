package main

import(
    "os"
    "fmt"
    "github.com/fabiocolacio/mercury"
)

const systemConf string = "~/.config/mercury/server.toml"

func main() {
    var server mercury.Server

    if args := os.Args[1:]; len(args) > 0 {
        confPath := args[0]

        var err error
        server, err = mercury.NewServer(confPath)

        if err != nil {
            fmt.Printf("Failed to load configuration file '%s': %s\n", confPath, err)
            return
        }
    } else if _, err := os.Stat(systemConf); err == nil {
        fmt.Println("No configuration file specified.")
        fmt.Printf("Falling back to %s\n", systemConf)
        server, err = mercury.NewServer(systemConf)

        if err != nil {
            fmt.Printf("Failed to load configuration file '%s': %s\n", systemConf, err)
            return
        }
    } else {
        fmt.Println("No configuration file found in ~/.config/mercury, and none specified")
        return
    }

    err := server.ListenAndServe()

    if err != nil {
        fmt.Println(err)
    }

    fmt.Println("Closing Server...")
}

