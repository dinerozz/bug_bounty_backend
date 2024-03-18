package migrate

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
)

//go:embed migrations/*.sql
var migrations embed.FS

func RunMigrations(dbPool *pgxpool.Pool, migrationsPath string) {
	migrationFiles, err := fs.ReadDir(migrations, migrationsPath)
	if err != nil {
		log.Fatal("Не удалось прочитать директорию с миграциями:", err)
	}

	for _, file := range migrationFiles {
		fileName := file.Name()
		fileContent, err := fs.ReadFile(migrations, fmt.Sprintf("%s/%s", migrationsPath, fileName))
		if err != nil {
			log.Fatal("Ошибка при чтении файла миграции:", err)
		}

		_, err = dbPool.Exec(context.Background(), string(fileContent))
		if err != nil {
			log.Fatalf("Ошибка при выполнении миграции %s: %v", fileName, err)
		}

		log.Printf("Миграция %s успешно применена", fileName)
	}
}
