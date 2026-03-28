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
