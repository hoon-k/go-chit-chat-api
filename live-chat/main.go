package main

import (
    "log"
    "net/http"

    "go-chit-chat-api/events"
    "go-chit-chat-api/middlewares"
    "go-chit-chat-api/middlewares/logger"
    "go-chit-chat-api/middlewares/validators"

    "github.com/julienschmidt/httprouter"
)

type messageReceivedHandler struct {
    writer http.ResponseWriter
    req *http.Request
    params httprouter.Params
}

func main() {
    router := initializeRouter()
    mrRouter := middlewares.CreateManagedRouter(router)

    mrRouter.Add(&logger.Logger{})
    mrRouter.Add(&request.SchemaValidator{})

    log.Fatal(http.ListenAndServe(":8085", mrRouter))
}

func initializeRouter() *httprouter.Router {
    router := httprouter.New()
    router.GET("/live-chat/push", pushMessage)
    router.GET("/live-chat/pull", pullMessage)

    return router
}

func pushMessage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    manager := event.ManagerInstance()
    manager.Publish(event.ChatMessagePublished, "some message")
}

func pullMessage(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    manager := event.ManagerInstance()
    manager.AddSubscription(event.ChatMessagePublished, &messageReceivedHandler{
        writer: w,
        req: r,
        params: p,
    })
    manager.WaitForMessagesForDispatching()
}

func (h *messageReceivedHandler) Handle(msg interface{}, e event.Event) {
    log.Printf("Handling %s event with message %s %v", string(e), msg.(string), h.writer)

    // res, _ := json.Marshal(msg)
    s := msg.(string)
    h.writer.Write([]byte(s))
}