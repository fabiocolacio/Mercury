package mercury

import(
    "os"
    "runtime"
)

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

