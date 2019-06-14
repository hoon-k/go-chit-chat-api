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

// BindQueueToExchange binds queue to exchange
func BindQueueToExchange(ch *amqp.Channel, exchangeName string, queueName string, routeKey string) {
    err := ch.QueueBind(
        queueName,      // queue name
        routeKey,       // routing key
        exchangeName,   // exchange
        false,
        nil,
    )

    failOnError(err, "Failed to bind a queue")
}

// ConsumeMessages consumes messages
func ConsumeMessages(ch *amqp.Channel, queueName string) <-chan amqp.Delivery {
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

func receiveMessage(queueName string, exchangeType ExchangeType, exchangeName string, routeKeys []string, isDurableQueue bool, isDurableExchange bool) (<-chan amqp.Delivery,*amqp.Connection, *amqp.Channel) {
    conn := Connect()

    ch := CreateChannel(conn)

    if exchangeType != DefaultExchange {
        DeclareExchange(ch, exchangeName, exchangeType, isDurableExchange)
    }

    q := DeclareQueue(ch, queueName, isDurableQueue)

    if exchangeType != DefaultExchange {
        if len(routeKeys) == 1 {
            BindQueueToExchange(ch, exchangeName, q.Name, routeKeys[0])
        } else {
            for _, routeKey := range routeKeys {
                BindQueueToExchange(ch, exchangeName, q.Name, routeKey)
            }
        }
    }

    msgs := ConsumeMessages(ch, q.Name)

    return msgs, conn, ch
}