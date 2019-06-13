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
    // msgs, conn, ch := mq.GetMessagesFromDefaultExchange("task_queue")
    msgs, conn, ch := mq.GetMessages("", mq.FanoutExchange, "chitchat")
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
        log.Printf("Received a message: %s", string(d.Body))
        d.Ack(false)
    }
}
