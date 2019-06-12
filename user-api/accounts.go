package main

import (
    "encoding/json"
    "io/ioutil"
    // "database/sql"
    "fmt"
    "log"
    "net/http"

    "github.com/julienschmidt/httprouter"
    _ "github.com/lib/pq"
    "github.com/streadway/amqp"
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

    _, err = db.Query(`CALL create_member_user($1, $2, $3, $4)`, req.UserName, req.Password, req.FirstName, req.LastName)

    failOnError(err, "Unable to create new user")

    conn := getMQConnection()
    failOnError(err, "Failed to connect to RabbitMQ")
    defer conn.Close()

    ch, err := conn.Channel()
    failOnError(err, "Failed to open a channel")
    defer ch.Close()

    q, err := ch.QueueDeclare(
        "task_queue",   // name
        true,           // durable
        false,          // delete when unused
        false,          // exclusive
        false,          // no-wait
        nil,            // arguments
    )

    failOnError(err, "Failed to declare a queue")

    body, _ := ioutil.ReadAll(r.Body)

    err = ch.Publish(
        "",     // exchange
        q.Name, // routing key
        false,  // mandatory
        false,  // immediate
        amqp.Publishing {
            ContentType: "text/plain",
            Body:        []byte(body),
    })
    // failOnError(err, "Failed to publish a message")
}