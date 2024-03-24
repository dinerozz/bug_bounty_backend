package db

import (
	"context"
	"github.com/jackc/pgx/v4"
	"log"
)

type LogLevel int

const (
	LogLevelTrace LogLevel = iota
	LogLevelDebug
	LogLevelInfo
	LogLevelWarn
	LogLevelError
	LogLevelNone
)

type logger struct{}

func (l *logger) Log(ctx context.Context, level pgx.LogLevel, msg string, data map[string]interface{}) {
	switch level {
	case pgx.LogLevelTrace, pgx.LogLevelDebug:
		log.Printf("DEBUG: %s: %+v\n", msg, data)
	case pgx.LogLevelInfo:
		log.Printf("INFO: %s: %+v\n", msg, data)
	case pgx.LogLevelWarn:
		log.Printf("WARN: %s: %+v\n", msg, data)
	case pgx.LogLevelError:
		log.Printf("ERROR: %s: %+v\n", msg, data)
	}
}
