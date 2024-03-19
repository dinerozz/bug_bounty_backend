package db

import (
	"context"
	"fmt"
	"github.com/dinerozz/bug_bounty_backend/pkg/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

var Pool *pgxpool.Pool

func ConnectToDB(dsn string) *pgxpool.Pool {
	var err error
	Pool, err = pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	fmt.Println("Successfully connected to database!")
	return Pool
}

func GetUserByID(dbPool *pgxpool.Pool, userID uuid.UUID) (*models.UserById, error) {
	var user models.UserById

	err := dbPool.QueryRow(context.Background(), "SELECT id, username, email FROM users WHERE id = $1", userID).Scan(&user.ID, &user.Username, &user.Email)
	if err != nil {
		return nil, fmt.Errorf("error fetching user from database: %w", err)
	}

	return &user, nil
}

func Close() {
	if Pool != nil {
		Pool.Close()
	}
}
