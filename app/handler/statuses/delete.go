package statuses

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"yatter-backend-go/app/handler/httperror"
	"yatter-backend-go/app/handler/request"
)

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
		if err == sql.ErrNoRows {
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
