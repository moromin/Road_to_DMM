package statuses

import (
	"encoding/json"
	"errors"
	"net/http"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/handler/auth"
	"yatter-backend-go/app/handler/httperror"
)

// Request body for `POST /v1/statuses`
type AddRequest struct {
	Content string
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
	if account == nil {
		httperror.InternalServerError(w, errors.New("failed to read ID"))
		return
	}

	status := &object.Status{
		Content: req.Content,
		Account: *account,
	}

	statusRepo := h.app.Dao.Status()
	id, err := statusRepo.Create(ctx, status.Account.ID, status.Content)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}
	status.ID = id

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(status); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}
