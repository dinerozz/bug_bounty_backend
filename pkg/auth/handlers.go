package auth

import (
	"encoding/json"
	db "github.com/dinerozz/bug_bounty_backend/config"
	"net/http"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Некорректный запрос", http.StatusBadRequest)
		return
	}

	if err := RegisterUser(req.Username, req.Email, req.Password); err != nil {
		http.Error(w, "Ошибка при регистрации: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode("Пользователь " + req.Username + " успешно зарегистрирован")
}

func AuthenticateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Некорректный запрос", http.StatusBadRequest)
		return
	}

	token, err := AuthenticateUser(req.Username, req.Password)
	if err != nil {
		http.Error(w, "Ошибка при аутентификации: "+err.Error(), http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}

func CurrentUserHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := GetUserFromJWT(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	user, err := db.GetUserByID(db.Pool, userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}
