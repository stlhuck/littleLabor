package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/stlhuck/littleLabor/internal/db"
)

// RewardHandler handles reward-related HTTP requests
type RewardHandler struct {
	repo *db.RewardRepository
}

// NewRewardHandler creates a new reward handler
func NewRewardHandler(repo *db.RewardRepository) *RewardHandler {
	return &RewardHandler{repo: repo}
}

// List returns all rewards
func (h *RewardHandler) List(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	rewards, err := h.repo.GetAllRewards(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if rewards == nil {
		rewards = []*db.Reward{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"rewards": rewards})
}
