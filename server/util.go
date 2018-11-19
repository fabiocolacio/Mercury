package server

import(
    "log"
    "net/http"
)

// Memcmp returns true if the first n bytes of two slices are
// equal and false otherwise, where n is Min(len(a), len(b)).
func Memcmp(a, b []byte) bool {
    var n int
    if len(a) < len(b) {
        n = len(a)
    } else {
        n = len(b)
    }

    for i := 0; i < n; i++ {
        if a[i] != b[i] {
            return false
        }
    }

    return true
}

// Assert checks if condition is false, and exits the program if it is.
// Assert logs message to the standard logger before exiting if the
// assertion fails.
func Assert(condition bool, message interface{}) {
    if !condition {
        log.Fatal(message)
    }
}

// Assertf is equivalent to Assert, but with a format string
func Assertf(condition bool, format string, args ...interface{}) {
    if !condition {
        log.Fatalf(format, args)
    }
}

// ReadBody reads the body of an http.Request into a []byte
func ReadBody(req *http.Request) (body []byte, err error) {
    body = make([]byte, req.ContentLength)

    read, err := req.Body.Read(body);

    if int64(read) == req.ContentLength {
        err = nil
    }

    return body, err
}

