package main

import(
    "net/http"
    "io/ioutil"
    "encoding/json"
    "crypto/rand"
    "log"

    "golang.org/x/crypto/scrypt"
)

const(
    SALT_SIZE      int = 16
    SHASH_SIZE     int = 32
    CHALLENGE_SIZE int = 16
)

var(
    ErrNoPassword restError = restError{ 400, "No password specified" }
)

func registerRoute(res http.ResponseWriter, req *http.Request) {
    payload, err := ioutil.ReadAll(req.Body)
    if err != nil {
        res.WriteHeader(500)
        log.Println(err)
        return
    }

    var body map[string]string
    if err := json.Unmarshal(payload, &body); err != nil {
        ErrMalformedMessage.SendResponse(res)
        log.Println(err)
        return
    }

    pass, ok := body["password"]
    if !ok {
        ErrNoPassword.SendResponse(res)
        log.Println(err)
        return
    }

    salt := make([]byte, SALT_SIZE)
    if _, err := rand.Read(salt); err != nil {
        ErrInternalServer.SendResponse(res)
        log.Println(err)
        return
    }

    shash, err := scrypt.Key([]byte(pass), salt, 32768, 8, 1, SHASH_SIZE)
    if err != nil {
        ErrInternalServer.SendResponse(res)
        log.Println(err)
        return
    }

    result, err := sqlDb.Exec(`insert into users (salt, shash) values (?, ?)`, salt, shash)
    if err != nil {
        ErrInternalServer.SendResponse(res)
        log.Println(err)
        return
    }

    uid, err := result.LastInsertId()
    if err != nil {
        ErrInternalServer.SendResponse(res)
        log.Println(err)
        return
    }

    response, err := json.Marshal(map[string]int64{ "uid": uid })
    if err != nil {
        ErrInternalServer.SendResponse(res)
        log.Println(err)
        return
    }

    res.Write(response)
}

