package server

import(
    "os"
    "fmt"
    "log"
    "net/http"
    "crypto/tls"
    "crypto/rand"
    "strings"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
)

// Server is a type that represents a Mercury Chat Server.
type Server struct {
    httpsServer     *http.Server
    config           Config
    logFile         *os.File
    macKey      [256]byte
    db              *sql.DB
}

// NewServerWithConf creates a new Server structure using the
// settings defined by the Config structure.
func NewServerWithConf(conf Config) (*Server, error) {
    var server  *Server
    var file    *os.File
    var err      error

    // Set the log file as specified in the config
    if conf.LogFile != "" {
        file, err = os.OpenFile(conf.LogFile, os.O_APPEND | os.O_CREATE | os.O_WRONLY, 0644)
        if err == nil {
            log.SetOutput(file)
        }
    }

    // Set default http and https address if none specified in config
    if conf.HttpAddr == "" {
        conf.HttpAddr = DefaultHttpAddr
    }
    if conf.HttpsAddr ==  "" {
        conf.HttpsAddr = DefaultHttpsAddr
    }

    // TLS configuration
    tlsConf := &tls.Config{
        MinVersion: tls.VersionTLS12,
        PreferServerCipherSuites: true,
        CipherSuites: []uint16{ tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384, tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
            tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
            tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
        },
    }

    // Create a custom https server with our TLS config
    httpsServer := &http.Server{
        Addr: conf.HttpsAddr,
        TLSConfig: tlsConf,
    }

    dataSource := fmt.Sprintf("%s:%s@/%s", conf.SQLUser, conf.SQLPass, conf.SQLDb)
    db, err := sql.Open("mysql", dataSource)
    if err != nil {
        return nil, err
    }

    server = &Server {
        httpsServer: httpsServer,
        config:  conf,
        logFile: file,
        db: db,
    }

    rand.Read(server.macKey[:])

    server.httpsServer.Handler = http.Handler(server)

    return server, err
}

// NewServer creates a new Server structure, with the configuration
// specified in a toml configuration file at location confPath.
func NewServer(confPath string) (*Server, error) {
    var(
        serv *Server
        conf  Config
        err   error
        e1    error
        e2    error
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
func (serv *Server) Close() (e error) {
    log.Println("Shutting down server.")

    if serv.logFile != nil {
        e = serv.logFile.Close()
    }

    if serv.db != nil {
        e = serv.db.Close()
    }

    return e
}

// Config returns a copy of the underlying Config structure for a
// particular Server instance.
func (serv *Server) Config() Config {
    return serv.config
}

// ListenAndServe is similar to go's http.ListenAndServe and https.ListenAndServeTLS
// functions. This starts the Mercury Server, and handles incoming connections.
// This is a blocking function, and should be started as a goroutine if it needs
// to run in the background. If it fails to bind one of the sockets, it will return
// with an error.
func (serv *Server) ListenAndServe() error {
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
        e <- serv.httpsServer.ListenAndServeTLS(conf.CertFile, conf.KeyFile)
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
func (serv *Server) ServeHTTP(res http.ResponseWriter, req *http.Request) {
    host := strings.Split(req.Host, ":")[0]
    port := strings.Split(serv.config.HttpsAddr, ":")[1]
    path := req.URL.Path

    // Redirect Non-HTTPS requests to HTTPS
    if req.TLS == nil {
        dest := fmt.Sprintf("https://%s:%s%s", host, port, path)
        log.Printf("Redirecting HTTP client '%s' to %s", req.RemoteAddr, dest)
        http.Redirect(res, req, dest, http.StatusTemporaryRedirect)
        return
    }

    res.Header().Add("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

    log.Printf("Handling client '%s'", req.RemoteAddr)

    switch path {
    case "/register":
        serv.RegisterRoute(res, req)

    case "/login":
        serv.LoginRoute(res, req)

    case "/test":
        serv.TestRoute(res, req)

    case "/send":
        serv.SendRoute(res, req)

    default:
        res.Write([]byte("<h1>Hello World!</h1>"))
    }
}
