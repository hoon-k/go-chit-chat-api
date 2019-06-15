package main

import (
    "database/sql"
    "fmt"
    "log"
    "net/http"

    "github.com/julienschmidt/httprouter"
    _ "github.com/lib/pq"
)

type myMiddleware interface {
    intercept(w http.ResponseWriter, r *http.Request) error
}

type myrouter struct {
    middlewares []myMiddleware
    router *httprouter.Router
}

type mw1 struct {}
type mw2 struct {}

// Index function
func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    db := getDBConnection()
    defer db.Close()

    rows, err := db.Query("SELECT first_name, last_name FROM users")
    failOnError(err, "Query execution failure")
    defer rows.Close()

    var firstName string
    var lastName string

    for rows.Next() {
        err := rows.Scan(&firstName, &lastName)
        if err != nil {
            log.Fatal(err)
        }

        log.Println(firstName, lastName)
        fmt.Fprintf(w, "Name is %s %s\n", firstName, lastName)
    }
}

func main() {
    router := httprouter.New()
    router.GET("/", Index)
    router.POST("/account/create", create)
    router.GET("/account/update", update)
    router.GET("/account/delete", delete)
    router.GET("/account/list", list)
    router.GET("/account/retreive/:id", single)

    mRouter := &myrouter {
        router: router,
    }

    mRouter.middlewares = []myMiddleware {
        &mw1{},
        &mw2{},
    }

    log.Fatal(http.ListenAndServe(":8081", mRouter))
}

func getDBConnection() *sql.DB {
    connStr := "user=postgres password=Password1! dbname=chitchat_users port=5433 sslmode=disable"
    db, err := sql.Open("postgres", connStr)
    failOnError(err, "Failed to connect to DB")
    return db
}

func failOnError(err error, msg string) {
    if err != nil {
        log.Fatalf("%s: %s", msg, err)
    }
}

func (h *myrouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    log.Printf("Serving %s", r.URL.Path)

    for _, mw := range h.middlewares {
        err := mw.intercept(w, r)
        failOnError(err, "Failed on middleware")
    }

    header := w.Header()
    header.Add("Content-Type", "application/json")
    h.router.ServeHTTP(w, r)
}

func (m mw1) intercept(w http.ResponseWriter, r *http.Request) error {
    log.Printf("Intercepted on m1")
    return nil
}

func (m mw2) intercept(w http.ResponseWriter, r *http.Request) error {
    log.Printf("Intercepted on m2")
    return nil
}