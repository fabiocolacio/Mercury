package server

import(
    "fmt"
    "log"
    "errors"
    "crypto/hmac"
    "crypto/sha256"
    "encoding/json"
    "encoding/base64"
    "golang.org/x/crypto/pbkdf2"
)

var(
    InvalidMACError error = errors.New("Computed MAC does not match the provided MAC")
    MalformedJWTError error = errors.New("The JWT is missing fields or corrupt")
)

// Session represents session data stored in a JWT.
// Session only contains the uid of the currently logged-in user.
type Session struct {
    Uid int
}

// HashAndSaltPassword takes a password and salt, and creates a hash
// of the password concatenated with the salt using pbkdf2.
//
// The number of iterations is defined by KeyHashIterations.
// The hash function used is defined by KeyHashAlgo.
// The length of the key that is created is defined by KeyHashLength.
func HashAndSaltPassword(passwd, salt []byte) []byte {
    return pbkdf2.Key(passwd, salt, KeyHashIterations, KeyHashLength, KeyHashAlgo)
}

// CreateSessionToken creates a JWT token that is sent to the client
// at login to represent a session.
func CreateSessionToken(uid int, macKey []byte) []byte {
    header := base64.URLEncoding.EncodeToString([]byte("{\"alg\": \"HS256\", \"typ\": \"JWT\"}"))
    payload := base64.URLEncoding.EncodeToString([]byte(fmt.Sprintf("{\"Uid\": %d}", uid)))
    jwt := []byte(header + "." + payload)

    mac := hmac.New(sha256.New, macKey)
    mac.Write(jwt)
    tag := mac.Sum(nil)
    encodedTag := make([]byte, base64.URLEncoding.EncodedLen(len(tag)))
    base64.URLEncoding.Encode(encodedTag, tag)

    log.Println(base64.URLEncoding.EncodeToString(macKey))

    jwt = append(jwt, '.')
    jwt = append(jwt, encodedTag...)

    return jwt
}

// UnwrapSessionToken verifies a JWT, and returns its payload if the integrity check passes.
// The session token's payload simply contains the uid of the currently logged in user.
func UnwrapSessionToken(jwt, macKey []byte) (Session, error) {
    var session Session

    separators := make([]int, 0, 2)
    for i := 0; i < len(jwt); i++ {
        if jwt[i] == '.' {
            separators = append(separators, i)
        }
    }

    if len(separators) > 2 {
        return session, MalformedJWTError
    }

    payload := jwt[separators[0] + 1:separators[1]]
    mac := jwt[separators[1] + 1:]

    log.Println(string(jwt[:separators[1]]))

    decodedMAC := make([]byte, base64.URLEncoding.DecodedLen(len(mac)))
    _, err := base64.URLEncoding.Decode(decodedMAC, mac)

    if err != nil {
        log.Println(string(mac))
        log.Println("Decoding Error:", err)
    }

    if !ValidateMAC(jwt[:separators[1]], decodedMAC, macKey) {
        return session, InvalidMACError
    }

    jsonPayload, err := base64.URLEncoding.DecodeString(string(payload))
    if err != nil {
        log.Println("Failed to decode payload:", err)
        return session, MalformedJWTError
    }

    err = json.Unmarshal(jsonPayload, &session)
    if err != nil {
        log.Println("Failed to unmarshal json:", err)
        return session, MalformedJWTError
    }

    return session, nil
}

// ValidateMAC computes a SHA256 HMAC tag for message, and compares it with messageMAC.
// The the tags match, ValidateMAC returns true, else it returns false.
func ValidateMAC(message, messageMAC, key []byte) bool {
    mac := hmac.New(sha256.New, key)
    mac.Write(message)
    expectedMAC := mac.Sum(nil)
    return Memcmp(messageMAC, expectedMAC)
}
