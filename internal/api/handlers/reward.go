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

type CreateRewardRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	PointsCost  int    `json:"points_cost" binding:"required"`
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

func (h *RewardHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateRewardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.Description == "" || req.PointsCost == 0 {
		http.Error(w, "Missing required fields: name, description, point_cost", http.StatusBadRequest)
		return
	}

	reward, err := h.repo.CreateReward(r.Context(), req.Name, req.Description, req.PointsCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(reward)

}
