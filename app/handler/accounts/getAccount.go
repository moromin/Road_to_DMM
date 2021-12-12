package accounts

import (
	"encoding/json"
	"net/http"
	"yatter-backend-go/app/handler/httperror"

	"github.com/go-chi/chi"
)

// Request body for `GET /v1/accounts/{username}`
type getAccountRequest struct {
	Username string
}

// Handle request for `GET /v1/accounts/{username}`
func (h *handler) GetAccount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req getAccountRequest
	req.Username = chi.URLParam(r, "username")

	repo := h.app.Dao.Account()
	account, err := repo.FindByUsername(ctx, req.Username)
	if account == nil {
		httperror.Error(w, http.StatusNotFound)
		return
	} else if err != nil {
		httperror.BadRequest(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(account); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}
