package mq

import (
    // "encoding/json"
    "log"
    "github.com/streadway/amqp"
)

// ExchangeType indicates the type of MQ exhange
type ExchangeType string

const (
    // DirectExchange type
    DirectExchange ExchangeType = "direct"

    // FanoutExchange type
    FanoutExchange ExchangeType = "fanout"

    // TopicExchange type
    TopicExchange ExchangeType = "topic"

    // HeadersExchange type
    HeadersExchange ExchangeType = "headers"

    // DefaultExchange type
    DefaultExchange ExchangeType = ""
)


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

func declareExchange(ch *amqp.Channel, exchangeName string, exchangeType ExchangeType, isDurable bool) {
    err := ch.ExchangeDeclare(
        exchangeName,           // name
        string(exchangeType),   // type
        isDurable,              // durable
        false,                  // auto-deleted
        false,                  // internal
        false,                  // no-wait
        nil,                    // arguments
    )

    failOnError(err, "Failed to declare an exchange")
}

func declareQueue(ch *amqp.Channel, queueName string, isDurable bool) amqp.Queue {
    isExclusive := queueName == ""

    q, err := ch.QueueDeclare(
        queueName,      // name
        isDurable,      // durable
        false,          // delete when unused
        isExclusive,    // exclusive
        false,          // no-wait
        nil,            // arguments
    )

    failOnError(err, "Failed to declare a queue")

    return q
}

func failOnError(err error, msg string) {
    if err != nil {
        log.Fatalf("%s: %s", msg, err)
    }
}