package db

import (
	"context"
	"fmt"
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

func Close() {
	if Pool != nil {
		Pool.Close()
	}
}
