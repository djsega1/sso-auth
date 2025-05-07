package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(username, hashedPassword string) (uuid.UUID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var id uuid.UUID
	err := r.db.QueryRow(ctx,
		"INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id",
		username, hashedPassword).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to save user: %v", err)
	}

	return id, nil
}

func (r *UserRepository) GetHashedPassword(username string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var hashedPassword string
	err := r.db.QueryRow(ctx,
		"SELECT password FROM users WHERE username = $1",
		username).Scan(&hashedPassword)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return "", fmt.Errorf("user not found")
		}
		return "", fmt.Errorf("database query error: %v", err)
	}

	return hashedPassword, nil
}
