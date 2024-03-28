package main

import (
	"fmt"
	"github.com/dinerozz/bug_bounty_backend/cmd/migrate"
	db "github.com/dinerozz/bug_bounty_backend/config"
	"github.com/dinerozz/bug_bounty_backend/pkg/auth"
	"github.com/dinerozz/bug_bounty_backend/pkg/report"
	"github.com/dinerozz/bug_bounty_backend/pkg/role"
	"github.com/dinerozz/bug_bounty_backend/pkg/task"
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
	router.GET("/refresh", auth.RefreshHandler)

	authRequired := router.Group("/")
	authRequired.Use(auth.JWTMiddleware())
	{
		authRequired.POST("/logout", auth.LogoutHandler)
		authRequired.GET("/current", auth.CurrentUserHandler)
		authRequired.POST("/team", team.CreateTeamHandler)
		authRequired.GET("/teams", team.GetTeamsHandler)
		authRequired.PATCH("/team/invite-token", team.UpdateInviteTokenHandler)
		authRequired.POST("/team/join", team.JoinTeamHandler)
		authRequired.GET("/team/members", team.GetTeamMembersHandler)
		authRequired.GET("/my-team", team.GetTeamHandler)
		authRequired.GET("/tasks", task.GetTasksHandler)
		authRequired.POST("/report", report.CreateReportHandler)
		authRequired.GET("/report", report.GetReportsHandler)
		authRequired.GET("/report/details", report.ReviewDetailsHandler)
	}

	adminRoutes := router.Group("/admin")
	adminRoutes.Use(auth.JWTMiddleware(), role.RolesMiddleware("ADMIN"))
	{
		adminRoutes.POST("/role", role.CreateRoleHandler)
		adminRoutes.POST("/user/role", role.SetUserRoleHandler)
		adminRoutes.POST("/tasks", task.CreateTaskHandler)
		adminRoutes.POST("/report/review", report.ReviewReportHandler)
	}

	log.Println("Запуск сервера на http://localhost:5555")
	router.Run("localhost:5555")
}
