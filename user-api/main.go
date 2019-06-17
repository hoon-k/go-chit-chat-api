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
    router := initializeRouter()
    mrRouter := middlewares.CreateManagedRouter(router)

    mrRouter.Add(&logger.Logger{})
    mrRouter.Add(&request.SchemaValidator{})

    log.Fatal(http.ListenAndServe(":8081", mrRouter))
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