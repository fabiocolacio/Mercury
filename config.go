package mercury

import(
    "log"
    "github.com/BurntSushi/toml"
)

// Config contains configuration data for use by Server.
type Config struct {
    HttpAddr  string
    HttpsAddr string
    CertFile  string
    KeyFile   string
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

