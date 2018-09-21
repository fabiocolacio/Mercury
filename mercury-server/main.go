package main

import(
    "flag"
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
    mercury.Assertf(err == nil, "Failed to load configuration file '%s': %s", confPath, err)

    // Start the server, logging errors to stdout
    err = server.ListenAndServe()
    mercury.Assert(err != nil, err)
}

