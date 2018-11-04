package mercury

import(
    "errors"
    "crypto/hmac"
    "crypto/sha256"
    "encoding/json"
    "encoding/base64"
)

var InvalidMACError error = errors.New("Computed MAC does not match the provided MAC")
var MalformedJWTError error = errors.New("The JWT is missing fields or corrupt")

// Credentials represents a userername and password combination
type Credentials struct {
    Username string
    Password string
}

// Session represents session data stored in a JWT.
// Session only contains the uid of the currently logged-in user.
type Session struct {
    uid string
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

