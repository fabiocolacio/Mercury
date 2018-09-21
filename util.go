package mercury

import(
    "log"
)

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

