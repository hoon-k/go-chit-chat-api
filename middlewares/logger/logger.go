package logger

import (
    "log"
    "net/http"
)

// Logger is a middleware that logs things
type Logger struct {}

// Intercept intercepts requests
func (m Logger) Intercept(w http.ResponseWriter, r *http.Request) error {
    log.Printf("Logging %s", r.URL.Path)
    return nil
}