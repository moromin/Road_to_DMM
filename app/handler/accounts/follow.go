package accounts

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"yatter-backend-go/app/handler/auth"
	"yatter-backend-go/app/handler/httperror"

	"github.com/go-chi/chi"
)

type Response struct {
	ID         int64 `json:"id"`
	Following  bool  `json:"following"`
	FollowedBy bool  `json:"followed_by"`
}

func (h *handler) Follow(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	follower := auth.AccountOf(r)
	username := chi.URLParam(r, "username")
	if username == follower.Username {
		httperror.BadRequest(w, errors.New("following yourself is forbidden"))
		return
	}

	repo := h.app.Dao.Account()
	followee, err := repo.FindByUsername(ctx, username)
	if err != nil {
		httperror.BadRequest(w, err)
		return
	}
	if followee == nil {
		httperror.BadRequest(w, fmt.Errorf("%s you want to follow is not found", username))
		return
	}

	id, followedBy, err := repo.Follow(ctx, follower.ID, followee.ID)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	res := &Response{
		ID:         id,
		Following:  true,
		FollowedBy: followedBy,
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}
