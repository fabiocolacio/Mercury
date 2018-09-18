package mercury

import(
    "fmt"
    "log"
    "net/http"
    "strings"
    "github.com/BurntSushi/toml"
)

// Server is a type that represents a Mercury Chat Server.
type Server struct {
    config Config
}

// Config contains configuration data for use by Server.
type Config struct {
    HttpAddr  string
    HttpsAddr string
    CertFile  string
    KeyFile   string
}

// NewServerWithConf creates a new Server structure using the
// settings defined by the Config structure.
func NewServerWithConf(conf Config) (Server) {
    return Server{ conf }
}

// NewServer creates a new Server structure, with the configuration
// specified in a toml configuration file at location confPath.
func NewServer(confPath string) (Server, error) {
    conf, err := LoadConfig(confPath)
    serv := NewServerWithConf(conf)
    return serv, err
}

// LoadConfig loads a toml-formatted configuration file at the location
// confPath, and returns a new Config structure to represent it.
func LoadConfig(confPath string) (Config, error){
    var conf Config
    _, err := toml.DecodeFile(confPath, &conf)

    if err == nil {
        log.Printf("Loaded configuration file '%s' successfully.", confPath)
    } else {
        log.Printf("Failed to load file '%s': %s", confPath, err)
    }

    return conf, err
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

