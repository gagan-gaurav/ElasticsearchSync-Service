package handlers

import (
	"database/sql"
	"encoding/json"
	"fold/internal/models"
	"fold/internal/repository"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	// Parse request data
	var newUser models.User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Insert user into the database
	err = repository.CreateUser(&newUser)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to create new user")
		return
	}

	// Respond with success message
	RespondWithJSON(w, http.StatusCreated, map[string]string{"message": "User created successfully"})
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	// Get user ID from URL parameters
	vars := mux.Vars(r)
	userIDStr := vars["id"]
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	// Query the database for the user
	var user models.User
	err = repository.GetUserById(userID, &user)
	if err != nil {
		if err == sql.ErrNoRows {
			RespondWithError(w, http.StatusNotFound, "User not found")
		} else {
			RespondWithError(w, http.StatusInternalServerError, "Failed to fetch user")
		}
		return
	}

	// Respond with the user's information
	RespondWithJSON(w, http.StatusOK, user)
}

func CreateHashtag(w http.ResponseWriter, r *http.Request) {
	// Parse request data
	var newHashtag models.Hashtag
	err := json.NewDecoder(r.Body).Decode(&newHashtag)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Insert hashtag into the database
	err = repository.CreateHashtag(&newHashtag)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to create new hashtag")
		return
	}

	// Respond with success message
	RespondWithJSON(w, http.StatusCreated, map[string]string{"message": "Hashtag created successfully"})
}

func GetHashtag(w http.ResponseWriter, r *http.Request) {
	// Get hashtag ID from URL parameters
	vars := mux.Vars(r)
	hashtagIDStr := vars["id"]
	hashtagID, err := strconv.Atoi(hashtagIDStr)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid hashtag ID")
		return
	}

	// Query the database for the hashtag
	var hashtag models.Hashtag
	err = repository.GetHashtagById(hashtagID, &hashtag)
	if err != nil {
		if err == sql.ErrNoRows {
			RespondWithError(w, http.StatusNotFound, "Hashtag not found")
		} else {
			RespondWithError(w, http.StatusInternalServerError, "Failed to fetch hashtag")
		}
		return
	}

	// Respond with the hashtag's information
	RespondWithJSON(w, http.StatusOK, hashtag)
}

func CreateProject(w http.ResponseWriter, r *http.Request) {
	// Parse request data
	var newProject models.Project
	err := json.NewDecoder(r.Body).Decode(&newProject)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Check if all users are valid.
	users := newProject.UserIds
	for _, userId := range users {
		if !repository.UserExists(userId) {
			RespondWithError(w, http.StatusInternalServerError, "Users do not Exist.")
			return
		}
	}

	// Check if all hashtags are valid.
	hashtags := newProject.HashtagIds
	for _, hashtagId := range hashtags {
		if !repository.HashtagExistsById(hashtagId) {
			RespondWithError(w, http.StatusInternalServerError, "Hashtags do not Exist.")
			return
		}
	}

	//Start project creation transaction to insert project into database.
	err = repository.ProjectCreationAndSyncTransaction(&newProject)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to create new project. Transaction failed.")
		return
	}

	// Respond with success message
	RespondWithJSON(w, http.StatusCreated, map[string]string{"message": "Project created successfully"})
}

func GetProject(w http.ResponseWriter, r *http.Request) {
	// Get project ID from URL parameters
	vars := mux.Vars(r)
	projectIDStr := vars["id"]
	projectID, err := strconv.Atoi(projectIDStr)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid project ID")
		return
	}

	// Query the database for the project
	var project models.Project
	err = repository.GetProjectById(projectID, &project)
	if err != nil {
		if err == sql.ErrNoRows {
			RespondWithError(w, http.StatusNotFound, "Project not found")
		} else {
			RespondWithError(w, http.StatusInternalServerError, "Failed to fetch project")
		}
		return
	}

	// Respond with the project's information
	RespondWithJSON(w, http.StatusOK, project)
}

func RespondWithError(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payload)
}
