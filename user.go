package main

import(
    "github.com/fabiocolacio/golibs/crypto"
    "crypto/rand"
)

const(
    SALT_SIZE     int = 16        // 128-bit salts
)

type User struct {
    User   string `json: "user"  bson: "user"`
    Salt  []byte  `json: "salt"  bson: "salt"`
    SHash []byte  `json: "shash" bson: "shash"`
}

func NewUserFromCreds(creds Credentials) (User, error) {
    var user User

    salt := make([]byte, SALT_SIZE)
    if _, err := rand.Read(salt); err != nil {
        return user, err
    }

    sHash, err := crypto.HashAndSaltPassword([]byte(creds.Pass), salt)
    if err != nil {
        return user, err
    }

    user.User = creds.User
    user.Salt = salt
    user.SHash = sHash

    return user, nil
}

