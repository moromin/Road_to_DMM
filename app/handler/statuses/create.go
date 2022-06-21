package statuses

import (
	"encoding/json"
	"net/http"
	"yatter-backend-go/app/handler/auth"
	"yatter-backend-go/app/handler/httperror"
)

// Request body for `POST /v1/statuses`
type AddRequest struct {
	Status string
}

// Handle request for `POST /v1/statuses`
func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req AddRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httperror.BadRequest(w, err)
		return
	}

	account := auth.AccountOf(r)

	statusRepo := h.app.Dao.Status()
	id, err := statusRepo.Create(ctx, account.ID, req.Status)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}
	status, err := statusRepo.FindByID(ctx, id)
	if err != nil || status == nil {
		httperror.InternalServerError(w, err)
		return
	}
	status.Account = *account

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(status); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}
