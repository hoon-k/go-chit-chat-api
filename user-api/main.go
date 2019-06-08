package main

import (
    "database/sql"
    "fmt"
    "log"

    _ "github.com/lib/pq"
)

func main() {
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
        fmt.Printf("Name is %v %v\n", firstName, lastName)
    }

	fmt.Printf("hello, world user api\n")
}
