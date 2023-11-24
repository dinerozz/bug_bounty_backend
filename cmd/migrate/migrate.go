package migrate

import (
	"database/sql"
	"embed"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func ConnectToDB(dsn string) *sql.DB {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Ошибка при подключении к базе данных:", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Ошибка при проверке соединения с базой данных:", err)
	}

	fmt.Println("Успешное подключение к базе данных")
	return db
}

//go:embed migrations
var migrations embed.FS

func RunMigrations(db *sql.DB, filename string) {
	content, err := os.ReadFile(filename)

	if err != nil {
		log.Fatal(err)
	}

	requests := string(content)
	_, err = db.Exec(requests)

	if err != nil {
		log.Fatal(err)
	}

	log.Print("Migrations have been successfully applied")

}
