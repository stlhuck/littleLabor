package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/stlhuck/littleLabor/internal/db"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	repo *db.UserRepository
}

// NewUserHandler creates a new user handler
func NewUserHandler(repo *db.UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

// CreateUserRequest represents the JSON payload for creating a user
type CreateUserRequest struct {
	Username     string `json:"username" binding:"required"`
	Email        string `json:"email" binding:"required"`
	PasswordHash string `json:"password_hash" binding:"required"`
}

// List returns all users
func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	users, err := h.repo.GetAllUsers(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if users == nil {
		users = []*db.User{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"users": users})
}

// Create creates a new user
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Email == "" || req.PasswordHash == "" {
		http.Error(w, "Missing required fields: username, email, password_hash", http.StatusBadRequest)
		return
	}

	user, err := h.repo.CreateUser(r.Context(), req.Username, req.Email, req.PasswordHash)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}
