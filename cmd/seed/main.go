package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"dispatch/internal/platform/auth"
	"dispatch/internal/platform/config"
)

// This seed creates a default admin user if it does not exist
// Username: admin
// Password: admin123

func main() {
	ctx := context.Background()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	pool, err := pgxpool.New(ctx, cfg.DB.DSN())
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	username := "admin"
	password := "admin123"

	var exists bool
	err = pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM users WHERE username=$1 AND deleted_at IS NULL)`, username).Scan(&exists)
	if err != nil {
		log.Fatal(err)
	}

	if exists {
		fmt.Println("admin user already exists")
		os.Exit(0)
	}

	hash, err := auth.HashPassword(password)
	if err != nil {
		log.Fatal(err)
	}

	userID := uuid.NewString()
	now := time.Now().UTC()

	tx, err := pool.Begin(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `
	INSERT INTO users (
		id, username, first_name, last_name, email, phone,
		password_hash, status, is_active, created_at, updated_at
	)
	VALUES ($1,$2,'System','Administrator','admin@dispatch.local','+256780000000', $3,'ACTIVE',true,$4,$4)
	`, userID, username, hash, now)
	if err != nil {
		log.Fatal(err)
	}

	var adminRoleID string
	err = tx.QueryRow(ctx, `SELECT id FROM roles WHERE name='Super Admin'`).Scan(&adminRoleID)
	if err != nil {
		log.Fatal("admin role not found. Ensure RBAC seed ran first.")
	}

	_, err = tx.Exec(ctx, `INSERT INTO user_roles (user_id, role_id) VALUES ($1,$2)`, userID, adminRoleID)
	if err != nil {
		log.Fatal(err)
	}

	if err := tx.Commit(ctx); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Default admin user created")
	fmt.Println("username: admin")
	fmt.Println("password: admin123")
}
