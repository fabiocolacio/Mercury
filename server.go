package mercury

import(
    "fmt"
    "net/http"
    "github.com/BurntSushi/toml"
)

type Server struct {
    config Config
}

type Config struct {
    HttpAddr  string
    HttpsAddr string
    CertFile  string
    KeyFile   string
}

func New(confPath string) (Config, error){
    var conf Config
    _, err := toml.DecodeFile(confPath, &conf)
    return conf, err
}

func (serv Server) Config() Config {
    return serv.config
}

func (serv Server) ListenAndServe() error {
    conf := serv.config
    e := make(chan error)

    go func() {
        e <- http.ListenAndServe(conf.HttpAddr, http.Handler(serv))
    }()

    go func() {
        e <- http.ListenAndServe(conf.HttpsAddr, http.Handler(serv))
    }()

    return <-e
}

func (serv Server) ServeHTTP(res http.ResponseWriter, req *http.Request) {
    // Redirect Non-HTTPS requests to HTTPS
    if req.TLS == nil {
        host := serv.config.HttpsAddr
        path := req.URL.Path
        dest := fmt.Sprintf("https://%s%s", host, path)
        http.Redirect(res, req, dest, http.StatusTemporaryRedirect)
        return
    }

    res.Write([]byte("<h1>Hello World!</h1>"))
}

