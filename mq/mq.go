package mq

import (
    "encoding/json"
    "log"
    "github.com/streadway/amqp"
)

var conn *amqp.Connection

// SendMessages sends MQ message
func SendMessages(msg interface{}, queueName string) {
    conn := connect()
    defer conn.Close()

    ch := createChannel(conn)
    defer ch.Close()

    q := declareQueue(ch, queueName)

    b, _ := json.Marshal(msg)

    publishMessage(ch, q.Name, b)
}

// GetMessages gets MQ message
func GetMessages(queueName string) (<-chan amqp.Delivery,*amqp.Connection, *amqp.Channel) {
    conn := connect()
    // defer conn.Close()

    ch := createChannel(conn)
    // defer ch.Close()

    msgs := consumeMessages(ch, queueName)

    return msgs, conn, ch
}

func connect() *amqp.Connection {
    conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
    failOnError(err, "Failed to connect to RabbitMQ")
    return conn
}

func createChannel(conn *amqp.Connection) *amqp.Channel {
    ch, err := conn.Channel()

    err = ch.Qos(
        1,     // prefetch count
        0,     // prefetch size
        false, // global
    )

    failOnError(err, "Failed to open a channel")
    return ch
}

func declareQueue(ch *amqp.Channel, queueName string) amqp.Queue {
    q, err := ch.QueueDeclare(
        queueName,      // name
        true,           // durable
        false,          // delete when unused
        false,          // exclusive
        false,          // no-wait
        nil,            // arguments
    )

    failOnError(err, "Failed to declare a queue")

    return q
}

func publishMessage(ch *amqp.Channel, queueName string, message []byte) {
    err := ch.Publish(
        "",     // exchange
        queueName, // routing key
        false,  // mandatory
        false,  // immediate
        amqp.Publishing {
            ContentType: "application/json",
            Body:        message,
    })

    failOnError(err, "Failed to publish a message")
}

func consumeMessages(ch *amqp.Channel, queueName string) <-chan amqp.Delivery {
    msgs, err := ch.Consume(
        queueName, // queue
        "",     // consumer
        false,   // auto-ack
        false,  // exclusive
        false,  // no-local
        false,  // no-wait
        nil,    // args
    )

    failOnError(err, "Failed to register a consumer")

    return msgs
}

func failOnError(err error, msg string) {
    if err != nil {
        log.Fatalf("%s: %s", msg, err)
    }
}

