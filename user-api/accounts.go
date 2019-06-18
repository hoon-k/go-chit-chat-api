package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"

    "github.com/julienschmidt/httprouter"
    _ "github.com/lib/pq"

    "go-chit-chat-api/events"
    "go-chit-chat-api/user-api/models"
    "go-chit-chat-api/user-api/services"
)

func list(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    services.GetAllUsers()
}

func create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    decoder := json.NewDecoder(r.Body)
    var req models.CreateUserRequest
    err := decoder.Decode(&req)
    if err != nil {
        panic(err)
    }

    log.Printf("Post is %s %s\n", req.FirstName, req.LastName)

    msg, err := services.CreateUser(&req)

    failOnError(err, "Unable to create new user")

    manager := event.ManagerInstance()
    manager.Publish(event.UserCreated, &msg)
}

func update(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    msg := &models.CreateUserMessage{
        UserName: "Updated!",
    }

    manager := event.ManagerInstance()
    manager.Publish(event.UserUpdated, msg)

    res, _ := json.Marshal(msg)
    w.Write(res)
}

func delete(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    uuid := p.ByName("id")
    msg, err := services.DeleteUser(uuid)

    failOnError(err, fmt.Sprintf("Unable to delete user %s", uuid))

    manager := event.ManagerInstance()
    manager.Publish(event.UserDeleted, msg)
}

func single(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    msg := &models.CreateUserMessage{
        UserName: "Updated!",
    }

    fmt.Fprintf(w, "Updated and sending msg: %s", p.ByName("id"))

    res, _ := json.Marshal(msg)
    w.Write(res)
}