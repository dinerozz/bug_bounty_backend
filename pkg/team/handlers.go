package team

import (
	"encoding/json"
	"github.com/dinerozz/bug_bounty_backend/pkg/models"
	"github.com/google/uuid"
	"net/http"
)

func CreateTeamHandler(w http.ResponseWriter, r *http.Request) {
	var newTeam models.Team
	if err := json.NewDecoder(r.Body).Decode(&newTeam); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userIdFromContext := r.Context().Value("userID")
	userID, ok := userIdFromContext.(uuid.UUID)

	if !ok {
		http.Error(w, "Could not find user ID", http.StatusUnauthorized)
		return
	}

	newTeam.OwnerID = userID

	err := CreateTeam(&newTeam)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTeam)
}
