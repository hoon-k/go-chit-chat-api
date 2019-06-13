package event

// Event type
type Event string

const (
    // UserCreated event
    UserCreated Event = "userCreated"

    // UserUpdated event
    UserUpdated Event = "userUpdated"

    // UserDeleted event
    UserDeleted Event = "userDeleted"

    // PostCreated event
    PostCreated Event = "postCreated"

    // PostDeleted event
    PostDeleted Event = "postDeleted"
)