package statuses

import (
	"encoding/json"
	"net/http"
	"yatter-backend-go/app/handler/httperror"
	"yatter-backend-go/app/handler/request"
)

// Handler request for `GET /v1/statuses/{id}`
func (h *handler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := request.IDOf(r)
	if err != nil {
		httperror.BadRequest(w, err)
		return
	}

	statusRepo := h.app.Dao.Status()
	status, err := statusRepo.FindByID(ctx, id)
	if status == nil {
		httperror.Error(w, http.StatusNotFound)
		return
	} else if err != nil {
		httperror.BadRequest(w, err)
		return
	}

	accountRepo := h.app.Dao.Account()
	account, err := accountRepo.FindByID(ctx, status.AccountID)
	if err != nil {
		httperror.BadRequest(w, err)
		return
	}
	status.Account = *account

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(status); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}
