package statuses

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"yatter-backend-go/app/handler/auth"
	"yatter-backend-go/app/handler/httperror"
	"yatter-backend-go/app/handler/request"
)

// Handle request for `POST /v1/statuses`
// Request body
type CreateRequest struct {
	Status   string  `json:"status"`
	MediaIDs []int64 `json:"media_ids"`
}

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httperror.BadRequest(w, err)
		return
	}

	account := auth.AccountOf(r)

	statusRepo := h.app.Dao.Status()
	id, err := statusRepo.Create(ctx, account.ID, req.Status, req.MediaIDs)
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

// Handle request for `GET /v1/statuses/{id}`
func (h *handler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := request.IDOf(r)
	if err != nil {
		httperror.BadRequest(w, err)
		return
	}

	statusRepo := h.app.Dao.Status()
	status, err := statusRepo.FindByID(ctx, id)
	if err != nil {
		httperror.BadRequest(w, err)
		return
	} else if status == nil {
		httperror.Error(w, http.StatusNotFound)
		return
	}

	accountRepo := h.app.Dao.Account()
	account, err := accountRepo.FindByID(ctx, status.AccountID)
	if err != nil {
		httperror.BadRequest(w, err)
		return
	} else if status == nil {
		httperror.Error(w, http.StatusNotFound)
		return
	}
	status.Account = *account

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(status); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}

// Handle request for `DELETE /v1/statuses/{id}`
func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := request.IDOf(r)
	if err != nil {
		httperror.BadRequest(w, err)
		return
	}

	statusRepo := h.app.Dao.Status()
	err = statusRepo.DeleteByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			httperror.Error(w, http.StatusNotFound)
		} else {
			httperror.BadRequest(w, err)
		}
		return
	}

	status := struct{}{}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(status); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}
