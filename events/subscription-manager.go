package event

import (
    "log"

    "github.com/streadway/amqp"

    "go-chit-chat-api/mq"
)

const exchange = "chitchat"

// IHandler interface
type IHandler interface {
    Handle(interface{}, Event)
}

var instance *SubscriptionManager

// SubscriptionManager to handle events and handlers
type SubscriptionManager struct {
    queueName string
    conn *amqp.Connection
    ch *amqp.Channel
    handlers map[string][]IHandler
}

// Publish publishes event to MQ
func (m *SubscriptionManager) Publish(event Event, msg interface{}) {
    routeKey := string(event)
    mq.SendMessageToRoute(&msg, exchange, routeKey)
}

// AddSubscription adds an event handler
func (m *SubscriptionManager) AddSubscription(event Event, handler IHandler) {
    routeKey := string(event)
    m.handlers[routeKey] = append(m.handlers[routeKey], handler)

    if m.ch != nil {
        mq.BindQueueToExchange(m.ch, exchange, m.queueName, routeKey)
    }
}

// RemoveSubscription removes an event handler
func (m *SubscriptionManager) RemoveSubscription(event Event, handler IHandler) {
    
}

// GetHandlersForEvent gets all handlers given an event
func (m *SubscriptionManager) GetHandlersForEvent(event Event) []IHandler {
    routeKey := string(event)
    return m.handlers[routeKey]
}

// Clear removes all event handlers and cleans things up
func (m *SubscriptionManager) Clear() {
    m.handlers = make(map[string][]IHandler)
}

// WaitForMessagesForDispatching waits for messages
func (m *SubscriptionManager) WaitForMessagesForDispatching() {
    if m.conn == nil || m.ch == nil {
        panic("MQ not ready!")
    }

    defer m.conn.Close()
    defer m.ch.Close()

    msgs := mq.ConsumeMessages(m.ch, m.queueName)

    forever := make(chan bool)

    go m.dispatchMessage(msgs)

    log.Printf(" [*] Waiting for messages. To exit press CTRL+C")

    <-forever
}

// ManagerInstance returns an instance of SubscriptionManager
func ManagerInstance() *SubscriptionManager {
    if instance == nil {
        instance = &SubscriptionManager {
            handlers: make(map[string][]IHandler),
        }
        instance.createQueue()
    }

    return instance;
}

func (m *SubscriptionManager) createQueue() {
    m.conn = mq.Connect()
    m.ch = mq.CreateChannel(m.conn)
    mq.DeclareExchange(m.ch, exchange, mq.DirectExchange, true)
    q := mq.DeclareQueue(m.ch, "", false)
    m.queueName = q.Name
}

func (m *SubscriptionManager) dispatchMessage(msgs <-chan amqp.Delivery) {
    for d := range msgs {
        // log.Printf("Received a message in report api: %s with route key of %s", string(d.Body), string(d.RoutingKey))
        routeKey := string(d.RoutingKey)
        for _, handler := range m.handlers[routeKey] {
            handler.Handle(string(d.Body), Event(routeKey))
        }
        d.Ack(false)
    }
}