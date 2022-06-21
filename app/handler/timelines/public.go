package timelines

import (
	"encoding/json"
	"math"
	"net/http"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/handler/httperror"
	"yatter-backend-go/app/handler/request"
)

type ListRequest struct {
	MaxID   int64
	SinceID int64
	Limit   int64
}

// Handle request for `GET /v1/timelines/public`
func (h *handler) Public(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	const (
		maxID   = "max_id"
		sinceID = "since_id"
		limit   = "limit"
	)

	options := []request.Option{
		{maxID, 0, 1, math.MaxInt64},
		{sinceID, 0, 1, math.MaxInt64},
		{limit, 40, 0, 80},
	}
	params, err := request.GetOptionParams(r, options)
	if err != nil {
		httperror.BadRequest(w, err)
		return
	}
	req := ListRequest{
		MaxID:   params[maxID],
		SinceID: params[sinceID],
		Limit:   params[limit],
	}

	statusRepo := h.app.Dao.Status()
	statuses, err := statusRepo.ListAll(ctx, req.MaxID, req.SinceID, req.Limit)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	} else if statuses == nil {
		httperror.Error(w, http.StatusNotFound)
		return
	}

	// TODO: delete this to solve N + 1 problem
	accounts := make(map[int]*object.Account)
	accountRepo := h.app.Dao.Account()
	for i, status := range statuses {
		account, ok := accounts[status.AccountID]
		if !ok {
			account, err = accountRepo.FindByID(ctx, status.AccountID)
			if err != nil {
				httperror.BadRequest(w, err)
				return
			}
			accounts[status.AccountID] = account
		}
		statuses[i].Account = *account
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(statuses); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}
