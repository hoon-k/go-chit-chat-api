package main

import (
    "database/sql"
    "fmt"
    "log"
    "net/http"

    "github.com/julienschmidt/httprouter"
    _ "github.com/lib/pq"

    "go-chit-chat-api/middlewares"
    "go-chit-chat-api/middlewares/logger"
    "go-chit-chat-api/middlewares/validators"
)

type myrouter struct {
    middlewares []middlewares.Middleware
    router *httprouter.Router
}

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
    mRouter := &myrouter {
        router: initializeRouter(),
    }

    mRouter.middlewares = []middlewares.Middleware {
        &logger.Logger{},
        &request.SchemaValidator{},
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
        err := mw.Intercept(w, r)
        failOnError(err, "Failed on middleware")
    }

    header := w.Header()
    header.Add("Content-Type", "application/json")
    h.router.ServeHTTP(w, r)
}