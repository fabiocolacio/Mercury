package mercury

import(
    "fmt"
    "net/http"
    "strings"
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

func NewServerWithConf(conf Config) (Server) {
    return Server{ conf }
}

func NewServer(confPath string) (Server, error) {
    conf, err := LoadConfig(confPath)
    serv := NewServerWithConf(conf)
    return serv, err
}

func LoadConfig(confPath string) (Config, error){
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
        e <- http.ListenAndServeTLS(conf.HttpsAddr, conf.CertFile, conf.KeyFile, http.Handler(serv))
    }()

    return <-e
}

func (serv Server) ServeHTTP(res http.ResponseWriter, req *http.Request) {
    // Redirect Non-HTTPS requests to HTTPS
    if req.TLS == nil {
        host := strings.Split(req.Host, ":")[0]
        port := strings.Split(serv.config.HttpsAddr, ":")[1]
        path := req.URL.Path
        dest := fmt.Sprintf("https://%s:%s%s", host, port, path)
        fmt.Println(dest)
        http.Redirect(res, req, dest, http.StatusTemporaryRedirect)
        return
    }

    res.Write([]byte("<h1>Hello World!</h1>"))
}

