package main

import (
    "encoding/json"
    // "io/ioutil"
    // "database/sql"
    "fmt"
    "log"
    "net/http"

    "github.com/julienschmidt/httprouter"
    _ "github.com/lib/pq"
    // "github.com/streadway/amqp"
    "go-chit-chat-api/mq"
)

type aData struct {
    Point string `json:"point"`
}

type createUserRequest struct {
    UserName string `json:"userName"`
    Password string `json:"password"`
    FirstName string `json:"firstName"`
    LastName string `json:"lastName"`
    Data aData `json:"aData"`

}

type createUserMessage struct {
    UserName string
    FirstName string
    LastName string
    Role string
}

func list(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    db := getDBConnection()
    defer db.Close()

    rows, _ := db.Query("SELECT first_name, last_name FROM users")
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

// Create user
func create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    decoder := json.NewDecoder(r.Body)
    var req createUserRequest
    err := decoder.Decode(&req)
    if err != nil {
        panic(err)
    }

    log.Printf("Post is %s %s\n", req.FirstName, req.LastName)
    db := getDBConnection()
    defer db.Close()

    rows, err := db.Query(`SELECT * FROM create_user($1, $2, $3, $4)`, req.UserName, req.Password, req.FirstName, req.LastName)

    failOnError(err, "Unable to create new user")
    
    msg := createUserMessage{}
    rows.Next()
    rows.Scan(&msg.FirstName, &msg.LastName, &msg.UserName, &msg.Role)

    mq.SendMessagesToDefaultExchange(&msg, "task_queue")
}