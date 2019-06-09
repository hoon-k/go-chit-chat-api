package main

import (
    "database/sql"
    "fmt"
    "log"
    "net/http"

    "github.com/julienschmidt/httprouter"
    _ "github.com/lib/pq"
)

// Index function
func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    connStr := "user=postgres password=Password1! dbname=chitchat_users port=5433 sslmode=disable"
    db, err := sql.Open("postgres", connStr)
    defer db.Close()
    if err != nil {
        log.Fatal(err)
    }

    rows, err := db.Query("SELECT first_name, last_name FROM users")
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

    log.Fatal(http.ListenAndServe(":8080", router))
}
