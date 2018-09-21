package mercury

import(
    "github.com/BurntSushi/toml"
)

// Config contains configuration data for use by Server.
type Config struct {
    HttpAddr  string
    HttpsAddr string
    CertFile  string
    KeyFile   string
    LogFile   string
}

// LoadConfig loads a toml-formatted configuration file at the location
// confPath, and returns a new Config structure to represent it.
func LoadConfig(confPath string) (Config, error){
    var conf Config
    _, err := toml.DecodeFile(confPath, &conf)
    return conf, err
}

