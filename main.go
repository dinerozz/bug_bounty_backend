package main

import (
	"fmt"
	"github.com/dinerozz/bug_bounty_backend/cmd/migrate"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Привет, мир!")
}

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

	fmt.Println(dsn)

	// TODO: функция для создания коннекта к базе данных, в migrate прокидывать
	db := migrate.ConnectToDB(dsn)

	migrate.RunMigrations(db, "cmd/migrate/migrations/000001_init_schema.up.sql")

	defer db.Close()
	fmt.Println("Successfully connected to database!")

	http.HandleFunc("/", handler)

	log.Println("Запуск сервера на http://localhost:8080")
	log.Fatal(http.ListenAndServe("localhost:8081", nil))
}
