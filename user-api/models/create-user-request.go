package models

// CreateUserRequest model
type CreateUserRequest struct {
    UserName string `json:"userName"`
    Password string `json:"password"`
    FirstName string `json:"firstName"`
    LastName string `json:"lastName"`
}

// User model
type User struct {
    FirstName string
    LastName string
    Role string
}

// UserResults model
type UserResults struct {
    TotalNumber int
    Users []User
}