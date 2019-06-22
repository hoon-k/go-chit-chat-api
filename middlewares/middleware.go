package middlewares

import (
    "log"
    "net/http"

    "github.com/julienschmidt/httprouter"
)

// Middleware is the interface
type Middleware interface {
    Intercept(w http.ResponseWriter, r *http.Request) error
}
// ManagedRouter is a router with middleware
type ManagedRouter struct {
    middlewares []Middleware
    router *httprouter.Router
}

// CreateManagedRouter creates the router to be used for the middlewares
func CreateManagedRouter(router *httprouter.Router) *ManagedRouter {
    return &ManagedRouter {
        router: router,
    }
}

// Add to middlewares
func (mr *ManagedRouter) Add(middleware Middleware) {
    mr.middlewares = append(mr.middlewares, middleware)
}

// ServeHTTP serves http
func (mr *ManagedRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    log.Printf("Serving %s", r.URL.Path)

    for _, mw := range mr.middlewares {
        err := mw.Intercept(w, r)
        if err != nil {
            log.Fatalf("%s: %s", "Failed on middleware", err)
        }
    }

    header := w.Header()
    header.Add("Content-Type", "application/json")
    header.Add("Access-Control-Allow-Origin", "*")
    header.Add("Access-Control-Allow-Headers", "origin")
    mr.router.ServeHTTP(w, r)
}