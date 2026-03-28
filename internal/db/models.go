package db

import "time"

// User represents a user in the system
type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Points    int       `json:"points"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Chore represents a chore task
type Chore struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Description string   `json:"description"`
	Frequency  string    `json:"frequency"`
	PointsValue int      `json:"points_value"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// Completion represents a completed chore
type Completion struct {
	ID           int       `json:"id"`
	UserID       int       `json:"user_id"`
	ChoreID      int       `json:"chore_id"`
	CompletedAt  time.Time `json:"completed_at"`
	PointsEarned int       `json:"points_earned"`
}

// Reward represents a reward users can earn
type Reward struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	PointCost   int       `json:"point_cost"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// UserReward represents a user's earned or redeemed reward
type UserReward struct {
	ID        int        `json:"id"`
	UserID    int        `json:"user_id"`
	RewardID  int        `json:"reward_id"`
	EarnedAt  time.Time  `json:"earned_at"`
	RedeemedAt *time.Time `json:"redeemed_at"`
	Status    string     `json:"status"`
}
