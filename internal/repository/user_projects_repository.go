package repository

import (
	"database/sql"
	"fold/internal/database"
	"log"
)

func CreateUserProjects(tx *sql.Tx, projectId int, userId int) error {
	_, err := tx.Exec("INSERT INTO user_projects(project_id, user_id) VALUES($1, $2)", projectId, userId)
	return err
}

func ExistsUserProjects(projectId int, userId int) bool {
	var exists bool
	err := database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM user_projects WHERE project_id = $1 AND user_id = $2)", projectId, userId).Scan(&exists)
	if err != nil {
		log.Fatal(err)
		return false
	}
	return exists
}
