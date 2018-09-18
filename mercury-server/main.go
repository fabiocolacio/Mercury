package main

import(
    "os"
    "fmt"
    "log"
    "github.com/fabiocolacio/mercury"
)

func main() {
    var server mercury.Server
    var err error

    // Load the configuration file.
    // If the user provides one as a cli argument, this is the one that is used.
    // If no path was provided, mercury-server looks for the file ~/.config/mercury/server.toml
    // If no configuration file could be successfully loaded, mercury-server exits
    if args := os.Args[1:]; len(args) > 0 {
        confPath := args[0]
        server, err = mercury.NewServer(confPath)
        mercury.Assertf(err != nil, "Failed to load configuration file '%s': %s", confPath, err)
    } else if _, err = os.Stat(systemConf()); err == nil {
        server, err = mercury.NewServer(systemConf())
        mercury.Assertf(err != nil, "Failed to load configuration file '%s': %s", systemConf(), err)
    } else {
        log.Fatalf("No config was specified, and file '%s' was not found.", systemConf())
    }

    // Start the server, logging errors to stdout
    err = server.ListenAndServe()
    mercury.Assert(err != nil, err)
}

// The default location of mercury-server's configuration file.
// This file will be the fallback if no file was specified.
func systemConf() string {
    return fmt.Sprintf("%s/.config/mercury/server.toml", mercury.HomeDir())
}

