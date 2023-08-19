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

func GetAllProjects() ([]models.Project, error) {
	var projects []models.Project

	rows, err := database.DB.Query("SELECT id, name, slug, description, created_at FROM projects")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var project models.Project
		err := rows.Scan(&project.ID, &project.Name, &project.Slug, &project.Description, &project.CreatedAt)
		if err != nil {
			return nil, err
		}
		project.UserIds, err = GetProjectUsersId(project.ID)
		if err != nil {
			return nil, err
		}
		project.HashtagIds, err = GetProjectHashtagsId(project.ID)
		if err != nil {
			return nil, err
		}

		projects = append(projects, project)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return projects, nil
}

func GetProjectByIdForTransaction(tx *sql.Tx, projectId int, project *models.Project) error {
	err := tx.QueryRow("SELECT * FROM projects WHERE id = $1", projectId).Scan(&project.ID, &project.Name, &project.Slug, &project.Description, &project.CreatedAt)
	return err
}

func UpdateProject(tx *sql.Tx, project *models.Project) error {
	_, err := tx.Exec("UPDATE projects SET name = $1, slug = $2, description = $3 WHERE id = $4", project.Name, project.Slug, project.Description, project.ID)
	return err
}

func DeleteProject(tx *sql.Tx, projectId int) error {
	_, err := tx.Exec("DELETE FROM projects WHERE id = $1", projectId)
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

func GetProjectUsersId(projectId int) ([]int, error) {
	var userIds []int
	rows, err := database.DB.Query("SELECT u.id FROM users u JOIN user_projects p ON u.id = p.user_id WHERE p.project_id = $1", projectId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var userId int
		err := rows.Scan(&userId)
		if err != nil {
			return nil, err
		}
		userIds = append(userIds, userId)
	}
	return userIds, nil
}

func GetProjectHashtagsId(projectId int) ([]int, error) {
	var hashtagIds []int

	rows, err := database.DB.Query("SELECT h.id FROM hashtags h JOIN project_hashtags p ON h.id = p.hashtag_id WHERE p.project_id = $1", projectId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var hashtagId int
		err := rows.Scan(&hashtagId)
		if err != nil {
			return nil, err
		}
		hashtagIds = append(hashtagIds, hashtagId)
	}

	return hashtagIds, nil
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
		tx.Rollback()
		return err
	}

	// Create entries in user_projects.
	users := project.UserIds

	for _, userId := range users {
		err = CreateProjectUsers(tx, projectId, userId)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// Create entries in project_hashtags.
	hashtags := project.HashtagIds

	for _, hashtagId := range hashtags {
		err = CreateProjectHashtags(tx, hashtagId, projectId)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// Sync ElasticSearch
	err = SyncElasticsearch(tx, projectId, "POST")
	if err != nil {
		tx.Rollback()
		return err
	}
	// Commit the transaction
	return tx.Commit()
}

func ProjectUpdateAndSyncTransaction(project *models.Project) error {
	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}

	// Update project into the database
	err = UpdateProject(tx, project)
	if err != nil {
		tx.Rollback()
		return err
	}

	//Remove old entries in user_projects.
	err = DeleteProjectUsers(tx, project.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Update entries in user_projects.
	users := project.UserIds

	for _, userId := range users {
		err = CreateProjectUsers(tx, project.ID, userId)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	//Remove old entries in project_hastags.
	err = DeleteProjectHashtags(tx, project.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Update entries in project_hashtags.
	hashtags := project.HashtagIds

	for _, hashtagId := range hashtags {
		err = CreateProjectHashtags(tx, hashtagId, project.ID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// Sync ElasticSearch
	err = SyncElasticsearch(tx, project.ID, "POST")
	if err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction
	return tx.Commit()
}

func ProjectDeleteAndSyncTransaction(projectId int) error {
	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}

	//Remove entries in user_projects.
	err = DeleteProjectUsers(tx, projectId)
	if err != nil {
		tx.Rollback()
		return err
	}

	//Remove entries in project_hastags.
	err = DeleteProjectHashtags(tx, projectId)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Delete document from elasticsearch
	err = SyncElasticsearch(tx, projectId, "DELETE")
	if err != nil {
		tx.Rollback()
		return err
	}

	// Delete project into the database
	err = DeleteProject(tx, projectId)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction
	return tx.Commit()
}

func SyncElasticsearch(tx *sql.Tx, projectId int, method string) error {
	// Perform Denormalization of project and send it to sqs queue.
	var project models.Project
	err := GetProjectByIdForTransaction(tx, projectId, &project)
	if err != nil {
		tx.Rollback()
		return err
	}

	doc := createDoc(projectId, &project)

	err = GetProjectUsers(tx, projectId, &doc)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = GetProjectHashtags(tx, projectId, &doc)
	if err != nil {
		tx.Rollback()
		return err
	}

	//create payload
	var payload models.Payload
	payload.Doc = doc
	payload.Method = method

	err = services.SQS(&payload)
	if err != nil {
		fmt.Println("SQS push failed")
		tx.Rollback()
		return err
	}

	return err
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
