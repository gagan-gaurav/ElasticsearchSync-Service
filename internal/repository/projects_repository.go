package repository

import (
	"database/sql"
	"fmt"
	"fold/internal/database"
	"fold/internal/models"
	"fold/internal/services"
	"log"
	"time"
)

func CreateProject(tx *sql.Tx, project *models.Project) (int, error) {
	var projectId int
	err := tx.QueryRow(
		"INSERT INTO projects (name, slug, description, created_at) VALUES ($1, $2, $3, $4) RETURNING id",
		project.Name, project.Slug, project.Description, time.Now()).Scan(&projectId)

	return projectId, err
}

func GetProjectById(projectId int, project *models.Project) error {
	err := database.DB.QueryRow("SELECT * FROM projects WHERE id = $1", projectId).Scan(&project.ID, &project.Name, &project.Slug, &project.Description, &project.CreatedAt)
	return err
}

func ProjectExists(projectId int) bool {
	var exists bool
	err := database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM projects WHERE id = $1)", projectId).Scan(&exists)
	if err != nil {
		log.Fatal(err)
		return false
	}
	return exists
}

func GetProjectUsers(tx *sql.Tx, projectId int, doc *models.DenormalizedProject) error {
	rows, err := tx.Query("SELECT u.id, u.name, u.created_at FROM users u JOIN user_projects p ON u.id = p.user_id WHERE p.project_id = $1", projectId)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Name, &user.CreatedAt)
		if err != nil {
			return err
		}
		doc.Users = append(doc.Users, user)
	}

	return err
}

func GetProjectHashtags(tx *sql.Tx, projectId int, doc *models.DenormalizedProject) error {
	rows, err := tx.Query("SELECT h.id, h.name, h.created_at FROM hashtags h JOIN project_hashtags p ON h.id = p.hashtag_id WHERE p.project_id = $1", projectId)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var hashtag models.Hashtag
		err := rows.Scan(&hashtag.ID, &hashtag.Name, &hashtag.CreatedAt)
		if err != nil {
			return err
		}
		doc.Hashtags = append(doc.Hashtags, hashtag)
	}

	return err
}

func ProjectCreationAndSyncTransaction(project *models.Project) error {
	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}

	// Insert project into the database
	projectId, err := CreateProject(tx, project)
	if err != nil {
		fmt.Println(err)
		tx.Rollback()
		return err
	}

	// Create entries in user_projects.
	users := project.UserIds

	for _, userId := range users {
		err = CreateUserProjects(tx, projectId, userId)
		if err != nil {
			fmt.Println(err)
			tx.Rollback()
			return err
		}
	}

	// Create entries in project_hashtags.

	hashtags := project.HashtagIds

	for _, hashtagId := range hashtags {
		err = CreateProjectHashtags(tx, hashtagId, projectId)
		if err != nil {
			fmt.Println(err)
			tx.Rollback()
			return err
		}
	}

	// Perform Denormalization of project and send it to sqs queue.

	doc := createDoc(projectId, project)

	err = GetProjectUsers(tx, projectId, &doc)
	if err != nil {
		fmt.Println(err)
		tx.Rollback()
		return err
	}

	err = GetProjectHashtags(tx, projectId, &doc)
	if err != nil {
		fmt.Println(err)
		tx.Rollback()
		return err
	}

	err = services.SQS(&doc)
	if err != nil {
		fmt.Println("SQS push failed")
		tx.Rollback()
		return err
	}

	// Commit the transaction
	return tx.Commit()
}

func createDoc(projectId int, project *models.Project) models.DenormalizedProject {
	var doc models.DenormalizedProject
	doc.ID = projectId
	doc.Slug = project.Slug
	doc.Name = project.Name
	doc.CreatedAt = project.CreatedAt
	doc.Description = project.Description
	return doc
}
