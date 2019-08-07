package main

import(
    "net/http"
    "io/ioutil"
    "encoding/json"
    "context"
    "time"
    "log"
    "go.mongodb.org/mongo-driver/mongo/options"
)

var(
    ErrUsernameTaken   restError = restError{ 400, "Username taken" }
    ErrInvalidUsername restError = restError{ 400, "Invalid username" }
)

func registerRoute(res http.ResponseWriter, req *http.Request) {
    body, err := ioutil.ReadAll(req.Body)
    if err != nil {
        res.WriteHeader(500)
        log.Println(err)
        return
    }

    var creds Credentials
    if err := json.Unmarshal(body, &creds); err != nil {
        res.WriteHeader(400)
        log.Println(err)
        return
    }

    user, err := NewUserFromCreds(creds)
    if err != nil {
        res.WriteHeader(500)
        log.Println(err)
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
    defer cancel()
    mongoDB.Collection("users").InsertOne(ctx, user, options.InsertOne())
}

