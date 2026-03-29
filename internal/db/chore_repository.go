package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

// ChoreRepository handles chore database operations
type ChoreRepository struct {
	conn *pgx.Conn
}

// NewChoreRepository creates a new chore repository
func NewChoreRepository(conn *pgx.Conn) *ChoreRepository {
	return &ChoreRepository{conn: conn}
}

// CreateChore creates a new chore
func (r *ChoreRepository) CreateChore(ctx context.Context, name, description, frequency string, pointsValue int) (*Chore, error) {
	const query = `
		INSERT INTO chores (name, description, frequency, points_value)
		VALUES ($1, $2, $3, $4)
		RETURNING id, name, description, frequency, points_value, created_at, updated_at
	`

	var chore Chore
	err := r.conn.QueryRow(ctx, query, name, description, frequency, pointsValue).Scan(
		&chore.ID,
		&chore.Name,
		&chore.Description,
		&chore.Frequency,
		&chore.PointsValue,
		&chore.CreatedAt,
		&chore.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create chore: %w", err)
	}

	return &chore, nil
}

// DeleteChore deletes a chore
func (r *ChoreRepository) DeleteChore(ctx context.Context, name string) error {
	const query = `
	DELETE FROM chores WHERE name = $1
	`

	_, err := r.conn.Exec(ctx, query, name)
	if err != nil {
		return fmt.Errorf("failed to delete chore: %w", err)
	}

	return nil
}

// GetChoreByID retrieves a chore by ID
func (r *ChoreRepository) GetChoreByID(ctx context.Context, id int) (*Chore, error) {
	const query = `
		SELECT id, name, description, frequency, points_value, created_at, updated_at
		FROM chores
		WHERE id = $1
	`

	var chore Chore
	err := r.conn.QueryRow(ctx, query, id).Scan(
		&chore.ID,
		&chore.Name,
		&chore.Description,
		&chore.Frequency,
		&chore.PointsValue,
		&chore.CreatedAt,
		&chore.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("chore not found")
		}
		return nil, fmt.Errorf("failed to get chore: %w", err)
	}

	return &chore, nil
}

// GetAllChores retrieves all chores
func (r *ChoreRepository) GetAllChores(ctx context.Context) ([]*Chore, error) {
	const query = `
		SELECT id, name, description, frequency, points_value, created_at, updated_at
		FROM chores
		ORDER BY created_at DESC
	`

	rows, err := r.conn.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query chores: %w", err)
	}
	defer rows.Close()

	var chores []*Chore
	for rows.Next() {
		var chore Chore
		err := rows.Scan(
			&chore.ID,
			&chore.Name,
			&chore.Description,
			&chore.Frequency,
			&chore.PointsValue,
			&chore.CreatedAt,
			&chore.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan chore: %w", err)
		}
		chores = append(chores, &chore)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error reading chores: %w", err)
	}

	return chores, nil
}

// AssignChoreToUser assigns a chore to a user
func (r *ChoreRepository) AssignChoreToUser(ctx context.Context, userID int, choreID int) error {
	const query = `
		INSERT INTO user_chores (user_id, chore_id)
		VALUES ($1, $2)
	`

	_, err := r.conn.Exec(ctx, query, userID, choreID)
	if err != nil {
		return fmt.Errorf("failed to assign chore to user: %w", err)
	}

	return nil
}

// GetChoresForUser retrieves all chores assigned to a user
func (r *ChoreRepository) GetChoresForUser(ctx context.Context, userID int) ([]*Chore, error) {
	const query = `
		SELECT c.id, c.name, c.description, c.frequency, c.points_value, c.created_at, c.updated_at
		FROM chores c
		INNER JOIN user_chores uc ON c.id = uc.chore_id
		WHERE uc.user_id = $1
		ORDER BY c.name
	`

	rows, err := r.conn.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query user chores: %w", err)
	}
	defer rows.Close()

	var chores []*Chore
	for rows.Next() {
		var chore Chore
		err := rows.Scan(
			&chore.ID,
			&chore.Name,
			&chore.Description,
			&chore.Frequency,
			&chore.PointsValue,
			&chore.CreatedAt,
			&chore.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan chore: %w", err)
		}
		chores = append(chores, &chore)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error reading chores: %w", err)
	}

	return chores, nil
}

// RecordCompletion records a chore completion
func (r *ChoreRepository) RecordCompletion(ctx context.Context, userID int, choreID int) (*Completion, error) {
	const query = `
		INSERT INTO completions (user_id, chore_id, points_earned)
		SELECT $1, $2, points_value
		FROM chores
		WHERE id = $2
		RETURNING id, user_id, chore_id, completed_at, points_earned
	`

	var completion Completion
	err := r.conn.QueryRow(ctx, query, userID, choreID).Scan(
		&completion.ID,
		&completion.UserID,
		&completion.ChoreID,
		&completion.CompletedAt,
		&completion.PointsEarned,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to record completion: %w", err)
	}

	return &completion, nil
}

// GetRecentCompletions retrieves recent completions
func (r *ChoreRepository) GetRecentCompletions(ctx context.Context, limit int) ([]*Completion, error) {
	const query = `
		SELECT id, user_id, chore_id, completed_at, points_earned
		FROM completions
		ORDER BY completed_at DESC
		LIMIT $1
	`

	rows, err := r.conn.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query completions: %w", err)
	}
	defer rows.Close()

	var completions []*Completion
	for rows.Next() {
		var completion Completion
		err := rows.Scan(
			&completion.ID,
			&completion.UserID,
			&completion.ChoreID,
			&completion.CompletedAt,
			&completion.PointsEarned,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan completion: %w", err)
		}
		completions = append(completions, &completion)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error reading completions: %w", err)
	}

	return completions, nil
}
