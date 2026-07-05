package repository

import (
	"context"
	"fmt"

	"reservation-system/backend/internal/domain"
)

func (r *Repository) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	const query = `
		INSERT INTO users (
			name,
			email,
			password_hash,
			role,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING
			id,
			name,
			email,
			password_hash,
			role,
			created_at,
			updated_at
	`

	createdUser := &domain.User{}
	if err := r.db.QueryRow(
		ctx,
		query,
		user.Name,
		user.Email,
		user.PasswordHash,
		user.Role,
	).Scan(
		&createdUser.ID,
		&createdUser.Name,
		&createdUser.Email,
		&createdUser.PasswordHash,
		&createdUser.Role,
		&createdUser.CreatedAt,
		&createdUser.UpdatedAt,
	); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	return createdUser, nil
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	const query = `
		SELECT
			id,
			name,
			email,
			password_hash,
			role,
			created_at,
			updated_at
		FROM users
		WHERE email = $1
	`

	user := &domain.User{}
	if err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		return nil, fmt.Errorf("get user by email: %w", err)
	}

	return user, nil
}
