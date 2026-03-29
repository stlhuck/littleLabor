package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

// UserRepository handles user database operations
type UserRepository struct {
	conn *pgx.Conn
}

// NewUserRepository creates a new user repository
func NewUserRepository(conn *pgx.Conn) *UserRepository {
	return &UserRepository{conn: conn}
}

// CreateUser creates a new user in the database
func (r *UserRepository) CreateUser(ctx context.Context, username, email, passwordHash string) (*User, error) {
	const query = `
		INSERT INTO users (username, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id, username, email, points, created_at, updated_at
	`

	var user User
	err := r.conn.QueryRow(ctx, query, username, email, passwordHash).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Points,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &user, nil
}

// DeleteUser deletes a user from the database
func (r *UserRepository) DeleteUser(ctx context.Context, username string) error {
	const query = `
		DELETE FROM users WHERE username = $1
	`
	
	_, err := r.conn.Exec(ctx, query, username)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	
	return nil
}

// GetUserByID retrieves a user by ID
func (r *UserRepository) GetUserByID(ctx context.Context, id int) (*User, error) {
	const query = `
		SELECT id, username, email, points, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user User
	err := r.conn.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Points,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// GetUserByUsername retrieves a user by username
func (r *UserRepository) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	const query = `
		SELECT id, username, email, points, created_at, updated_at
		FROM users
		WHERE username = $1
	`

	var user User
	err := r.conn.QueryRow(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Points,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// GetAllUsers retrieves all users
func (r *UserRepository) GetAllUsers(ctx context.Context) ([]*User, error) {
	const query = `
		SELECT id, username, email, points, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
	`

	rows, err := r.conn.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.Points,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error reading users: %w", err)
	}

	return users, nil
}

// AddPointsToUser adds points to a user
func (r *UserRepository) AddPointsToUser(ctx context.Context, userID, points int) error {
	const query = `
		UPDATE users
		SET points = points + $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`

	result, err := r.conn.Exec(ctx, query, points, userID)
	if err != nil {
		return fmt.Errorf("failed to add points: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}
