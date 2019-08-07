package main

import(
    "net/http"
)

type restError struct {
    Status  int
    Message string
}

var(
    ErrInternalServer   restError = restError{ 500, "Internal server error" }
    ErrMalformedMessage restError = restError{ 400, "Malformed message" }
    ErrForbidden        restError = restError{ 403, "Forbidden" }
    ErrInvalidRequest   restError = restError{ 404, "Invalid request" }
)

func (err *restError) SendResponse(res http.ResponseWriter) {
    res.WriteHeader(err.Status)
    res.Write([]byte(err.Message))
}

func (err *restError) Error() string {
    return err.Message
}

