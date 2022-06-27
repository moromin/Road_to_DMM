package accounts

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"yatter-backend-go/app/handler/auth"
	"yatter-backend-go/app/handler/httperror"
	"yatter-backend-go/app/handler/media"
	"yatter-backend-go/app/handler/request"
)

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
