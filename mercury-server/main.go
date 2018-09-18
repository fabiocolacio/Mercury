package main

import(
    "os"
    "fmt"
    "github.com/fabiocolacio/mercury"
)

func main() {
    var server mercury.Server

    // Load the configuration file.
    // If the user provides one as a cli argument, this is the one that is used.
    // If no path was provided, mercury-server looks for the file ~/.config/mercury/server.toml
    // If no configuration file could be successfully loaded, mercury-server exits
    if args := os.Args[1:]; len(args) > 0 {
        confPath := args[0]
        var err error
        if server, err = mercury.NewServer(confPath); err != nil {
            fmt.Printf("Failed to load configuration file '%s': %s\n", confPath, err)
            return
        }
    } else if _, err := os.Stat(systemConf()); err == nil {
        fmt.Println("No configuration file specified.")
        fmt.Printf("Falling back to %s\n", systemConf())
        if server, err = mercury.NewServer(systemConf()); err != nil {
            fmt.Printf("Failed to load configuration file '%s': %s\n", systemConf(), err)
            return
        }
    } else {
        fmt.Println("No configuration file found in ~/.config/mercury, and none specified")
        return
    }

    // Start the server, logging errors to stdout
    if err := server.ListenAndServe(); err != nil {
        fmt.Println(err)
    }
}

// The default location of mercury-server's configuration file.
// This file will be the fallback if no file was specified.
func systemConf() string {
    return fmt.Sprintf("%s/.config/mercury/server.toml", mercury.HomeDir())
}

