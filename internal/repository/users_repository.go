package repository

import (
	"fold/internal/database"
	"fold/internal/models"
	"log"
	"time"
)

func CreateUser(user *models.User) error {
	_, err := database.DB.Exec("INSERT INTO users (name, created_at) VALUES ($1, $2)", user.Name, time.Now())
	return err
}

func GetUserById(userId int, user *models.User) error {
	err := database.DB.QueryRow("SELECT * FROM users WHERE id = $1", userId).Scan(&user.ID, &user.Name, &user.CreatedAt)
	return err
}

func UserExists(userId int) bool {
	var exists bool
	err := database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)", userId).Scan(&exists)
	if err != nil {
		log.Fatal(err)
		return false
	}
	return exists
}
