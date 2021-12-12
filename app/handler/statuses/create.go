package statuses

import (
	"encoding/json"
	"net/http"
	"strings"
	"yatter-backend-go/app/domain/object"
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

	a := r.Header.Get("Authentication")
	pair := strings.SplitN(a, " ", 2)
	username := pair[1]

	accountRepo := h.app.Dao.Account()
	account, _ := accountRepo.FindByUsername(ctx, username)

	status := new(object.Status)
	status.Content = req.Content
	status.Account = *account

	statusRepo := h.app.Dao.Status()
	if err := statusRepo.CreateStatus(ctx, status.Content, status.Account.ID); err != nil {
		httperror.BadRequest(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(status); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}
