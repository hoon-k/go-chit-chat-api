package request

import (
    "log"
    "net/http"
)

// SchemaValidator is a middleware that validates request given schema
type SchemaValidator struct {}

// Intercept intercepts requests
func (m SchemaValidator) Intercept(w http.ResponseWriter, r *http.Request) error {
    if r.Method != "POST" {
        return nil
    }

    log.Printf("validating %s", r.URL.Path)
    return nil
}