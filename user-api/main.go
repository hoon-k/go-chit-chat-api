package main

import (
    "log"
    "net/http"

    "go-chit-chat-api/middlewares"
    "go-chit-chat-api/middlewares/logger"
    "go-chit-chat-api/middlewares/validators"

    _ "github.com/lib/pq"
)

func main() {
    router := initializeRouter()
    mrRouter := middlewares.CreateManagedRouter(router)

    mrRouter.Add(&logger.Logger{})
    mrRouter.Add(&request.SchemaValidator{})

    log.Fatal(http.ListenAndServe(":8081", mrRouter))
}

func failOnError(err error, msg string) {
    if err != nil {
        log.Fatalf("%s: %s", msg, err)
    }
}