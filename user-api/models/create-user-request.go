package models

// CreateUserRequest model
type CreateUserRequest struct {
    UserName string `json:"userName"`
    Password string `json:"password"`
    FirstName string `json:"firstName"`
    LastName string `json:"lastName"`
}