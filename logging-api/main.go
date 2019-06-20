package main

import (
    "log"

    "go-chit-chat-api/events"
)

type userCreatedHandler struct {}
type userDeletedHandler struct {}

func main() {
    manager := event.ManagerInstance()

    manager.AddSubscription(event.UserCreated, &userCreatedHandler{})
    manager.AddSubscription(event.UserDeleted, &userDeletedHandler{})

    manager.WaitForMessagesForDispatching()
}

func (h *userCreatedHandler) Handle(msg []byte, e event.Event) {
    log.Printf("Handling %s event with message %s", string(e), msg)
}

func (h *userDeletedHandler) Handle(msg []byte, e event.Event) {
    log.Printf("Handling %s event with message %s", string(e), msg)
}