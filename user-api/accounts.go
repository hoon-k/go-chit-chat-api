package main

import (
    "encoding/json"
    // "database/sql"
    "fmt"
    // "log"
    "net/http"

    "github.com/julienschmidt/httprouter"
    _ "github.com/lib/pq"
    // "github.com/streadway/amqp"
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

// Create user
func Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    decoder := json.NewDecoder(r.Body)
    var req createUserRequest
    err := decoder.Decode(&req)
    if err != nil {
        panic(err)
    }

    fmt.Fprintf(w, "Post is %s %s\n", req.FirstName, req.Data.Point)
    // db, err := sql.Open("postgres", connStr)
    // defer db.Close()
    // if err != nil {
    //     log.Fatal(err)
    // }

    // r.

    // rows, err := db.Query("CALL create_member_user(?, ?, ?, ?)", )
    // defer rows.Close()

    // conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
    // failOnError(err, "Failed to connect to RabbitMQ")
    // defer conn.Close()

    // ch, err := conn.Channel()
    // failOnError(err, "Failed to open a channel")
    // defer ch.Close()

    // q, err := ch.QueueDeclare(
    //     "hello", // name
    //     false,   // durable
    //     false,   // delete when unused
    //     false,   // exclusive
    //     false,   // no-wait
    //     nil,     // arguments
    // )

    // failOnError(err, "Failed to declare a queue")

    // body := "Hello World!"

    // err = ch.Publish(
    //     "",     // exchange
    //     q.Name, // routing key
    //     false,  // mandatory
    //     false,  // immediate
    //     amqp.Publishing {
    //         ContentType: "text/plain",
    //         Body:        []byte(body),
    // })
    // failOnError(err, "Failed to publish a message")
}