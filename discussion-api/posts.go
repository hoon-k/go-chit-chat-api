package main

import(
    // "encoding/json"
    // "io/ioutil"
    // "database/sql"
    "fmt"
    // "log"
    "net/http"

    "github.com/julienschmidt/httprouter"
    _ "github.com/lib/pq"
    // "github.com/streadway/amqp"
    "go-chit-chat-api/mq"
)

type post struct {
    Body string
}

func create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    msg := post{
        Body: "Created!",
    }
    mq.SendMessageToRoute(&msg, "chitchat", "postCreated")

    fmt.Fprintf(w, "Created and sending msg: %s", "postCreated")
}