package timelines

import (
	"encoding/json"
	"math"
	"net/http"
	"yatter-backend-go/app/handler/auth"
	"yatter-backend-go/app/handler/httperror"
	"yatter-backend-go/app/handler/request"
)

// Handle request for `GET /v1/timelines/public`
type ListRequest struct {
	MaxID   int64
	SinceID int64
	Limit   int64
}

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

	repo := h.app.Dao.Status()
	statuses, err := repo.ListAll(ctx, req.MaxID, req.SinceID, req.Limit)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	} else if statuses == nil {
		httperror.Error(w, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(statuses); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}

// Handle request for `GET /v1/timelines/home``
func (h *handler) Home(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	account := auth.AccountOf(r)

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

	repo := h.app.Dao.Status()
	statuses, err := repo.ListByID(ctx, account.ID, req.MaxID, req.SinceID, req.Limit)
	if err != nil {
		httperror.InternalServerError(w, err)
	} else if statuses == nil {
		httperror.Error(w, http.StatusNotFound)
		return
	}

	for i := range statuses {
		statuses[i].Account = *account
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(statuses); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}
