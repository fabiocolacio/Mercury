package mercury

import(
    "fmt"
    "errors"
    "hash"
    "crypto/hmac"
    "crypto/sha256"
    "encoding/json"
    "encoding/base64"
    "golang.org/x/crypto/pbkdf2"
)

var(
    InvalidMACError error = errors.New("Computed MAC does not match the provided MAC")
    MalformedJWTError error = errors.New("The JWT is missing fields or corrupt")

    KeyHashAlgo func() hash.Hash = sha256.New
    KeyHashIterations int = 250000
    KeyHashLength int = 32
)

// Credentials represents a userername and password combination
type Credentials struct {
    Username string
    Password string
}

// Session represents session data stored in a JWT.
// Session only contains the uid of the currently logged-in user.
type Session struct {
    Uid string
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
    header := []byte("{\"alg\": \"HS256\", \"typ\": \"JWT\"}")
    payload := []byte(fmt.Sprintf("{\"Uid\": %d}", uid))

    jwt := make([]byte, 0, len(header) + len(payload) + 1)
    jwt = append(jwt, header...)
    jwt = append(jwt, '.')
    jwt = append(jwt, payload...)

    mac := hmac.New(sha256.New, macKey)
    mac.Write(jwt)
    tag := mac.Sum(nil)

    jwt = append(jwt, '.')
    jwt = append(jwt, tag...)

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

    payload := jwt[separators[0]:separators[1]]
    mac := jwt[separators[1]:]

    if !ValidateMAC(jwt[:separators[1]], mac, macKey) {
        return session, InvalidMACError
    }

    jsonPayload, err := base64.URLEncoding.DecodeString(string(payload))
    if err != nil {
        return session, MalformedJWTError
    }

    err = json.Unmarshal(jsonPayload, &session)
    if err != nil {
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
    return hmac.Equal(messageMAC, expectedMAC)
}

