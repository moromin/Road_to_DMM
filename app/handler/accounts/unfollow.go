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

func (h *handler) Unfollow(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	follower := auth.AccountOf(r)
	username := chi.URLParam(r, "username")
	if username == follower.Username {
		httperror.BadRequest(w, errors.New("unfollowing yourself is forbidden"))
		return
	}

	repo := h.app.Dao.Account()
	followee, err := repo.FindByUsername(ctx, username)
	if err != nil {
		httperror.BadRequest(w, err)
		return
	} else if followee == nil {
		httperror.BadRequest(w, fmt.Errorf("%s you want to unfollow is not found", username))
		return
	}

	id, followedBy, err := repo.Unfollow(ctx, follower.ID, followee.ID)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	res := &FollowResponse{
		ID:         id,
		Following:  false,
		FollowedBy: followedBy,
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}
