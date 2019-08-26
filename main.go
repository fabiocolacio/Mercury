package main

import(
    "strings"
    "fmt"
    "net/http"
    "crypto/tls"
    "crypto/rand"
    "database/sql"
    "flag"
    "log"

    "github.com/go-chi/chi"
    "github.com/go-chi/chi/middleware"
    "github.com/go-chi/jwtauth"
    _ "github.com/go-sql-driver/mysql"
)

const(
    HMAC_KEY_SIZE  int = 32        // 256-bit HMAC key
)

var(
    sqlDb *sql.DB
)

func main() {
    var (
        privKeyFile string
        certFile    string
        httpAddr    string
        httpsAddr   string
        sqlUser     string
        sqlPass     string
        sqlDb       string
    )

    // Parse command-line arguments
    flag.StringVar(&certFile, "cert", "", "The location of your ssl certificate")
    flag.StringVar(&privKeyFile, "privkey", "", "The location of your ssl private key")
    flag.StringVar(&httpAddr, "http-addr", ":8080", "The address from which to listen to http requests")
    flag.StringVar(&httpsAddr, "https-addr", ":9090", "The address from which to listen to https requests")
    flag.StringVar(&sqlUser, "sql-user", "", "The SQL user to use")
    flag.StringVar(&sqlPass, "sql-pass", "", "The SQL password to use")
    flag.StringVar(&sqlDb, "sql-db", "mercury", "The name of the SQL database to use")
    flag.Parse()

    dataSource := fmt.Sprintf("%s:%s@/%s", sqlUser, sqlPass, sqlDb)
    db, err := sql.Open("mysql", dataSource)
    if err != nil {
        log.Fatal("Failed to open SQL database!", err)
    }

    if err = initDB(db); err != nil {
        log.Fatal("Error creating database tables!", err)
    }

    if certFile == "" {
        log.Fatal("No certificate specified")
    }
    if privKeyFile == "" {
        log.Fatal("No private key specified")
    }

    // Create a random HMAC key
    hmacKey := make([]byte, HMAC_KEY_SIZE)
    if _, err := rand.Read(hmacKey); err != nil {
        log.Fatal(err)
    }

    // Create a JWT authenticator
    tokenAuth := jwtauth.New("H256", hmacKey, hmacKey)

    // HTTP routing
    router := chi.NewRouter()
    router.Use(middleware.Logger)
    router.Use(middleware.SetHeader("Strict-Transport-Security", "max-age=31536000; includeSubDomains"))
    router.Post("/register", registerRoute)
    // router.Get("/login", requestChallengeRoute)
    // router.Post("/login", loginRoute)

    // Protected Routes (must be logged in)
    router.Group(func(prouter chi.Router) {
        prouter.Use(jwtauth.Verifier(tokenAuth))
        prouter.Use(jwtauth.Authenticator)
    })

    httpsServer := &http.Server {
        Addr: httpsAddr,
        Handler: router,
        TLSConfig: &tls.Config {
            MinVersion: tls.VersionTLS12,
            PreferServerCipherSuites: true,
            CipherSuites: []uint16 {
                tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
                tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
                tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
                tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
            },
        },
    }

    // Redirects http traffic to https
    tlsRedirectHandler := func (res http.ResponseWriter, req *http.Request) {
        host := strings.Split(req.Host, ":")[0]
        port := strings.Split(httpsAddr, ":")[1]
        path := req.URL.Path
        dest := fmt.Sprintf("https://%s:%s%s", host, port, path)
        log.Printf("Redirecting HTTP client '%s' to %s", req.RemoteAddr, dest)
    }

    e := make(chan error)

    // Listen for HTTP requests concurrently
    go func() {
        e <- http.ListenAndServe(httpAddr, http.HandlerFunc(tlsRedirectHandler))
    }()

    // Listen for HTTPS requests concurrently
    go func() {
        e <- httpsServer.ListenAndServeTLS(certFile, privKeyFile)
    }()

    // Wait for one of the servers to fail or close
    log.Println(<-e)
}

