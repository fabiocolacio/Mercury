package mercury

import(
    "os"
    "log"
    "runtime"
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

// HomeDir gets the home directory for the current user.
// Pulled from: https://github.com/golang/go/blob/go1.8rc2/src/go/build/build.go#L260-L277
func HomeDir() string {
    env := "HOME"
    if runtime.GOOS == "windows" {
        env = "USERPROFILE"
    } else if runtime.GOOS == "plan9" {
        env = "home"
    }
    return os.Getenv(env)
}

