package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
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
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	// Insert user into the database
	err = repository.CreateUser(&newUser)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to create new user", err)
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
		RespondWithError(w, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	// Query the database for the user
	var user models.User
	err = repository.GetUserById(userID, &user)
	if err != nil {
		if err == sql.ErrNoRows {
			RespondWithError(w, http.StatusNotFound, "User not found", err)
		} else {
			RespondWithError(w, http.StatusInternalServerError, "Failed to fetch user", err)
		}
		return
	}

	// Respond with the user's information
	RespondWithJSON(w, http.StatusOK, user)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	// Get user ID from URL parameters
	vars := mux.Vars(r)
	userIDStr := vars["id"]
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	// Parse request data
	var updatedUser models.User
	err = json.NewDecoder(r.Body).Decode(&updatedUser)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	// Perform transaction to update user in the database
	updatedUser.ID = userID // Set the ID for the user to be updated
	err = repository.UpdateUserTransaction(&updatedUser)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to update user", err)
		return
	}

	// Respond with success message
	RespondWithJSON(w, http.StatusOK, map[string]string{"message": "User updated successfully"})
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	// Get user ID from URL parameters
	vars := mux.Vars(r)
	userIDStr := vars["id"]
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	// Perform transaction to delete user in the database
	err = repository.DeleteUserTransaction(userID)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to delete user", err)
		return
	}

	// Respond with success message
	RespondWithJSON(w, http.StatusOK, map[string]string{"message": "User deleted successfully"})
}

func CreateHashtag(w http.ResponseWriter, r *http.Request) {
	// Parse request data
	var newHashtag models.Hashtag
	err := json.NewDecoder(r.Body).Decode(&newHashtag)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	// Insert hashtag into the database
	err = repository.CreateHashtag(&newHashtag)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to create new hashtag", err)
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
		RespondWithError(w, http.StatusBadRequest, "Invalid hashtag ID", err)
		return
	}

	// Query the database for the hashtag
	var hashtag models.Hashtag
	err = repository.GetHashtagById(hashtagID, &hashtag)
	if err != nil {
		if err == sql.ErrNoRows {
			RespondWithError(w, http.StatusNotFound, "Hashtag not found", err)
		} else {
			RespondWithError(w, http.StatusInternalServerError, "Failed to fetch hashtag", err)
		}
		return
	}

	// Respond with the hashtag's information
	RespondWithJSON(w, http.StatusOK, hashtag)
}

func UpdateHashtag(w http.ResponseWriter, r *http.Request) {
	// Get hashtag ID from URL parameters
	vars := mux.Vars(r)
	hashtagIDStr := vars["id"]
	hashtagID, err := strconv.Atoi(hashtagIDStr)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid hashtag ID", err)
		return
	}

	// Parse request data
	var updatedHashtag models.Hashtag
	err = json.NewDecoder(r.Body).Decode(&updatedHashtag)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	// Perform transaction to update hashtags in the database
	updatedHashtag.ID = hashtagID // Set the ID for the hashtag to be updated
	err = repository.UpdateHashtagTransaction(&updatedHashtag)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to update hashtag", err)
		return
	}

	// Respond with success message
	RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Hashtag updated successfully"})
}

func DeleteHashtag(w http.ResponseWriter, r *http.Request) {
	// Get hashtag ID from URL parameters
	vars := mux.Vars(r)
	hashtagIDStr := vars["id"]
	hashtagID, err := strconv.Atoi(hashtagIDStr)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid hashtag ID", err)
		return
	}

	// Perform transaction to delete hashtag in the database
	err = repository.DeleteHashtagTransaction(hashtagID)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to delte hashtag", err)
		return
	}

	// Respond with success message
	RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Hashtag deleted successfully"})
}

func CreateProject(w http.ResponseWriter, r *http.Request) {
	// Parse request data
	var newProject models.Project
	err := json.NewDecoder(r.Body).Decode(&newProject)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	// Check if all users are valid.
	users := newProject.UserIds
	for _, userId := range users {
		if !repository.UserExists(userId) {
			RespondWithError(w, http.StatusInternalServerError, "Users do not Exist.", err)
			return
		}
	}

	// Check if all hashtags are valid.
	hashtags := newProject.HashtagIds
	for _, hashtagId := range hashtags {
		if !repository.HashtagExists(hashtagId) {
			RespondWithError(w, http.StatusInternalServerError, "Hashtags do not Exist.", err)
			return
		}
	}

	//Start project creation transaction to insert project into database.
	err = repository.ProjectCreationAndSyncTransaction(&newProject)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to create new project. Transaction failed.", err)
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
		RespondWithError(w, http.StatusBadRequest, "Invalid project ID", err)
		return
	}

	// Query the database for the project
	var project models.Project
	err = repository.GetProjectById(projectID, &project)
	if err != nil {
		if err == sql.ErrNoRows {
			RespondWithError(w, http.StatusNotFound, "Project not found", err)
		} else {
			RespondWithError(w, http.StatusInternalServerError, "Failed to fetch project", err)
		}
		return
	}

	// Respond with the project's information
	RespondWithJSON(w, http.StatusOK, project)
}

func UpdateProject(w http.ResponseWriter, r *http.Request) {
	// Get project ID from URL parameters
	vars := mux.Vars(r)
	projectIDStr := vars["id"]
	projectID, err := strconv.Atoi(projectIDStr)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid project ID", err)
		return
	}

	// Parse request data
	var newProject models.Project
	err = json.NewDecoder(r.Body).Decode(&newProject)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	// Check if all users are valid.
	users := newProject.UserIds
	for _, userId := range users {
		if !repository.UserExists(userId) {
			RespondWithError(w, http.StatusInternalServerError, "Users do not Exist.", err)
			return
		}
	}

	// Check if all hashtags are valid.
	hashtags := newProject.HashtagIds
	for _, hashtagId := range hashtags {
		if !repository.HashtagExists(hashtagId) {
			RespondWithError(w, http.StatusInternalServerError, "Hashtags do not Exist.", err)
			return
		}
	}

	newProject.ID = projectID // set project id

	//Start project update transaction to update project into database.
	err = repository.ProjectUpdateAndSyncTransaction(&newProject)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to update project. Update Transaction failed.", err)
		return
	}

	// Respond with success message
	RespondWithJSON(w, http.StatusCreated, map[string]string{"message": "Project updated successfully"})
}

func DeleteProject(w http.ResponseWriter, r *http.Request) {
	// Get project ID from URL parameters
	vars := mux.Vars(r)
	projectIDStr := vars["id"]
	projectID, err := strconv.Atoi(projectIDStr)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid project ID", err)
		return
	}

	//Start project Delete transaction to delete project into database.
	err = repository.ProjectDeleteAndSyncTransaction(projectID)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to delete project. Delete Transaction failed.", err)
		return
	}

	// Respond with success message
	RespondWithJSON(w, http.StatusCreated, map[string]string{"message": "Project deleted successfully"})
}

func RespondWithError(w http.ResponseWriter, code int, message string, err error) {
	fmt.Println(err)
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payload)
}
