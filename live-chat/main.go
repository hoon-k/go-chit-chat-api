package main

import (
    "encoding/json"
    "log"
    "net/http"
    "io"

    "go-chit-chat-api/events"
    "go-chit-chat-api/middlewares"
    "go-chit-chat-api/middlewares/logger"
    "go-chit-chat-api/middlewares/validators"

    "github.com/julienschmidt/httprouter"
    "github.com/google/uuid"
)

type messageReceivedHandler struct {}

type message struct {
    UserID int `json:"userID"`
    ChannelID string `json:"channelID"`
    SentTime string `json:"sentTime"`
    Message string `json:"message"`
}

type channel struct {
    ChannelID string `json:"channelID"`
}

var chatChannels = make(map[string]chan string)
// var messages = make(chan string)

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
    router.GET("/live-chat/create", createChannel)

    return router
}

func pushMessage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    decoder := json.NewDecoder(r.Body)
    var msg message
    err := decoder.Decode(&msg)
    log.Printf("Message %v", msg)
    if err != nil {
        panic(err)
    }

    manager := event.ManagerInstance()

    log.Printf("Message %v", msg)

    manager.Publish(event.ChatMessagePublished, &msg)
}

func pollMessage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    chID := r.Header.Get("X-Channel-ID")
    log.Printf("Serving channel %s", chID)

    messages := chatChannels[chID]
    io.WriteString(w, <-messages)
}

func createChannel(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    uuid, _ := uuid.NewUUID()
    messages := make(chan string)
    chatChannels[uuid.String()] = messages
    
    res, _ := json.Marshal(&channel{
        ChannelID: uuid.String(),
    })
    log.Printf("Channdel ID %s ", res)
    w.Write(res)
}

func (h *messageReceivedHandler) Handle(msg []byte, e event.Event) {
    log.Printf("Handling %s event with message %s", string(e), msg)
    s := string(msg)

    var receivedMsg message
    json.Unmarshal(msg, &receivedMsg)
    log.Printf("Channdel ID %s ", receivedMsg.ChannelID)
    messages := chatChannels[receivedMsg.ChannelID]
    messages <- s
}