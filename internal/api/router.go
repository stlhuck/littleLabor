package api

import (
	"net/http"

	"github.com/stlhuck/littleLabor/internal/api/handlers"
	"github.com/stlhuck/littleLabor/internal/db"
)

// Router sets up all application routes
type Router struct {
	userHandler   *handlers.UserHandler
	choreHandler  *handlers.ChoreHandler
	rewardHandler *handlers.RewardHandler
}

// NewRouter creates a new router with all handlers
func NewRouter(userRepo *db.UserRepository, choreRepo *db.ChoreRepository, rewardRepo *db.RewardRepository) *Router {
	return &Router{
		userHandler:   handlers.NewUserHandler(userRepo),
		choreHandler:  handlers.NewChoreHandler(choreRepo),
		rewardHandler: handlers.NewRewardHandler(rewardRepo),
	}
}

// RegisterRoutes registers all HTTP route handlers
func (r *Router) RegisterRoutes() {
	// Health checks
	http.HandleFunc("/health", r.handleHealth)
	http.HandleFunc("/", r.handleRoot)

	// User routes
	http.HandleFunc("/api/users", r.handleUsers)

	// Chore routes
	http.HandleFunc("/api/chores", r.handleChores)

	// Reward routes
	http.HandleFunc("/api/rewards", r.handleRewards)
}

func (r *Router) handleRoot(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status": "littleLabor API is running", "version": "0.1.0"}`))
}

func (r *Router) handleHealth(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status": "healthy"}`))
}

func (r *Router) handleUsers(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		r.userHandler.List(w, req)
	case http.MethodPost:
		r.userHandler.Create(w, req)
	case http.MethodDelete:
		r.userHandler.Delete(w, req)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
func (r *Router) handleChores(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		r.choreHandler.List(w, req)
	case http.MethodPost:
		r.choreHandler.Create(w, req)
	case http.MethodDelete:
		r.choreHandler.Delete(w, req)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (r *Router) handleRewards(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		r.rewardHandler.List(w, req)
	case http.MethodPost:
		r.rewardHandler.Create(w, req)
	case http.MethodDelete:
		r.rewardHandler.Delete(w, req)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
