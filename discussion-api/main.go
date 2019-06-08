package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func main() {
    connStr := "user=postgres password=Password1! dbname=pqgotest port=5433"
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal(err)
    }

    rows, _ := db.Query("SELECT first_name, last_name FROM users")
    defer rows.Close()

    firstName := ""
    lastName := ""

    for rows.Next() {
        err := rows.Scan(&firstName, &lastName)
        if err != nil {
            log.Fatal(err)
        }

        log.Println(firstName, lastName)
        fmt.Printf("Name is %v %v", firstName, lastName)
    }

	fmt.Printf("hello, world discussion api\n")
}
