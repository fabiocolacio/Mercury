package main

import(
    "fmt"
    "flag"
    "io"
    "io/ioutil"
    "crypto/rsa"
    "crypto/aes"
    "crypto/hmac"
    "crypto/cipher"
    "crypto/x509"
    "crypto/rand"
    "crypto/sha256"
    "encoding/pem"
    "github.com/fabiocolacio/mercury"
)

var(
    pubKey string
    message string
)

func init() {
    flag.StringVar(&pubKey, "k", "", "The public key to encrypt the message with")
    flag.StringVar(&message, "m", "thisismymessage!", "The message to encrypt")
    flag.Parse()

    mercury.Assert(pubKey != "", "Please specify a key to encrypt the message with!")
    mercury.Assert(message != "", "Please specify a message to encrypt!")
}

func main() {
    fmt.Printf("Plain Text:\n%s\n", message)

    // TODO: The plaintext message needs to be a length that is a multiple of
    // the AES block length (16 bytes). Padding needs to be added to the message
    // if it is not!

    // Load RSA public key into the 'rsaKey' variable
    data, err := ioutil.ReadFile(pubKey)
    mercury.Assertf(err == nil, "Failed to read file '%s': %s", pubKey, err)
    block, _ := pem.Decode(data)
    rsaKey, err := x509.ParsePKIXPublicKey(block.Bytes)
    mercury.Assertf(err == nil, "Unable to parse public key: %s", err)

    // Create a random 256-bit AES key in the variable 'aesKey'
    // crypto/rand produces cryptographically secure random values
    // For implementation details of crypto/rand, see https://golang.org/pkg/crypto/rand/
    aesKey := make([]byte, 32)
    _, err = rand.Read(aesKey)

    // Create a random 256-bit HMAC key in the variable 'hmacKey'
    hmacKey := make([]byte, 32)
    _, err = rand.Read(hmacKey)

    // Encrypt the message with our AES key into 'cipherText'
    blockCipher, err := aes.NewCipher(aesKey)
    cipherText := make([]byte, aes.BlockSize + len(message))
    iv := cipherText[:aes.BlockSize]
    io.ReadFull(rand.Reader, iv)
    encrypter := cipher.NewCBCEncrypter(blockCipher, iv)
    encrypter.CryptBlocks(cipherText[aes.BlockSize:], []byte(message))
    fmt.Printf("Cipher Text:\n%x\n\n", cipherText)

    // Create an HMAC tag in the variable 'hash'
    var hash []byte
    hash = hmac.New(sha256.New, hmacKey).Sum(hash)
    fmt.Printf("HMAC tag:\n%x\n\n", hash)

    // Concatenate AES and HMAC keys into the variable 'conKey'
    conKey := make([]byte, 64)
    copy(conKey, aesKey)
    copy(conKey[32:], hmacKey)

    // Encrypt the keys with the RSA public key into the variable 'encryptedKeys'
    encryptedKeys, err := rsa.EncryptOAEP(
        sha256.New(),
        rand.Reader,
        rsaKey.(*rsa.PublicKey),
        conKey,
        []byte("message"))
    fmt.Printf("Encrypted Keys:\n%x\n\n", encryptedKeys)

    // TODO: Format this into a JSON message
}

