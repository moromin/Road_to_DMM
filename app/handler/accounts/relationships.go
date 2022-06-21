package accounts

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"yatter-backend-go/app/handler/auth"
	"yatter-backend-go/app/handler/httperror"
)

type Relationship struct {
	ID         int64 `json:"id"`
	Following  bool  `json:"following"`
	FollowedBy bool  `json:"followed_by"`
}

func (h *handler) Relationships(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	repo := h.app.Dao.Account()

	user := auth.AccountOf(r)

	usernames := strings.Split(r.URL.Query().Get("username"), ",")

	accounts := make(map[string]int64)
	for _, username := range usernames {
		if username == user.Username {
			httperror.BadRequest(w, errors.New("specifying yourself is forbidden"))
			return
		}
		account, err := repo.FindByUsername(ctx, username)
		if err != nil {
			httperror.InternalServerError(w, err)
			return
		} else if account == nil {
			httperror.BadRequest(w, errors.New("account you want to know relationship is not found"))
			return
		}
		accounts[username] = account.ID
	}

	relationships := make([]Relationship, 0)
	for _, targetID := range accounts {
		following, followedBy, err := repo.FindRelationship(ctx, user.ID, targetID)
		if err != nil {
			httperror.InternalServerError(w, err)
			return
		}
		relationship := Relationship{
			ID:         targetID,
			Following:  following,
			FollowedBy: followedBy,
		}
		relationships = append(relationships, relationship)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(relationships); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}
