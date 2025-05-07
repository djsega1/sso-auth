package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/djsega1/sso-auth/config"
	"golang.org/x/crypto/bcrypt"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func InitDB(cfg *config.Config) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s&timezone=%s",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.DBSSLMode, cfg.DBTimeZone)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("Failed to create connection pool: %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	DB = pool
	fmt.Println("Database connection established successfully.")

	_, err = DB.Exec(ctx, `CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL
	)`)
	if err != nil {
		log.Fatalf("Failed to run migration: %v", err)
	}
}

func RegisterUser(username, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = DB.Exec(ctx, "INSERT INTO users (username, password) VALUES ($1, $2)", username, string(hashedPassword))
	if err != nil {
		return fmt.Errorf("failed to save user: %v", err)
	}

	return nil
}

func AuthenticateUser(username, password string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var hashedPassword string
	err := DB.QueryRow(ctx, "SELECT password FROM users WHERE username = $1", username).Scan(&hashedPassword)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return false, fmt.Errorf("user not found")
		}
		return false, fmt.Errorf("database query error: %v", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return false, nil // incorrect password
	}

	return true, nil
}
