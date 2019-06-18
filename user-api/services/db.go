package services

import (
    "database/sql"
    "log"
)

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