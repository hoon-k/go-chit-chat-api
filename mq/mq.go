package mq

import (
    "encoding/json"
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

var conn *amqp.Connection

// SendMessagesToDefaultExchange sends message to default exchange
func SendMessagesToDefaultExchange(msg interface{}, queueName string) {
    SendMessages(msg, DefaultExchange, "", queueName)
}

// SendMessages sends message
func SendMessages(msg interface{}, exchangeType ExchangeType, exchangeName string, queueName string) {
    conn := connect()
    defer conn.Close()

    ch := createChannel(conn)
    defer ch.Close()

    if exchangeType != DefaultExchange {
        declareExchange(ch, exchangeName, exchangeType)
    }

    q := declareQueue(ch, queueName)

    msgStr, _ := json.Marshal(msg)

    publishMessage(ch, exchangeName, q.Name, msgStr)
}

// GetMessagesFromDefaultExchange gets message from default exchange
func GetMessagesFromDefaultExchange(queueName string) (<-chan amqp.Delivery,*amqp.Connection, *amqp.Channel) {
    return GetMessages(queueName, DefaultExchange, "")
}

// GetMessages gets message
func GetMessages(queueName string, exchangeType ExchangeType, exchangeName string) (<-chan amqp.Delivery,*amqp.Connection, *amqp.Channel) {
    conn := connect()
    // defer conn.Close()

    ch := createChannel(conn)
    // defer ch.Close()

    if exchangeType != DefaultExchange {
        declareExchange(ch, exchangeName, exchangeType)
    }

    q := declareQueue(ch, queueName)

    if exchangeType != DefaultExchange {
        bindQueueToExchange(ch, exchangeName, q.Name)
    }

    msgs := consumeMessages(ch, q.Name)

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

func declareExchange(ch *amqp.Channel, exchangeName string, exchangeType ExchangeType) {
    err := ch.ExchangeDeclare(
        exchangeName,           // name
        string(exchangeType),   // type
        true,                   // durable
        false,                  // auto-deleted
        false,                  // internal
        false,                  // no-wait
        nil,                    // arguments
    )

    failOnError(err, "Failed to declare an exchange")
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

func bindQueueToExchange(ch *amqp.Channel, exchangeName string, queueName string) {
    err := ch.QueueBind(
        queueName,      // queue name
        "",             // routing key
        exchangeName,   // exchange
        false,
        nil,
    )

    failOnError(err, "Failed to bind a queue")
}

func publishMessage(ch *amqp.Channel, exchangeName string, queueName string, message []byte) {
    err := ch.Publish(
        exchangeName,   // exchange
        queueName,      // routing key
        false,          // mandatory
        false,          // immediate
        amqp.Publishing {
            ContentType: "text/plain",
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

