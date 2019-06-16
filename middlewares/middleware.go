package middlewares

import (
    "net/http"
)

// Middleware is the interface
type Middleware interface {
    Intercept(w http.ResponseWriter, r *http.Request) error
}
