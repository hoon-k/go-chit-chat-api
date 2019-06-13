package main

import (
    // "database/sql"
    "fmt"
    "log"
    "go-chit-chat-api/mq"
    "net/http"

    "github.com/julienschmidt/httprouter"
    "github.com/streadway/amqp"
    // _ "github.com/lib/pq"
)

func main() {
    router := httprouter.New()
    router.GET("/posts/create", create)
    go http.ListenAndServe(":8082", router)
    receiveMessages()
}

func receiveMessages() {
    msgs, conn, ch := mq.ReceiveMessageFromMultipleRoutes("chitchat", []string {"userCreated", "userDeleted"})
    defer conn.Close()
    defer ch.Close()

    fmt.Printf("hello, world discussion api\n")

    forever := make(chan bool)

    go processMessages(msgs)

    log.Printf(" [*] Waiting for messages. To exit press CTRL+C")

    <-forever
}

func createAuthor() {
    
}

func processMessages(msgs <-chan amqp.Delivery) {
    for d := range msgs {
        log.Printf("Received a message in discussion api: %s with route key of %s", string(d.Body), string(d.RoutingKey))
        d.Ack(false)
    }
}
