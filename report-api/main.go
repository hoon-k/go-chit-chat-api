package main

import (
    "log"

    "go-chit-chat-api/events"
)

type userCreatedHandler struct {}
type userDeletedHandler struct {}
type postCreatedHandler struct {}

func main() {
    manager := event.ManagerInstance()

    manager.AddSubscription(event.UserCreated, &userCreatedHandler{})
    manager.AddSubscription(event.UserDeleted, &userDeletedHandler{})
    manager.AddSubscription(event.PostCreated, &postCreatedHandler{})

    manager.WaitForMessagesForDispatching()
}

func (h *userCreatedHandler) Handle(msg interface{}, e event.Event) {
    log.Printf("Handling %s event with message %s", string(e), msg)
}

func (h *userDeletedHandler) Handle(msg interface{}, e event.Event) {
    log.Printf("Handling %s event with message %s", string(e), msg)
}

func (h *postCreatedHandler) Handle(msg interface{}, e event.Event) {
    log.Printf("Handling %s event with message %s", string(e), msg)
}
