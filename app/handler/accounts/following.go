package accounts

import (
	"encoding/json"
	"net/http"
	"yatter-backend-go/app/handler/httperror"
	"yatter-backend-go/app/handler/request"

	"github.com/go-chi/chi"
)

func (h *handler) Following(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	username := chi.URLParam(r, "username")

	const limit = "limit"
	options := []request.Option{{limit, 40, 0, 80}}
	params, err := request.GetOptionParams(r, options)
	if err != nil {
		httperror.BadRequest(w, err)
		return
	}

	repo := h.app.Dao.Account()
	follower, err := repo.FindByUsername(ctx, username)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}
	if follower == nil {
		httperror.Error(w, http.StatusNotFound)
		return
	}

	accounts, err := repo.FindFollowing(ctx, follower.ID, params[limit])
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
