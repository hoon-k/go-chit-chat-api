package main

import (
    // "database/sql"
    "fmt"
    "log"
    "go-chit-chat-api/mq"
    "go-chit-chat-api/events"
    "github.com/streadway/amqp"
    // _ "github.com/lib/pq"
)

type myHandler struct {}

func main() {
    manager := event.ManagerInstance()

    manager.AddSubscription(event.UserCreated, &myHandler{})

    receiveMessages();
}

func receiveMessages() {
    msgs, conn, ch := mq.ReceiveMessageFromMultipleRoutes("chitchat", []string {"userCreated", "userDeleted", "postCreated"})
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
        log.Printf("Received a message in report api: %s with route key of %s", string(d.Body), string(d.RoutingKey))
        d.Ack(false)
    }
}

func (h *myHandler) Handle(interface{}, event.Event) {

}
