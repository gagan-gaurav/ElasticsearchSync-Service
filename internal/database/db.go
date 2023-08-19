package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func MakeDatabaseConnection() {
	dbUrl := os.Getenv("POSTGRES_URL")
	db, connectionErr := sql.Open("postgres", dbUrl)
	if connectionErr != nil {
		fmt.Println("Failed to make connection to database.", connectionErr)
	}

	DB = db

	// Create tables if not exist
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name VARCHAR,
			created_at TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS hashtags (
			id SERIAL PRIMARY KEY,
			name VARCHAR,
			created_at TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS projects (
			id SERIAL PRIMARY KEY,
			name VARCHAR,
			slug VARCHAR,
			description TEXT,
			created_at TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS project_hashtags (
			hashtag_id INT REFERENCES hashtags(id),
			project_id INT REFERENCES projects(id)
		)`,
		`CREATE TABLE IF NOT EXISTS user_projects (
			project_id INT REFERENCES projects(id),
			user_id INT REFERENCES users(id)
		)`,
	}

	for _, query := range queries {
		_, err := db.Exec(query)
		if err != nil {
			fmt.Println("Error updating table schema.", err)
			return
		}
	}

	fmt.Println("Connection to Database Successfull")
}
