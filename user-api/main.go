package main

import (
    "database/sql"
    "fmt"
    "log"
    "net/http"

    "github.com/julienschmidt/httprouter"
    _ "github.com/lib/pq"
    "github.com/streadway/amqp"
)

var conn *amqp.Connection

// Index function
func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    db := getDBConnection()
    defer db.Close()

    rows, err := db.Query("SELECT first_name, last_name FROM users")
    defer rows.Close()

    var firstName string
    var lastName string

    for rows.Next() {
        err := rows.Scan(&firstName, &lastName)
        if err != nil {
            log.Fatal(err)
        }

        log.Println(firstName, lastName)
        fmt.Fprintf(w, "Name is %s %s\n", firstName, lastName)
    }

    conn := getMQConnection()
    defer conn.Close()

    ch, err := conn.Channel()
    failOnError(err, "Failed to open a channel")
    defer ch.Close()

    q, err := ch.QueueDeclare(
        "hello", // name
        false,   // durable
        false,   // delete when unused
        false,   // exclusive
        false,   // no-wait
        nil,     // arguments
    )

    failOnError(err, "Failed to declare a queue")

    body := "Hello World!"

    err = ch.Publish(
        "",     // exchange
        q.Name, // routing key
        false,  // mandatory
        false,  // immediate
        amqp.Publishing {
            ContentType: "text/plain",
            Body:        []byte(body),
    })
    failOnError(err, "Failed to publish a message")
}

func main() {
    router := httprouter.New()
    router.GET("/", Index)
    router.POST("/account/create", Create)
    // router.POST("/account/list", List)
    // router.POST("/account/:id", Single)

    log.Fatal(http.ListenAndServe(":8081", router))
}

func getDBConnection() *sql.DB {
    connStr := "user=postgres password=Password1! dbname=chitchat_users port=5433 sslmode=disable"
    db, err := sql.Open("postgres", connStr)
    failOnError(err, "Failed to connect to DB")
    return db
}

func getMQConnection() *amqp.Connection {
    conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
    failOnError(err, "Failed to connect to RabbitMQ")
    return conn
}

func failOnError(err error, msg string) {
    if err != nil {
        log.Fatalf("%s: %s", msg, err)
    }
}