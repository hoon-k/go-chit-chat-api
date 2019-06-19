package services

import (
    "go-chit-chat-api/user-api/models"
)

// CreateUser creates a user and returns a message to broadcast on success
func CreateUser(req *models.CreateUserRequest) (*models.CreateUserMessage, error) {
    db := getDBConnection()
    defer db.Close()

    rows, err := db.Query(`SELECT * FROM create_user($1, $2, $3, $4)`, req.UserName, req.Password, req.FirstName, req.LastName)

    failOnError(err, "Unable to create new user")

    if err != nil {
        return nil, err
    }
    
    msg := &models.CreateUserMessage{}
    rows.Next()
    rows.Scan(&msg.FirstName, &msg.LastName, &msg.UserName, &msg.Role)

    return msg, nil
}

// DeleteUser deletes user
func DeleteUser(id string) (*models.DeleteUserMessage, error) {
    db := getDBConnection()
    defer db.Close()

    rows, err := db.Query(`SELECT * FROM delete_user($1)`, id)

    failOnError(err, "Unable to delete a user")

    msg := &models.DeleteUserMessage{}
    rows.Next()
    rows.Scan(&msg.UserName, &msg.FirstName, &msg.LastName)

    return msg, err
}

// GetAllUsers lists all users
func GetAllUsers() *models.UserResults {
    db := getDBConnection()
    defer db.Close()

    rows, _ := db.Query("SELECT first_name, last_name FROM users")
    defer rows.Close()

    userResults := &models.UserResults{}

    count := 0
    for rows.Next() {
        user := &models.User{}
        err := rows.Scan(&user.FirstName, &user.LastName)
        
        failOnError(err, "Unable to get a list of users")

        userResults.Users = append(userResults.Users, *user)

        count = count + 1
    }

    userResults.TotalNumber = count

    return userResults
}