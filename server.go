package mercury

import(
    "os"
    "fmt"
    "log"
    "net/http"
    "strings"
)

// Server is a type that represents a Mercury Chat Server.
type Server struct {
    config Config
    logFile *os.File
}

// NewServerWithConf creates a new Server structure using the
// settings defined by the Config structure.
func NewServerWithConf(conf Config) (Server, error) {
    var file  *os.File
    var err    error

    if conf.LogFile != "" {
        file, err = os.OpenFile(conf.LogFile, os.O_APPEND | os.O_CREATE | os.O_WRONLY, 0644)
        if err == nil {
            log.SetOutput(file)
        }
    }

    return Server{
        config:  conf,
        logFile: file }, err
}

// NewServer creates a new Server structure, with the configuration
// specified in a toml configuration file at location confPath.
func NewServer(confPath string) (Server, error) {
    var(
        serv Server
        conf Config
        err  error
        e1   error
        e2   error
    )

    conf, e1 = LoadConfig(confPath)
    serv, e2 = NewServerWithConf(conf)

    if e1 != nil {
        err = e1
    } else {
        err = e2
    }

    return serv, err
}

// Closes resources allocated by the server
func (serv Server) Close() (e error) {
    log.Println("Shutting down server.")
    if serv.logFile != nil {
        e = serv.logFile.Close()
    }
    return e
}

// Config returns a copy of the underlying Config structure for a
// particular Server instance.
func (serv Server) Config() Config {
    return serv.config
}

// ListenAndServe is similar to go's http.ListenAndServe and https.ListenAndServeTLS
// functions. This starts the Mercury Server, and handles incoming connections.
// This is a blocking function, and should be started as a goroutine if it needs
// to run in the background. If it fails to bind one of the sockets, it will return
// with an error.
func (serv Server) ListenAndServe() error {
    conf := serv.config
    e := make(chan error)

    log.Printf("Listening to HTTP requests on %s", conf.HttpAddr)
    log.Printf("Listening to HTTPS requests on %s", conf.HttpsAddr)

    // HTTP and TLS servers are bound to the socket addresses defined in the
    // Config structure, and respond to requests concurrently.
    // The responses are generated in the ServeHTTP function below.
    go func() {
        e <- http.ListenAndServe(conf.HttpAddr, http.Handler(serv))
    }()
    go func() {
        e <- http.ListenAndServeTLS(conf.HttpsAddr, conf.CertFile, conf.KeyFile, http.Handler(serv))
    }()

    // Block execution until one of the functions returns with a critical error.
    // This may fail if you are trying to bind to a port that is in use, or if
    // you do not have proper permissions to bind to that port.
    err := <-e

    if err != nil {
        log.Println(err)
    }

    return err
}

// ServeHTTP generates an HTTP response to an HTTP request. See the go
// http.Handler interface for more information.
func (serv Server) ServeHTTP(res http.ResponseWriter, req *http.Request) {
    // Redirect Non-HTTPS requests to HTTPS
    if req.TLS == nil {
        host := strings.Split(req.Host, ":")[0]
        port := strings.Split(serv.config.HttpsAddr, ":")[1]
        path := req.URL.Path
        dest := fmt.Sprintf("https://%s:%s%s", host, port, path)

        log.Printf("Redirecting HTTP client '%s' to %s", req.RemoteAddr, dest)
        http.Redirect(res, req, dest, http.StatusTemporaryRedirect)
        return
    }

    log.Printf("Handling client '%s'", req.RemoteAddr)
    res.Write([]byte("<h1>Hello World!</h1>"))
}

