package statuses

import (
	"encoding/json"
	"fmt"
	"net/http"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/handler/auth"
	"yatter-backend-go/app/handler/httperror"
)

// Request body for `POST /v1/statuses`
type AddRequest struct {
	Content string
}

type errFailedToReadContext struct{}

func (e errFailedToReadContext) Error() string {
	return fmt.Sprint("failed to read context value")
}

// Handle request for `POST /v1/statuses`
// TODO: empty is OK?
func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req AddRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httperror.BadRequest(w, err)
		return
	}

	account := auth.AccountOf(r)
	if account == nil {
		httperror.InternalServerError(w, errFailedToReadContext{})
		return
	}

	status := &object.Status{
		Content: req.Content,
		Account: *account,
	}

	statusRepo := h.app.Dao.Status()
	if err := statusRepo.Create(ctx, status.Content, status.Account.ID); err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(status); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}
