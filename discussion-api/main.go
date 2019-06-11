package main

import (
    // "database/sql"
    "fmt"
    "log"

    "github.com/streadway/amqp"
    // _ "github.com/lib/pq"
)

func main() {
    conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
    failOnError(err, "Failed to connect to RabbitMQ")
    defer conn.Close()

    ch, err := conn.Channel()
    failOnError(err, "Failed to open a channel")
    defer ch.Close()    

    fmt.Printf("hello, world discussion api\n")
    
    q, err := ch.QueueDeclare(
        "task_queue", // name
        true,   // durable
        false,   // delete when usused
        false,   // exclusive
        false,   // no-wait
        nil,     // arguments
    )

    failOnError(err, "Failed to declare a queue")

    err = ch.Qos(
        1,     // prefetch count
        0,     // prefetch size
        false, // global
    )

    msgs, err := ch.Consume(
        q.Name, // queue
        "",     // consumer
        false,   // auto-ack
        false,  // exclusive
        false,  // no-local
        false,  // no-wait
        nil,    // args
    )

    failOnError(err, "Failed to register a consumer")

    forever := make(chan bool)

    go receiveMessage(msgs)

    log.Printf(" [*] Waiting for messages. To exit press CTRL+C")

    <-forever
}

func createAuthor() {
    
}

func receiveMessage(msgs <-chan amqp.Delivery) {
    for d := range msgs {
        log.Printf("Received a message: %s", d.Body)
        d.Ack(false)
    }
}

func failOnError(err error, msg string) {
    if err != nil {
        log.Fatalf("%s: %s", msg, err)
    }
  }
