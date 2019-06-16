package main

import "github.com/julienschmidt/httprouter"

func initializeRouter() *httprouter.Router {
    router := httprouter.New()
    router.GET("/", Index)
    router.POST("/accounts/create", create)
    router.GET("/accounts/update", update)
    router.GET("/accounts/delete", delete)
    router.GET("/accounts/list", list)
    router.GET("/accounts/retreive/:id", single)

    return router
}