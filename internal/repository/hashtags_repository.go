package repository

import (
	"fold/internal/database"
	"fold/internal/models"
	"log"
	"time"
)

func CreateHashtag(hashtag *models.Hashtag) error {
	_, err := database.DB.Exec("INSERT INTO hashtags (name, created_at) VALUES ($1, $2)", hashtag.Name, time.Now())
	return err
}

func GetHashtagById(hashtagID int, hashtag *models.Hashtag) error {
	err := database.DB.QueryRow("SELECT * FROM hashtags WHERE id = $1", hashtagID).Scan(&hashtag.ID, &hashtag.Name, &hashtag.CreatedAt)
	return err
}

func GetHashtagByName(hashtagName string, hashtag *models.Hashtag) error {
	err := database.DB.QueryRow("SELECT * FROM hashtags WHERE name = $1", hashtagName).Scan(&hashtag.ID, &hashtag.Name, &hashtag.CreatedAt)
	return err
}

func HashtagExistsByName(hashtagName string) bool {
	var exists bool
	err := database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM hashtags WHERE name = $1)", hashtagName).Scan(&exists)
	if err != nil {
		log.Fatal(err)
		return false
	}
	return exists
}

func HashtagExistsById(hashtagId int) bool {
	var exists bool
	err := database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM hashtags WHERE id = $1)", hashtagId).Scan(&exists)
	if err != nil {
		log.Fatal(err)
		return false
	}
	return exists
}
