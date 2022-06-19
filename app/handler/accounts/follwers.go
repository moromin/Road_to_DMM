package accounts

import (
	"encoding/json"
	"math"
	"net/http"
	"yatter-backend-go/app/handler/httperror"
	"yatter-backend-go/app/handler/request"

	"github.com/go-chi/chi"
)

func (h *handler) Followers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	username := chi.URLParam(r, "username")

	const (
		maxID   = "max_id"
		sinceID = "since_id"
		limit   = "limit"
	)

	options := []request.Option{
		{maxID, 0, 1, math.MaxInt64},
		{sinceID, 0, 1, math.MaxInt64},
		{limit, 40, 0, 80},
	}
	params, err := request.GetOptionParams(r, options)
	if err != nil {
		httperror.BadRequest(w, err)
		return
	}

	repo := h.app.Dao.Account()
	followee, err := repo.FindByUsername(ctx, username)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	} else if followee == nil {
		httperror.Error(w, http.StatusNotFound)
		return
	}

	accounts, err := repo.FindFollowers(ctx, followee.ID, params[maxID], params[sinceID], params[limit])
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(accounts); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}
