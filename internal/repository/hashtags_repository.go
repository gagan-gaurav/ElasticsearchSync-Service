package repository

import (
	"database/sql"
	"fmt"
	"fold/internal/database"
	"fold/internal/models"
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

func UpdateHashtag(tx *sql.Tx, hashtag *models.Hashtag) error {
	_, err := database.DB.Exec("UPDATE hashtags SET name = $1 WHERE id = $2", hashtag.Name, hashtag.ID)
	return err
}

func DeleteHashtag(tx *sql.Tx, hashtagId int) error {
	_, err := tx.Exec("DELETE FROM hashtags WHERE id = $1", hashtagId)
	return err
}

func HashtagExists(hashtagId int) bool {
	var exists bool
	err := database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM hashtags WHERE id = $1)", hashtagId).Scan(&exists)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return exists
}

func GetHashtagProjectIds(tx *sql.Tx, hashtagId int, projectIds *[]int) error {
	rows, err := tx.Query("SELECT project_id FROM project_hashtags WHERE hashtag_id = $1", hashtagId)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		if err != nil {
			return err
		}
		*projectIds = append(*projectIds, id)
	}

	return err
}

func DeleteHashtagProjectIds(tx *sql.Tx, hashtagId int) error {
	_, err := tx.Exec("DELETE FROM project_hashtags WHERE hashtag_id = $1", hashtagId)
	return err
}

func UpdateHashtagTransaction(hashtag *models.Hashtag) error {

	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}

	//Update user in database
	err = UpdateHashtag(tx, hashtag)
	if err != nil {
		tx.Rollback()
		return err
	}

	//Get list of projectIds that need to be changed.
	var projectIds []int
	err = GetHashtagProjectIds(tx, hashtag.ID, &projectIds)
	if err != nil {
		tx.Rollback()
		return err
	}

	//Sync Elastic Search for every project edited.
	for _, projectId := range projectIds {
		err = SyncElasticsearch(tx, projectId, "POST")
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func DeleteHashtagTransaction(hashtagId int) error {

	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}

	//Get list of projectIds that need to be changed.
	var projectIds []int
	err = GetHashtagProjectIds(tx, hashtagId, &projectIds)
	if err != nil {
		tx.Rollback()
		return err
	}

	//Delete rows from project_hashtags
	err = DeleteHashtagProjectIds(tx, hashtagId)
	if err != nil {
		tx.Rollback()
		return err
	}

	//Sync Elastic Search for every project deleted.
	for _, projectId := range projectIds {
		err = SyncElasticsearch(tx, projectId, "POST")
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	//Update user in database
	err = DeleteUser(tx, hashtagId)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
