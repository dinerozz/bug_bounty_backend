package main

import (
	"fmt"
	"github.com/dinerozz/bug_bounty_backend/cmd/migrate"
	db "github.com/dinerozz/bug_bounty_backend/config"
	"github.com/dinerozz/bug_bounty_backend/pkg/auth"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func main() {
	godotenv.Load(".env")
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	dbUser := os.Getenv("DATABASE_USER")
	dbPassword := os.Getenv("DATABASE_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DATABASE_NAME")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	dbPool := db.ConnectToDB(dsn)
	defer dbPool.Close()

	migrate.RunMigrations(dbPool, "migrations")

	defer db.Close()
	fmt.Println("Successfully connected to database!")

	mux := http.NewServeMux()

	mux.HandleFunc("/register", auth.RegisterHandler)
	mux.HandleFunc("/authenticate", auth.AuthenticateHandler)
	mux.Handle("/current", auth.JWTMiddleware(http.HandlerFunc(auth.CurrentUserHandler)))

	handler := auth.CORSMiddleware(mux)

	log.Println("Запуск сервера на http://localhost:5555")
	log.Fatal(http.ListenAndServe("localhost:5555", handler))
	//log.Fatal(http.ListenAndServe("0.0.0.0:5555", nil))
}
