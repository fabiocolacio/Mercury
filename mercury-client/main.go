package main

import(
    "fmt"
    "flag"
    "io/ioutil"
    "github.com/fabiocolacio/mercury"
)

var(
    keyFile string
    message string
    decrypt bool
)

func init() {
    flag.StringVar(&keyFile, "k", "", "The RSA key to encrypt the message with")
    flag.StringVar(&message, "m", "thisismymessage!", "The message to encrypt/decrypt")
    flag.BoolVar(&decrypt, "d", false, "Decrypt the message")
    flag.Parse()

    mercury.Assert(keyFile != "", "Please specify a key to encrypt the message with!")
    mercury.Assert(message != "", "Please specify a message to encrypt!")
}

func main() {
    // Load RSA key into buffer
    keyData, _ := ioutil.ReadFile(keyFile)

    if decrypt {
        // Decrypt the JSON message into plaintext
        plaintext, err := mercury.Decrypt(keyData, []byte(message))
        if err == nil {
            fmt.Println(string(plaintext))
        } else {
            fmt.Println(err)
        }
    } else {
        // Encrypt the message into JSON format
        json, err := mercury.Encrypt(keyData, []byte(message))
        if err == nil {
            fmt.Println(string(json))
        } else {
            fmt.Println(err)
        }
    }
}

