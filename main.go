package main

import(
    "fmt"
    "flag"
    "os"
    "os/signal"
    "syscall"
    "github.com/fabiocolacio/mercury/server"
)

var(
    confPath string
    flagInit bool
)

func init() {
    dconf := fmt.Sprintf("%s/.config/mercury/config.toml", os.Getenv("HOME"))
    flag.StringVar(&confPath, "config", dconf, "The configuration file to load.")
    flag.BoolVar(&flagInit, "init", false, "Creates necessary database tables if they do not exist")
    flag.Parse()
}

func main() {
    // Creates a new server with the details from the configuration file.
    // If there was an error loading the file, the program quits.
    serv, err := server.NewServer(confPath)

    // Free resources allocated by the Server after exiting main
    defer serv.Close()

    // Exit if there was an error creating the server
    server.Assertf(err == nil, "Failed to load configuration file '%s': %s", confPath, err)


    if flagInit {
        err = serv.ResetDB()
        server.Assertf(err == nil, "Failed to initialize database: %s", err)
    }

    // Start handling connections
    go serv.ListenAndServe()

    // Handle these signals so that the server can cleanly exit before closing the program
    sig := make(chan os.Signal, 1)
    signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
    <-sig
}

