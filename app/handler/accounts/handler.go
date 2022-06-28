package accounts

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strings"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/handler/auth"
	"yatter-backend-go/app/handler/httperror"
	"yatter-backend-go/app/handler/media"
	"yatter-backend-go/app/handler/request"
	"yatter-backend-go/app/handler/validate"

	"github.com/go-chi/chi"
)

// Handle request for `POST /v1/accounts`
// Request body
type CreateRequest struct {
	Username string `validate:"required"`
	Password string `validate:"required"`
}

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httperror.BadRequest(w, err)
		return
	}

	if err := validate.Validate(h.validator, req); err != nil {
		httperror.BadRequest(w, err)
		return
	}

	account := new(object.Account)
	account.Username = req.Username
	if err := account.SetPassword(req.Password); err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	repo := h.app.Dao.Account()
	if err := repo.CreateAccount(ctx, account.Username, account.PasswordHash); err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	res, err := repo.FindByUsername(ctx, account.Username)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}
	account.CreateAt = res.CreateAt

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(account); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}

// Handle request for `GET /v1/accounts/{username}`
func (h *handler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	username := chi.URLParam(r, "username")

	repo := h.app.Dao.Account()
	account, err := repo.FindByUsername(ctx, username)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}
	if account == nil {
		httperror.Error(w, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(account); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}

// Relationship represents relationship between two accounts
type Relationship struct {
	ID         int64 `json:"id"`
	Following  bool  `json:"following"`
	FollowedBy bool  `json:"followed_by"`
}

// Handle request for `POST /v1/accounts/{username}/follow`
func (h *handler) Follow(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	follower := auth.AccountOf(r)
	username := chi.URLParam(r, "username")
	if username == follower.Username {
		httperror.BadRequest(w, errors.New("following yourself is forbidden"))
		return
	}

	repo := h.app.Dao.Account()
	followee, err := repo.FindByUsername(ctx, username)
	if err != nil {
		httperror.BadRequest(w, err)
		return
	} else if followee == nil {
		httperror.BadRequest(w, fmt.Errorf("%s you want to follow is not found", username))
		return
	}

	id, followedBy, err := repo.Follow(ctx, follower.ID, followee.ID)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	res := &Relationship{
		ID:         id,
		Following:  true,
		FollowedBy: followedBy,
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}

// Handle request for `POST /v1/accounts/{username}/unfollow`
func (h *handler) Unfollow(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	follower := auth.AccountOf(r)
	username := chi.URLParam(r, "username")
	if username == follower.Username {
		httperror.BadRequest(w, errors.New("unfollowing yourself is forbidden"))
		return
	}

	repo := h.app.Dao.Account()
	followee, err := repo.FindByUsername(ctx, username)
	if err != nil {
		httperror.BadRequest(w, err)
		return
	} else if followee == nil {
		httperror.BadRequest(w, fmt.Errorf("%s you want to unfollow is not found", username))
		return
	}

	id, followedBy, err := repo.Unfollow(ctx, follower.ID, followee.ID)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	res := &Relationship{
		ID:         id,
		Following:  false,
		FollowedBy: followedBy,
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}

// Handle request for `GET /v1/accounts/{username}/following`
func (h *handler) Following(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	username := chi.URLParam(r, "username")

	const limit = "limit"
	options := []request.Option{{limit, 40, 0, 80}}
	params, err := request.GetOptionParams(r, options)
	if err != nil {
		httperror.BadRequest(w, err)
		return
	}

	repo := h.app.Dao.Account()
	follower, err := repo.FindByUsername(ctx, username)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}
	if follower == nil {
		httperror.Error(w, http.StatusNotFound)
		return
	}

	accounts, err := repo.FindFollowing(ctx, follower.ID, params[limit])
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(accounts); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}

// Handle request for `GET /v1/accounts/{username}/followers`
func (h *handler) Followers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	username := chi.URLParam(r, "username")

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

	repo := h.app.Dao.Account()
	followee, err := repo.FindByUsername(ctx, username)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	} else if followee == nil {
		httperror.Error(w, http.StatusNotFound)
		return
	}

	accounts, err := repo.FindFollowers(ctx, followee.ID, params[maxID], params[sinceID], params[limit])
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(accounts); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}

// Relationships handles request for `GET /v1/accounts/relationships`
func (h *handler) Relationships(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	repo := h.app.Dao.Account()

	user := auth.AccountOf(r)

	usernames := strings.Split(r.URL.Query().Get("username"), ",")

	accounts := make(map[string]int64)
	for _, username := range usernames {
		if username == user.Username {
			httperror.BadRequest(w, errors.New("specifying yourself is forbidden"))
			return
		}
		account, err := repo.FindByUsername(ctx, username)
		if err != nil {
			httperror.InternalServerError(w, err)
			return
		} else if account == nil {
			httperror.BadRequest(w, errors.New("account you want to know relationship is not found"))
			return
		}
		accounts[username] = account.ID
	}

	relationships := make([]Relationship, 0)
	for _, targetID := range accounts {
		following, followedBy, err := repo.FindRelationship(ctx, user.ID, targetID)
		if err != nil {
			httperror.InternalServerError(w, err)
			return
		}
		relationship := Relationship{
			ID:         targetID,
			Following:  following,
			FollowedBy: followedBy,
		}
		relationships = append(relationships, relationship)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(relationships); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}

// UpdateCredentials handles request for `POST /v1/accounts/update_credentials`
func (h *handler) UpdateCredentials(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	repo := h.app.Dao.Account()

	account := auth.AccountOf(r)

	displayName := r.FormValue("display_name")
	note := r.FormValue("note")

	avatar, code, err := h.uploadFormFile(r, ctx, "avatar")
	if err != nil {
		http.Error(w, err.Error(), code)
		return
	}

	header, code, err := h.uploadFormFile(r, ctx, "header")
	if err != nil {
		http.Error(w, err.Error(), code)
		return
	}

	if err := repo.UpdateCredentials(ctx, account.ID, displayName, note, avatar, header); err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	account, err = repo.FindByID(ctx, account.ID)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(account); err != nil {
		httperror.InternalServerError(w, err)
		return
	}

}

func (h *handler) uploadFormFile(r *http.Request, ctx context.Context, name string) (string, int, error) {
	file, header, err := r.FormFile(name)
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			return "", http.StatusOK, nil
		}
		return "", http.StatusBadRequest, err
	}
	defer file.Close()

	filetype, err := request.DetectAttachmentType(file)
	if err != nil {
		return "", http.StatusInternalServerError, err
	} else if filetype != request.Image {
		return "", http.StatusBadRequest, errors.New("invalid file type, please image (jpeg, png, etc.)")
	}

	repo := h.app.Dao.Attachment()
	attachment, err := repo.UploadFile(ctx, file, media.FilesDir, header.Filename, filetype, "")
	if err != nil {
		return "", http.StatusInternalServerError, err
	}

	return attachment.URL, http.StatusOK, nil
}
