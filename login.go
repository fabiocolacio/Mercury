package main

import(
    "net/http"
    "context"
    "fmt"
    "time"
    "io/ioutil"
    "encoding/json"
    "crypto/sha256"
    "crypto/rand"
    "crypto/hmac"
    "go.mongodb.org/mongo-driver/mongo/options"
)

var(
    ErrNoSuchUser      restError = restError{ 400, "No such user" }
    ErrChallengeFailed restError = restError{ 300, "Challenge failed (request a new one)" }
    ErrAccountLocked   restError = restError{ 423, "Account is locked" }
)

func requestChallengeRoute(res http.ResponseWriter, req *http.Request) {
    username := req.URL.Query().Get("user")

    ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
    defer cancel()
    result := mongoDB.Collection("users").FindOne(ctx, map[string]string{ "user": username }, options.FindOne())
    if result == nil {
        res.WriteHeader(500)
        return
    }

    var user User
    if err := result.Decode(&user); err != nil {
        res.WriteHeader(500)
        return
    }

    challenge := make([]byte, 16)
    if _, err := rand.Read(challenge); err != nil {
        res.WriteHeader(500)
        return
    }

    payload, err := json.Marshal(struct{ C, S []byte }{ challenge, user.Salt })
    if err != nil {
        res.WriteHeader(500)
        fmt.Println(err)
        return
    }


    mac := hmac.New(sha256.New, user.SHash)
    mac.Write(challenge)
    user.Chal = mac.Sum(nil)

    ctx, cancel = context.WithTimeout(context.Background(), 10 * time.Second)
    defer cancel()
    if _, err = mongoDB.Collection("users").UpdateOne(ctx, map[string]string{ "user": username }, map[string]interface{}{ "$set": user }, options.Update()); err != nil {
        res.WriteHeader(500)
        fmt.Println(err)
        return
    } else {
        res.Write(payload)
    }
}

func loginRoute(res http.ResponseWriter, req *http.Request) {
    username := req.URL.Query().Get("user")

    ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
    defer cancel()
    result := mongoDB.Collection("users").FindOne(ctx, map[string]string{ "user": username }, options.FindOne())
    if result == nil {
        res.WriteHeader(500)
        return
    }

    var user User
    if err := result.Decode(&user); err != nil {
        res.WriteHeader(500)
        return
    }

    payload, err := ioutil.ReadAll(req.Body)
    if err != nil {
        res.WriteHeader(500)
        return
    }

    if !hmac.Equal(user.Chal, payload) {
        res.WriteHeader(500)
        return
    }

    res.WriteHeader(200)
}

