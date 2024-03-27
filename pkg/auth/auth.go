package auth

import (
	"context"
	"fmt"
	db "github.com/dinerozz/bug_bounty_backend/config"
	"github.com/dinerozz/bug_bounty_backend/pkg/models"
	"github.com/dinerozz/bug_bounty_backend/pkg/role"
	_ "github.com/dinerozz/bug_bounty_backend/pkg/team"
	team2 "github.com/dinerozz/bug_bounty_backend/pkg/team"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
	"time"
)

type Claims struct {
	UserID uuid.UUID `json:"userId"`
	jwt.RegisteredClaims
}

func RegisterUser(username, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	newUserUUID := uuid.New()
	_, err = db.Pool.Exec(context.Background(), "INSERT into users (id, username, email, password) VALUES ($1, $2, $3, $4)", newUserUUID, username, email, hashedPassword)
	if err != nil {
		return fmt.Errorf("ошибка при добавлении пользователя в базу данных: %w", err)
	}

	var userRoleID uuid.UUID
	err = db.Pool.QueryRow(context.Background(), "SELECT id FROM roles WHERE name = 'USER'").Scan(&userRoleID)
	if err != nil {
		fmt.Println("error", err)
		return fmt.Errorf("ошибка при получении ID роли USER: %w", err)
	}
	fmt.Println("role id", userRoleID)
	_, err = db.Pool.Exec(context.Background(), "INSERT INTO user_roles (user_id, role_id) VALUES ($1, $2)", newUserUUID, userRoleID)
	if err != nil {
		return fmt.Errorf("ошибка при добавлении роли пользователя в базу данных: %w", err)
	}

	return nil
}

func AuthenticateUser(username, password string) (*models.AuthResponse, error) {
	var (
		UserID         uuid.UUID
		HashedPassword string
		Email          string
		Team           models.Team
	)

	var user *models.CurrentUser

	log.Printf("Попытка аутентификации пользователя: %s", username)

	err := db.Pool.QueryRow(context.Background(), "SELECT u.id, u.username, u.email, u.password, t.name, t.id, t.owner_id, t.invite_token FROM users u LEFT JOIN teams t on u.id = t.owner_id WHERE u.username = $1", username).Scan(&UserID, &username, &Email, &HashedPassword, &Team.Name, &Team.ID, &Team.OwnerID, &Team.InviteToken)
	if err != nil {
		return nil, fmt.Errorf("пользователь не найден: %w", err)
	}

	user, _ = GetUserByID(db.Pool, UserID)

	if err = bcrypt.CompareHashAndPassword([]byte(HashedPassword), []byte(password)); err != nil {
		return nil, fmt.Errorf("неверный пароль: %w", err)
	}

	accessTTL := time.Now().Add(60 * time.Minute)
	refreshTTL := time.Now().Add(24 * 7 * time.Hour)

	accessToken, refreshToken, err := GenerateTokens(UserID)

	if err != nil {
		return nil, fmt.Errorf("ошибка при создании токена: %w", err)
	}

	isAdmin, err := ValidateRole(UserID)
	if err != nil {
		fmt.Println("err", err)
		return nil, fmt.Errorf("ошибка при валидации роли пользователя: %w", err)
	}

	user.IsAdmin = isAdmin

	return &models.AuthResponse{
		Token:        accessToken,
		RefreshToken: refreshToken,
		AccessTTL:    accessTTL,
		RefreshTTL:   refreshTTL,
		User:         user,
	}, nil
}

func ValidateRole(userID uuid.UUID) (bool, error) {
	userRole, err := role.GetUserRole(userID)
	if err != nil {
		return false, fmt.Errorf("ошибка при получении роли пользователя: %w", err)
	}

	if userRole == "ADMIN" {
		return true, nil
	}

	return false, nil
}

func Refresh(refreshToken string) (*models.AuthResponse, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_KEY")), nil
	})

	if err != nil || !token.Valid || time.Unix(claims.ExpiresAt.Unix(), 0).Sub(time.Now()) <= 0 {
		return nil, fmt.Errorf("неверный refresh токен: %w", err)
	}

	accessTokenStr, refreshTokenStr, err := GenerateTokens(claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании токена: %w", err)
	}

	accessTTL := time.Now().Add(60 * time.Minute)
	refreshTTL := time.Now().Add(24 * 7 * time.Hour)

	return &models.AuthResponse{
		Token:        accessTokenStr,
		RefreshToken: refreshTokenStr,
		AccessTTL:    accessTTL,
		RefreshTTL:   refreshTTL,
	}, nil
}

func GenerateTokens(userID uuid.UUID) (accessTokenStr, refreshTokenStr string, err error) {
	jwtKey := []byte(os.Getenv("JWT_KEY"))

	accessExpirationTime := time.Now().Add(15 * time.Minute)
	accessClaims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExpirationTime),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenStr, err = accessToken.SignedString(jwtKey)
	if err != nil {
		return "", "", err
	}

	refreshExpirationTime := time.Now().Add(24 * time.Hour * 7)
	refreshClaims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExpirationTime),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenStr, err = refreshToken.SignedString(jwtKey)
	if err != nil {
		return "", "", err
	}

	return accessTokenStr, refreshTokenStr, nil
}

func GetUserByID(dbPool *pgxpool.Pool, userID uuid.UUID) (*models.CurrentUser, error) {
	var user models.CurrentUser
	var team models.Team
	var members []models.Member

	err := dbPool.QueryRow(context.Background(), `
    SELECT
      u.id,
      u.username,
      u.points,
      u.email,
      t.name,
      t.id AS team_id,
      t.owner_id,
      CASE WHEN u.id = t.owner_id THEN t.invite_token ELSE NULL END AS invite_token
    FROM
      users u
    LEFT JOIN
      team_members tm ON u.id = tm.user_id
    LEFT JOIN
      teams t ON tm.team_id = t.id OR u.id = t.owner_id
    WHERE
      u.id = $1`, userID).Scan(&user.ID, &user.Username, &user.Points, &user.Email, &team.Name, &team.ID, &team.OwnerID, &team.InviteToken)

	if err != nil {
		return nil, fmt.Errorf("error fetching user from database: %w", err)
	}

	if team.OwnerID != nil && *team.OwnerID != user.ID {
		team.InviteToken = nil
	}

	if team.ID != nil {
		user.Team = &models.Team{
			Name:        team.Name,
			ID:          team.ID,
			OwnerID:     team.OwnerID,
			InviteToken: team.InviteToken,
			TeamMembers: []models.Member{},
		}

		members, err = team2.GetTeamMembers(userID)
		if err != nil {
			return nil, fmt.Errorf("ошибка при получении участников команды: %w", err)
		}

		user.Team.TeamMembers = members
	} else {
		user.Team = nil
	}

	isAdmin, err := ValidateRole(user.ID)
	if err != nil {
		return nil, fmt.Errorf("ошибка при валидации роли пользователя: %w", err)
	}

	user.IsAdmin = isAdmin

	return &user, nil
}
