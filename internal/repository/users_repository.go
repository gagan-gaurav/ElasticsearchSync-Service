package repository

import (
	"database/sql"
	"fmt"
	"fold/internal/database"
	"fold/internal/models"
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

func UpdateUser(tx *sql.Tx, user *models.User) error {
	_, err := database.DB.Exec("UPDATE users SET name = $1 WHERE id = $2", user.Name, user.ID)
	return err
}

func DeleteUser(tx *sql.Tx, userId int) error {
	_, err := tx.Exec("DELETE FROM users WHERE id = $1", userId)
	return err
}

func UserExists(userId int) bool {
	var exists bool
	err := database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)", userId).Scan(&exists)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return exists
}

func GetUserProjectIds(tx *sql.Tx, userId int, projectIds *[]int) error {
	rows, err := tx.Query("SELECT project_id FROM user_projects WHERE user_id = $1", userId)
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

func DeleteUserProjectIds(tx *sql.Tx, userId int) error {
	_, err := tx.Exec("DELETE FROM user_projects WHERE user_id = $1", userId)
	return err
}

func UpdateUserTransaction(user *models.User) error {

	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}

	//Update user in database
	err = UpdateUser(tx, user)
	if err != nil {
		tx.Rollback()
		return err
	}

	//Get list of projectIds that need to be changed.
	var projectIds []int
	err = GetUserProjectIds(tx, user.ID, &projectIds)
	if err != nil {
		tx.Rollback()
		return err
	}

	//Sync Elastic Search for every project edited.
	for _, projectId := range projectIds {
		fmt.Println("seomthing")
		err = SyncElasticsearch(tx, projectId, "POST")
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func DeleteUserTransaction(userId int) error {

	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}

	//Get list of projectIds that need to be changed.
	var projectIds []int
	err = GetUserProjectIds(tx, userId, &projectIds)
	if err != nil {
		tx.Rollback()
		return err
	}

	//Delete rows from user_projects
	err = DeleteUserProjectIds(tx, userId)
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
	err = DeleteUser(tx, userId)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
