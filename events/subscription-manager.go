package event

import (
    // "go-chit-chat-api/mq"
)

// IHandler interface
type IHandler interface {
	Handle(interface{}, Event)
}

var instance *SubscriptionManager

// SubscriptionManager to handle events and handlers
type SubscriptionManager struct {
    handlers map[Event]IHandler
}

// AddSubscription adds an event handler
func (m *SubscriptionManager) AddSubscription(event Event, handler IHandler) {
    
}

// RemoveSubscription removes an event handler
func (m *SubscriptionManager) RemoveSubscription(event Event, handler IHandler) {
    
}

// GetHandlersForEvent gets all handlers given an event
func (m *SubscriptionManager) GetHandlersForEvent(event Event) {

}

// Clear removes all event handlers and cleans things up
func (m *SubscriptionManager) Clear() {
}

// ManagerInstance returns an instance of SubscriptionManager
func ManagerInstance() *SubscriptionManager {
    if instance == nil {
        instance = &SubscriptionManager {
            handlers: make(map[Event]IHandler),
        }
    }

    return instance;
}