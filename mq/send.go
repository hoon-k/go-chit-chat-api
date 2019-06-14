package mq

import (
    "encoding/json"

    "github.com/streadway/amqp"
)

// SendMessageToQueue sends message to a specific queue
func SendMessageToQueue(msg interface{}, queueName string) {
    sendMessage(msg, DefaultExchange, "", queueName, "", false)
}

// SendMessageToExchange sends message to a specific exchange.
func SendMessageToExchange(msg interface{}, exchangeName string) {
    sendMessage(msg, FanoutExchange, exchangeName, "", "", true)
}

// SendMessageToRoute sends message to a specific route.
func SendMessageToRoute(msg interface{}, exchangeName string, routeKey string) {
    sendMessage(msg, DirectExchange, exchangeName, "", routeKey, true)
}

func sendMessage(msg interface{}, exchangeType ExchangeType, exchangeName string, queueName string, routeKey string, isDurable bool) {
    conn := Connect()
    defer conn.Close()

    ch := CreateChannel(conn)
    defer ch.Close()

    if exchangeType != DefaultExchange {
        DeclareExchange(ch, exchangeName, exchangeType, isDurable)
    } else {
        q := DeclareQueue(ch, queueName, isDurable)
        routeKey = q.Name
    }

    msgStr, _ := json.Marshal(msg)

    publishMessage(ch, exchangeName, routeKey, msgStr, isDurable)
}

func publishMessage(ch *amqp.Channel, exchangeName string, routeKey string, message []byte, isPersistent bool) {
    deliveryMode := amqp.Transient
    if (isPersistent) {
        deliveryMode = amqp.Persistent
    }

    err := ch.Publish(
        exchangeName,   // exchange
        routeKey,       // routing key
        false,          // mandatory
        false,          // immediate
        amqp.Publishing {
            DeliveryMode: deliveryMode,
            ContentType: "text/plain",
            Body:        message,
    })

    failOnError(err, "Failed to publish a message")
}