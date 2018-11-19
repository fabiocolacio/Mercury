package server

import(
    "net/http"
    "encoding/json"
    "log"
)

func (serv *Server) TestRoute(res http.ResponseWriter, req *http.Request) {
    body, err := ReadBody(req)

    if err != nil {
        log.Printf("Failed to read request body: %s", err)
        res.WriteHeader(http.StatusUnauthorized)
        return
    }

    _, err = UnwrapSessionToken(body, serv.macKey[:])
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
        res.Write([]byte("Malformed request"))
        return
    }

    err = json.Unmarshal(body, &creds)
    if err != nil {
        log.Println(err)
        res.Write([]byte("ERROR: Invalid JSON object"))
        return
    }

    jwt, err := serv.LoginUser(creds)
    if err != nil {
        res.Write([]byte(err.Error()))
    } else {
        res.Write(jwt)
    }
}
