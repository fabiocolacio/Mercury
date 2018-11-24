package server

import(
    "net/http"
    "encoding/json"
    "log"
    "fmt"
)

func (serv *Server) LookupRoute(res http.ResponseWriter, req *http.Request) {
    query := req.URL.Query()
    user := query.Get("user")
    if len(user) < 1 {
        res.WriteHeader(400)
        res.Write([]byte("No user specified"))
        return
    }

    uid, err := serv.LookupUser(user)
    if err != nil {
        res.WriteHeader(400)
        res.Write([]byte(err.Error()))
        return
    }

    fmt.Fprint(res, uid)
}

func (serv *Server) GetRoute(res http.ResponseWriter, req *http.Request) {
    sess, err := serv.VerifyUser(req)
    if err != nil {
        log.Printf("Unauthorized request: %s", err)
        res.WriteHeader(http.StatusUnauthorized)
        return
    }

    query := req.URL.Query()
    peer := query.Get("peer")
    since := query.Get("since")

    if len(peer) < 1 {
        res.Write([]byte("No peer specified"))
        return
    }

    messages, err := serv.MsgFetch(peer, sess.Uid, since)
    if err != nil {
        res.Write([]byte("Failed to fetch messages"))
        return
    }

    res.Write(messages)
}

func (serv *Server) SendRoute(res http.ResponseWriter, req *http.Request) {
    sess, err := serv.VerifyUser(req)
    if err != nil {
        log.Printf("Unauthorized request: %s", err)
        res.WriteHeader(http.StatusUnauthorized)
        return
    }

    recipient := req.URL.Query().Get("to")
    if len(recipient) < 1 {
        res.Write([]byte("No recipient specified"))
        return
    }

    message, _ := ReadBody(req)

    err = serv.SendMsg(message, recipient, sess.Uid)
    if err != nil {
        res.Write([]byte(err.Error()))
        return
    }

    res.WriteHeader(http.StatusOK)
}

func (serv *Server) TestRoute(res http.ResponseWriter, req *http.Request) {
    _, err := serv.VerifyUser(req)
    if err != nil {
        log.Printf("Unauthorized request: %s", err)
        res.WriteHeader(http.StatusUnauthorized)
        return
    }

    res.WriteHeader(http.StatusOK)
}

func (serv *Server) RegisterRoute(res http.ResponseWriter, req *http.Request) {
    res.Header().Set("Content-Type", "text/plain; charset=utf-8")
    if req.ContentLength > 0 {
        var(
            creds  Credentials
            err    error
        )

        body, err := ReadBody(req)
        if err != nil {
            res.Write([]byte("Malformed request"))
            return
        }

        err = json.Unmarshal(body, &creds)
        if err != nil {
            log.Println(err)
            res.Write([]byte("ERROR: Invalid JSON object"))
            return
        }

        err = serv.RegisterUser(creds)
        if err == nil {
            res.WriteHeader(http.StatusOK)
        } else {
            res.Write([]byte(err.Error()))
        }
    }
}

func (serv *Server) LoginRoute(res http.ResponseWriter, req *http.Request) {
    res.Header().Set("Content-Type", "text/plain; charset=utf-8")

    var creds Credentials

    body, err := ReadBody(req)
    if err != nil {
        res.WriteHeader(400)
        res.Write([]byte("Malformed request"))
        return
    }

    err = json.Unmarshal(body, &creds)
    if err != nil {
        log.Println(err)
        res.WriteHeader(400)
        res.Write([]byte("ERROR: Invalid JSON object"))
        return
    }

    jwt, err := serv.LoginUser(creds)
    if err != nil {
        res.WriteHeader(400)
        res.Write([]byte(err.Error()))
    } else {
        res.Write(jwt)
    }
}
