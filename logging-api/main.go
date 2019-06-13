package main

import (
    // "database/sql"
    "fmt"
    "log"
    "go-chit-chat-api/mq"
    "github.com/streadway/amqp"
    // _ "github.com/lib/pq"
)

func main() {
    receiveMessages()
}

func receiveMessages() {
    msgs, conn, ch := mq.ReceiveMessageFromMultipleRoutes("chitchat", []string {"userCreated", "userUpdated"})
    defer conn.Close()
    defer ch.Close()

    fmt.Printf("hello, world reports api\n")

    forever := make(chan bool)

    go processMessages(msgs)

    log.Printf(" [*] Waiting for messages. To exit press CTRL+C")

    <-forever
}

func processMessages(msgs <-chan amqp.Delivery) {
    for d := range msgs {
        log.Printf("Received a message in logging api: %s with route key of %s", string(d.Body), string(d.RoutingKey))
        d.Ack(false)
    }
}