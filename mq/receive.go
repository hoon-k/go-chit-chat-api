package mq

import (
    "github.com/streadway/amqp"
)

// ReceiveMessageFromQueue gets message from a specific queue.
func ReceiveMessageFromQueue(queueName string) (<-chan amqp.Delivery,*amqp.Connection, *amqp.Channel) {
    return receiveMessage(queueName, DefaultExchange, "", []string{""}, false, false)
}

// ReceiveMessageFromExchange gets message from a specific exchange
func ReceiveMessageFromExchange(exchangeName string) (<-chan amqp.Delivery,*amqp.Connection, *amqp.Channel) {
    return receiveMessage("", FanoutExchange, exchangeName, []string{""}, false, true)
}

// ReceiveMessageFromRoute gets message from a specifc route
func ReceiveMessageFromRoute(exchangeName string, routeKey string) (<-chan amqp.Delivery,*amqp.Connection, *amqp.Channel) {
    return receiveMessage("", DirectExchange, exchangeName, []string{routeKey}, false, true)
}

// ReceiveMessageFromMultipleRoutes gets message from multiple routes
func ReceiveMessageFromMultipleRoutes(exchangeName string, routeKeys []string) (<-chan amqp.Delivery,*amqp.Connection, *amqp.Channel) {
    return receiveMessage("", DirectExchange, exchangeName, routeKeys, false, true)
}

func receiveMessage(queueName string, exchangeType ExchangeType, exchangeName string, routeKeys []string, isDurableQueue bool, isDurableExchange bool) (<-chan amqp.Delivery,*amqp.Connection, *amqp.Channel) {
    conn := connect()

    ch := createChannel(conn)

    if exchangeType != DefaultExchange {
        declareExchange(ch, exchangeName, exchangeType, isDurableExchange)
    }

    q := declareQueue(ch, queueName, isDurableQueue)

    if exchangeType != DefaultExchange {
        if len(routeKeys) == 1 {
            bindQueueToExchange(ch, exchangeName, q.Name, routeKeys[0])
        } else {
            for _, routeKey := range routeKeys {
                bindQueueToExchange(ch, exchangeName, q.Name, routeKey)
            }
        }
    }

    msgs := consumeMessages(ch, q.Name)

    return msgs, conn, ch
}

func bindQueueToExchange(ch *amqp.Channel, exchangeName string, queueName string, routeKey string) {
    err := ch.QueueBind(
        queueName,      // queue name
        routeKey,       // routing key
        exchangeName,   // exchange
        false,
        nil,
    )

    failOnError(err, "Failed to bind a queue")
}

func consumeMessages(ch *amqp.Channel, queueName string) <-chan amqp.Delivery {
    msgs, err := ch.Consume(
        queueName,  // queue
        "",         // consumer
        false,      // auto-ack
        false,      // exclusive
        false,      // no-local
        false,      // no-wait
        nil,        // args
    )

    failOnError(err, "Failed to register a consumer")

    return msgs
}