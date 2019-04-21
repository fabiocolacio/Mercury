package main

import(
    "github.com/go-chi/chi"
    "github.com/go-chi/chi/middleware"
    "github.com/go-chi/jwtauth"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/bson"
    "context"
    "strings"
    "time"
    "fmt"
    "encoding/json"
    "net/http"
    "io/ioutil"
    "crypto/tls"
    "crypto/rand"
    "flag"
    "log"
)

var(
    mongoDB      *mongo.Database  // Our database
    hmacKey     []byte            // HMAC Key used to sign messages
    tokenAuth    *jwtauth.JWTAuth // JWT authenticator
    certFile      string          // Path to SSL certificate file
    privKeyFile   string          // Path to SSL private key file
    httpAddr      string          // Address and port to listen for HTTP requests from
    httpsAddr     string          // Address and port to listen for HTTPS requests from 
)

const(
    HMAC_KEY_SIZE  int = 32        // 256-bit HMAC key
)

func main() {
    // Parse command-line arguments
    flag.StringVar(&certFile, "cert", "", "The location of your ssl certificate")
    flag.StringVar(&privKeyFile, "privkey", "", "The location of your ssl private key")
    flag.StringVar(&httpAddr, "http-addr", ":8080", "The address from which to listen to http requests")
    flag.StringVar(&httpsAddr, "https-addr", ":9090", "The address from which to listen to https requests")
    flag.Parse()

    if certFile == "" {
        log.Fatal("No certificate specified")
    }
    if privKeyFile == "" {
        log.Fatal("No private key specified")
    }

    // Create a random HMAC key
    hmacKey = make([]byte, HMAC_KEY_SIZE)
    if _, err := rand.Read(hmacKey); err != nil {
        log.Fatal(err)
    }

    // Create a JWT authenticator
    tokenAuth = jwtauth.New("H256", hmacKey, hmacKey)

    // Create a connection to MongoDB
    mongoCon, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost"))
    if err != nil {
        log.Fatal(err)
    }

    ctx, cancel := context.WithTimeout(context.Background(), 15 * time.Second)
    err = mongoCon.Connect(ctx)
    if err != nil {
        log.Fatal(err)
    }
    cancel()

    mongoDB = mongoCon.Database("mercury")

    ctx, cancel = context.WithTimeout(context.Background(), 15 * time.Second)
    mongoDB.Collection("users").Indexes().CreateOne(
        context.Background(),
        mongo.IndexModel{ bson.M{ "user": 1 }, options.Index().SetUnique(true)},
        options.CreateIndexes(),
    )
    cancel()

    router := chi.NewRouter()

    router.Use(middleware.Logger)
    router.Use(middleware.SetHeader("Strict-Transport-Security", "max-age=31536000; includeSubDomains"))

    router.Post("/register", registerRoute)

    // Protected Routes
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

func tlsRedirectHandler(res http.ResponseWriter, req *http.Request) {
    host := strings.Split(req.Host, ":")[0]
    port := strings.Split(httpsAddr, ":")[1]
    path := req.URL.Path

    dest := fmt.Sprintf("https://%s:%s%s", host, port, path)
    log.Printf("Redirecting HTTP client '%s' to %s", req.RemoteAddr, dest)
}

func registerRoute(res http.ResponseWriter, req *http.Request) {
    body, err := ioutil.ReadAll(req.Body)
    if err != nil {
        res.WriteHeader(500)
        log.Println(err)
        return
    }
    
    var creds Credentials
    if err := json.Unmarshal(body, &creds); err != nil {
        res.WriteHeader(400)
        log.Println(err)
        return
    }

    user, err := NewUserFromCreds(creds)
    if err != nil {
        res.WriteHeader(500)
        log.Println(err)
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
    defer cancel()
    mongoDB.Collection("users").InsertOne(ctx, user, options.InsertOne())
}

