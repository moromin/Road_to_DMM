package timelines

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/handler/httperror"
)

type ListRequest struct {
	MaxID   int64
	SinceID int64
	Limit   int64
}

type Option struct {
	name         string
	defaultValue int64
	min          int64
	max          int64
}

// Handle request for `GET /v1/timelines/public`
func (h *handler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	options := []Option{
		{"max_id", 0, 1, math.MaxInt64},
		{"since_id", 0, 1, math.MaxInt64},
		{"limit", 40, 0, 80},
	}
	params, err := getOptionParams(r, options)
	if err != nil {
		httperror.BadRequest(w, err)
		return
	}
	req := ListRequest{
		MaxID:   params["max_id"],
		SinceID: params["since_id"],
		Limit:   params["limit"],
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

func getOptionParams(r *http.Request, options []Option) (map[string]int64, error) {
	params := make(map[string]int64)
	var err error

	for _, op := range options {
		params[op.name], err = paramOf(r.URL.Query().Get(op.name), op.defaultValue, op.min, op.max)
		if err != nil {
			return nil, err
		}
	}

	return params, nil
}

func paramOf(strParam string, defaultValue, min, max int64) (int64, error) {
	if strParam == "" {
		return defaultValue, nil
	}

	param, err := strconv.ParseInt(strParam, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("%q is invalid format for option", strParam)
	}
	if param < min || max < param {
		return 0, fmt.Errorf("%d is over valid range [%d, %d]", param, min, max)
	}

	return param, nil
}
