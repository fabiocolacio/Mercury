package main

import(
    "flag"
    "os"
    "os/signal"
    "syscall"
    "github.com/fabiocolacio/mercury"
)

var confPath string

func init() {
    flag.StringVar(&confPath, "c",
        "/usr/local/share/com.github.fabiocolacio.mercury-server/config.toml",
        "The configuration file to load.")
    flag.Parse()
}

func main() {
    // Creates a new server with the details from the configuration file.
    // If there was an error loading the file, the program quits.
    server, err := mercury.NewServer(confPath)

    // Free resources allocated by the Server after exiting main
    defer server.Close()

    // Exit if there was an error creating the server
    mercury.Assertf(err == nil, "Failed to load configuration file '%s': %s", confPath, err)

    // Start handling connections
    go server.ListenAndServe()

    // Handle these signals so that the server can cleanly exit before closing the program
    sig := make(chan os.Signal, 1)
    signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
    <-sig
}

