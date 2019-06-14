package main

import (
    "log"
    "net/http"

    "github.com/julienschmidt/httprouter"

    "go-chit-chat-api/events"
)

type userCreatedHandler struct {}
type userDeletedHandler struct {}
type userUpdatedHandler struct {}

func main() {
    router := httprouter.New()
    router.GET("/posts/create", create)
    go http.ListenAndServe(":8082", router)
    
    manager := event.ManagerInstance()

    manager.AddSubscription(event.UserCreated, &userCreatedHandler{})
    manager.AddSubscription(event.UserUpdated, &userUpdatedHandler{})
    manager.AddSubscription(event.UserDeleted, &userDeletedHandler{})

    manager.WaitForMessagesForDispatching()
}

func (h *userCreatedHandler) Handle(msg interface{}, e event.Event) {
    log.Printf("Handling %s event with message %s", string(e), msg)
}

func (h *userDeletedHandler) Handle(msg interface{}, e event.Event) {
    log.Printf("Handling %s event with message %s", string(e), msg)
}

func (h *userUpdatedHandler) Handle(msg interface{}, e event.Event) {
    log.Printf("Handling %s event with message %s", string(e), msg)
}
