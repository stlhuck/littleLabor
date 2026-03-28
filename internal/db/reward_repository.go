package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

// RewardRepository handles reward database operations
type RewardRepository struct {
	conn *pgx.Conn
}

// NewRewardRepository creates a new reward repository
func NewRewardRepository(conn *pgx.Conn) *RewardRepository {
	return &RewardRepository{conn: conn}
}

// CreateReward creates a new reward
func (r *RewardRepository) CreateReward(ctx context.Context, name, description string, pointCost int) (*Reward, error) {
	const query = `
		INSERT INTO rewards (name, description, point_cost)
		VALUES ($1, $2, $3)
		RETURNING id, name, description, point_cost, created_at, updated_at
	`

	var reward Reward
	err := r.conn.QueryRow(ctx, query, name, description, pointCost).Scan(
		&reward.ID,
		&reward.Name,
		&reward.Description,
		&reward.PointCost,
		&reward.CreatedAt,
		&reward.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create reward: %w", err)
	}

	return &reward, nil
}

// GetRewardByID retrieves a reward by ID
func (r *RewardRepository) GetRewardByID(ctx context.Context, id int) (*Reward, error) {
	const query = `
		SELECT id, name, description, point_cost, created_at, updated_at
		FROM rewards
		WHERE id = $1
	`

	var reward Reward
	err := r.conn.QueryRow(ctx, query, id).Scan(
		&reward.ID,
		&reward.Name,
		&reward.Description,
		&reward.PointCost,
		&reward.CreatedAt,
		&reward.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("reward not found")
		}
		return nil, fmt.Errorf("failed to get reward: %w", err)
	}

	return &reward, nil
}

// GetAllRewards retrieves all rewards
func (r *RewardRepository) GetAllRewards(ctx context.Context) ([]*Reward, error) {
	const query = `
		SELECT id, name, description, point_cost, created_at, updated_at
		FROM rewards
		ORDER BY point_cost ASC
	`

	rows, err := r.conn.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query rewards: %w", err)
	}
	defer rows.Close()

	var rewards []*Reward
	for rows.Next() {
		var reward Reward
		err := rows.Scan(
			&reward.ID,
			&reward.Name,
			&reward.Description,
			&reward.PointCost,
			&reward.CreatedAt,
			&reward.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan reward: %w", err)
		}
		rewards = append(rewards, &reward)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error reading rewards: %w", err)
	}

	return rewards, nil
}

// RedeemReward marks a reward as redeemed for a user
func (r *RewardRepository) RedeemReward(ctx context.Context, userID int, rewardID int) error {
	const query = `
		UPDATE user_rewards
		SET redeemed_at = CURRENT_TIMESTAMP, status = 'redeemed'
		WHERE user_id = $1 AND reward_id = $2 AND status = 'earned'
	`

	result, err := r.conn.Exec(ctx, query, userID, rewardID)
	if err != nil {
		return fmt.Errorf("failed to redeem reward: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("reward not found or already redeemed")
	}

	return nil
}

// EarnReward marks a reward as earned for a user
func (r *RewardRepository) EarnReward(ctx context.Context, userID int, rewardID int) (*UserReward, error) {
	const query = `
		INSERT INTO user_rewards (user_id, reward_id, status)
		VALUES ($1, $2, 'earned')
		RETURNING id, user_id, reward_id, earned_at, redeemed_at, status
	`

	var ur UserReward
	err := r.conn.QueryRow(ctx, query, userID, rewardID).Scan(
		&ur.ID,
		&ur.UserID,
		&ur.RewardID,
		&ur.EarnedAt,
		&ur.RedeemedAt,
		&ur.Status,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to earn reward: %w", err)
	}

	return &ur, nil
}

// GetUserRewards retrieves all rewards for a user
func (r *RewardRepository) GetUserRewards(ctx context.Context, userID int) ([]*UserReward, error) {
	const query = `
		SELECT id, user_id, reward_id, earned_at, redeemed_at, status
		FROM user_rewards
		WHERE user_id = $1
		ORDER BY earned_at DESC
	`

	rows, err := r.conn.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query user rewards: %w", err)
	}
	defer rows.Close()

	var rewards []*UserReward
	for rows.Next() {
		var ur UserReward
		err := rows.Scan(
			&ur.ID,
			&ur.UserID,
			&ur.RewardID,
			&ur.EarnedAt,
			&ur.RedeemedAt,
			&ur.Status,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan reward: %w", err)
		}
		rewards = append(rewards, &ur)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error reading rewards: %w", err)
	}

	return rewards, nil
}
