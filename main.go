package main

import (
	"fmt"
	"github.com/dinerozz/bug_bounty_backend/cmd/migrate"
	db "github.com/dinerozz/bug_bounty_backend/config"
	"github.com/dinerozz/bug_bounty_backend/pkg/auth"
	"github.com/dinerozz/bug_bounty_backend/pkg/team"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	router := gin.Default()
	router.Use(auth.CORSMiddleware())

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

	router.POST("/register", auth.RegisterHandler)
	router.POST("/authenticate", auth.AuthenticateHandler)

	authRequired := router.Group("/")
	authRequired.Use(auth.JWTMiddleware())
	{
		authRequired.GET("/refresh", auth.RefreshHandler)
		authRequired.GET("/current", auth.CurrentUserHandler)
		authRequired.POST("/team", team.CreateTeamHandler)
	}

	log.Println("Запуск сервера на http://localhost:5555")
	router.Run(":5555")
}
