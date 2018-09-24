package mercury

import(
    "io"
    "errors"
    "bytes"
    "crypto/rsa"
    "crypto/aes"
    "crypto/hmac"
    "crypto/cipher"
    "crypto/x509"
    "crypto/rand"
    "crypto/sha256"
    "encoding/pem"
    "encoding/json"
)

type JSONMessage struct {
    Key []byte
    Tag []byte
    Msg []byte
}

var rsaTag []byte

func init() {
    rsaTag = []byte("message")
}

// Encrypt encrypts a message and returns a JSON object.
//
// The algorithm is as follows:
// 1. Generate a random AES key
// 2. Encrypt the plaintext with the AES key
// 3. Generate a random HMAC key
// 4. Create an HMAC tag using the ciphertext and HMAC key
// 5. Concatenate the AES and HMAC keys
// 6. Encrypt the two keys with the recipient's public RSA key
// 7. Send a JSON message to the recipient like this
//
// {
//     Key: "The encrypted keys go here."
//     Tag: "The HMAC tag goes here.",
//     Msg: "the encrypted message here",
// }
//
// The recipient can decrypt the message using Decrypt
//
// key is the recipient's public key.
// plaintext is the message to encrypt
func Encrypt(key, plaintext []byte) ([]byte, error) {
    // Create padding if message isn't a multiple of 16
    if offset := len(plaintext) % aes.BlockSize; offset != 0 {
        padding := make([]byte, aes.BlockSize - offset)
        plaintext = append(plaintext, padding...)
    }

    // Decode the key data
    pemData, _ := pem.Decode(key)
    rsaKey, err := x509.ParsePKIXPublicKey(pemData.Bytes)
    if err != nil {
        return nil, err
    }

    // Make a random AES key
    aesKey := make([]byte, 32)
    _, err = rand.Read(aesKey)
    if err != nil {
        return nil, err
    }

    // Make a random HMAC key
    hmacKey := make([]byte, 32)
    _, err = rand.Read(hmacKey)
    if err != nil {
        return nil, err
    }

    // Create an AES structure with our key
    block, err := aes.NewCipher(aesKey)
    if err != nil {
        return nil, err
    }

    // Create a buffer for the ciphertext
    ciphertext := make([]byte, aes.BlockSize + len(plaintext))

    iv := ciphertext[:aes.BlockSize]
    _, err = io.ReadFull(rand.Reader, iv)
    if err != nil {
        return nil, err
    }

    // Encrypt the message with a CBC encrypter
    encrypter := cipher.NewCBCEncrypter(block, iv)
    encrypter.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)

    // Create an HMAC tag
    tag := make([]byte, 0, 32)
    tag = hmac.New(sha256.New, hmacKey).Sum(tag)

    // Concatenate AES and HMAC keys
    keys := make([]byte, 64)
    copy(keys, aesKey)
    copy(keys[32:], hmacKey)

    // Encrypt the keys
    keys, err = rsa.EncryptOAEP(
        sha256.New(),
        rand.Reader,
        rsaKey.(*rsa.PublicKey),
        keys,
        rsaTag)
    if err != nil {
        return nil, err
    }

    // Encode as JSON object
    message := JSONMessage{
        Key: keys,
        Tag: tag,
        Msg: ciphertext,
    }
    jsonMessage, err := json.Marshal(message)
    if err != nil {
        return nil, err
    }

    return jsonMessage, nil
}

// Decrypt decrypts a message encrypted with Encrypt
//
// The algorithm is as follows:
// 1. Decrypt the keys with your private RSA key
// 2. Extract the AES key (first 256 bits) from the keys
// 3. Extract the HMAC key (last 256 bits) from the keys
// 4. Run HMAC with the HMAC key, and compare the tags
// 5. Decrypt the message with the AES key
//
// key is your private key.
// json is the JSON-encoded message produced by Encrypt
func Decrypt(key, jsonData []byte) ([]byte, error) {
    // Decode the JSON message
    var message JSONMessage
    err := json.Unmarshal(jsonData, &message)
    if err != nil {
        return nil, err
    }

    // Decode the RSA key data
    pemData, _ := pem.Decode(key)
    rsaKey, err := x509.ParsePKCS1PrivateKey(pemData.Bytes)
    if err != nil {
        return nil, err
    }

    // Decrypt the keys
    message.Key, err = rsa.DecryptOAEP(
        sha256.New(),
        rand.Reader,
        rsaKey,
        message.Key,
        rsaTag)
    if err != nil {
        return nil, err
    }

    // Extract AES and HMAC keys
    aesKey := message.Key[:32]
    hmacKey := message.Key[32:]

    // Verify HMAC tag with HMAC key
    tag := make([]byte, 0, 32)
    tag = hmac.New(sha256.New, hmacKey).Sum(tag)
    if bytes.Compare(tag, message.Tag) != 0 {
        return nil, errors.New("HMAC tags don't match.")
    }

    // Create AES structure with AES key
    block, err := aes.NewCipher(aesKey)
    if err != nil {
        return nil, err
    }

    // Decrypt the message with AES key
    iv := message.Msg[:aes.BlockSize]
    message.Msg = message.Msg[aes.BlockSize:]
    decrypter := cipher.NewCBCDecrypter(block, iv)
    decrypter.CryptBlocks(message.Msg, message.Msg)

    return message.Msg, nil
}

