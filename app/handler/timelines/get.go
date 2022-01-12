package timelines

import (
	"encoding/json"
	"net/http"
	"strconv"
	"yatter-backend-go/app/handler/httperror"
)

type ListRequest struct {
	MaxID   int64
	SinceID int64
	Limit   int64
}

// Handle request for `GET /v1/timelines/*`
func (h *handler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	strMaxID := r.URL.Query().Get("max_id")
	strSinceID := r.URL.Query().Get("since_id")
	strLimit := r.URL.Query().Get("limit")

	req := ListRequest{
		MaxID:   0,
		SinceID: 0,
		Limit:   40,
	}

	if strMaxID != "" {
		maxID, err := strconv.ParseInt(strMaxID, 10, 64)
		if err != nil {
			httperror.BadRequest(w, err)
			return
		} else if maxID <= 0 {
			httperror.Error(w, http.StatusBadRequest)
			return
		}
		req.MaxID = maxID
	}

	if strSinceID != "" {
		sinceID, err := strconv.ParseInt(strSinceID, 10, 64)
		if err != nil {
			httperror.BadRequest(w, err)
			return
		} else if sinceID <= 0 {
			httperror.Error(w, http.StatusBadRequest)
			return
		}
		req.SinceID = sinceID
	}

	if strLimit != "" {
		limit, err := strconv.ParseInt(strLimit, 10, 64)
		if err != nil {
			httperror.BadRequest(w, err)
			return
		} else if limit < 0 || 80 < limit {
			httperror.Error(w, http.StatusBadRequest)
			return
		}
		req.Limit = limit
	}

	statusRepo := h.app.Dao.Status()
	statuses, err := statusRepo.List(ctx, req.MaxID, req.SinceID, req.Limit)
	if err != nil {
		httperror.BadRequest(w, err)
		return
	} else if statuses == nil {
		httperror.Error(w, http.StatusNotFound)
		return
	}

	accountRepo := h.app.Dao.Account()
	for i, status := range statuses {
		account, err := accountRepo.FindByID(ctx, status.AccountID)
		if err != nil {
			httperror.BadRequest(w, err)
			return
		}
		statuses[i].Account = *account
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(statuses); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}
