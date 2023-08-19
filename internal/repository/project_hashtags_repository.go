package repository

import (
	"database/sql"
	"fold/internal/database"
	"log"
)

func CreateProjectHashtags(tx *sql.Tx, hashtagId int, projectId int) error {
	_, err := tx.Exec("INSERT INTO project_hashtags(hashtag_id, project_id) VALUES($1, $2)", hashtagId, projectId)
	return err
}

func ExistsProjectHastags(hashtagId int, projectId int) bool {
	var exists bool
	err := database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM project_hashtags WHERE hashtag_id = $1 AND project_id = $2)", hashtagId, projectId).Scan(&exists)
	if err != nil {
		log.Fatal(err)
		return false
	}
	return exists
}

func DeleteProjectHashtags(tx *sql.Tx, projectId int) error {
	_, err := tx.Exec("DELETE FROM project_hashtags WHERE project_id = $1", projectId)
	return err
}
