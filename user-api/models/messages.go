package models

// CreateUserMessage model
type CreateUserMessage struct {
    UserName string
    FirstName string
    LastName string
    Role string
}

// DeleteUserMessage model
type DeleteUserMessage struct {
    UserName string
    FirstName string
    LastName string
}