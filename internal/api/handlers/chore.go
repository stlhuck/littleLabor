package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/stlhuck/littleLabor/internal/db"
)

// ChoreHandler handles chore-related HTTP requests
type ChoreHandler struct {
	repo *db.ChoreRepository
}

// NewChoreHandler creates a new chore handler
func NewChoreHandler(repo *db.ChoreRepository) *ChoreHandler {
	return &ChoreHandler{repo: repo}
}

type CreateChoreRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	Frequency   string `json:"frequency" binding:"required"`
	PointsValue int    `json:"points_value" binding:"required"`
}

type DeleteChoreRequest struct {
	Name string `json:"name" binding:"required"`
}

// List returns all chores
func (h *ChoreHandler) List(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	chores, err := h.repo.GetAllChores(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if chores == nil {
		chores = []*db.Chore{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"chores": chores})
}

// Create creates a new chore
func (h *ChoreHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateChoreRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.Description == "" || req.Frequency == "" || req.PointsValue == 0 {
		http.Error(w, "Missing required fields: name, description, frequency, points_value", http.StatusBadRequest)
		return
	}

	chore, err := h.repo.CreateChore(r.Context(), req.Name, req.Description, req.Frequency, req.PointsValue)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(chore)
}

// Delete deletes a chore
func (h *ChoreHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req DeleteChoreRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "Missing required fields: name", http.StatusBadRequest)
		return
	}

	err := h.repo.DeleteChore(r.Context(), req.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)

}
