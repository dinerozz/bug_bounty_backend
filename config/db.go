package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

// TODO REFACTOR (DI)
var Pool *pgxpool.Pool

func ConnectToDB(dsn string) *pgxpool.Pool {
	var err error

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatalf("Unable to parse database configuration: %v\n", err)
	}

	config.ConnConfig.Logger = &logger{}

	config.ConnConfig.LogLevel = pgx.LogLevelDebug

	Pool, err = pgxpool.ConnectConfig(context.Background(), config)
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
