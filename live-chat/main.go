package main

import (
    "encoding/json"
    "log"
    "net/http"
    "io"
    // "time"

    "go-chit-chat-api/events"
    "go-chit-chat-api/middlewares"
    "go-chit-chat-api/middlewares/logger"
    "go-chit-chat-api/middlewares/validators"

    "github.com/julienschmidt/httprouter"
)

type messageReceivedHandler struct {}

type message struct {
    SentTime string `json:"sentTime"`
    Message string `json:"message"`
}

var messages = make(chan string)

func main() {
    router := initializeRouter()
    mrRouter := middlewares.CreateManagedRouter(router)

    mrRouter.Add(&logger.Logger{})
    mrRouter.Add(&request.SchemaValidator{})

    go http.ListenAndServe(":8085", mrRouter)

    manager := event.ManagerInstance()
    manager.AddSubscription(event.ChatMessagePublished, &messageReceivedHandler{})
    manager.WaitForMessagesForDispatching()
}

func initializeRouter() *httprouter.Router {
    router := httprouter.New()
    router.POST("/live-chat/push", pushMessage)
    router.GET("/live-chat/poll", pollMessage)

    return router
}

func pushMessage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    decoder := json.NewDecoder(r.Body)
    var msg message
    err := decoder.Decode(&msg)
    if err != nil {
        panic(err)
    }

    manager := event.ManagerInstance()

    // msg := &message{
    //     SentTime: time.Now().Local().String(),
    //     Message: "some message",
    // }

    log.Printf("Message %s", msg)

    manager.Publish(event.ChatMessagePublished, msg)
}

func pollMessage(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    io.WriteString(w, <-messages)
}

func (h *messageReceivedHandler) Handle(msg []byte, e event.Event) {
    log.Printf("Handling %s event with message %s", string(e), msg)
    s := string(msg)
    messages <- s
}