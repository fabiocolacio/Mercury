package server

import(
    "database/sql"
    "crypto/rand"
    "crypto/subtle"
    "errors"
    "net/http"
)

var(
    ErrInvalidCredentials error = errors.New("Invalid username or password.")
    ErrUsernameTaken error = errors.New("Username already taken.")
    ErrRegistrationFailed error = errors.New("Failed to register user.")
    ErrNotLoggedIn error = errors.New("Please log in.")
)

// Credentials represents a userername and password combination
type Credentials struct {
    Username string
    Password string
}

// LoginUser creates a JWT session token if the credentials creds are valid
func (serv *Server) LoginUser(creds Credentials) (jwt []byte, err error) {
    row := serv.db.QueryRow(
        "select uid, salt, saltedhash from users where username = ?;",
        creds.Username)

    var(
        uid int
        salt []byte
        saltedHash []byte
    )

    if row.Scan(&uid, &salt, &saltedHash) == sql.ErrNoRows {
        err = ErrInvalidCredentials
    } else {
        key := HashAndSaltPassword([]byte(creds.Password), salt)

        if subtle.ConstantTimeCompare(key, saltedHash) != 1{
            err = ErrInvalidCredentials
            return
        }

        jwt = CreateSessionToken(uid, serv.macKey[:])
    }

    return
}

// RegisterUser attempts to creates a user in the with the credentials creds
func (serv *Server) RegisterUser(creds Credentials) (err error) {
    row := serv.db.QueryRow("select uid from users where username = ?;", creds.Username)

    var uid int
    if row.Scan(&uid) == sql.ErrNoRows {
        salt := make([]byte, SaltLength)
        rand.Read(salt)

        saltedHash := HashAndSaltPassword([]byte(creds.Password), salt)

        _, err := serv.db.Exec(
            `insert into users (username, salt, saltedhash) values (?, ?, ?);`,
            creds.Username, salt, saltedHash)

        if err != nil {
            err = ErrRegistrationFailed
        }
    } else {
        err = ErrUsernameTaken
    }

    return
}

// VerifyUser extracts the Session data from a request.
func (serv *Server) VerifyUser(req *http.Request) (Session, error) {
    jwt := []byte(req.Header.Get("Session"))
    if len(jwt) <= 0 {
        return Session{}, ErrNotLoggedIn
    }
    return UnwrapSessionToken(jwt, serv.macKey[:])
}
