package media

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"unicode/utf8"
	"yatter-backend-go/app/handler/httperror"
	"yatter-backend-go/app/handler/request"
)

const maxDescriptionLength = 420

func (h *handler) Upload(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	repo := h.app.Dao.Attachment()

	description := r.FormValue("description")
	if utf8.RuneCountInString(description) > maxDescriptionLength {
		httperror.BadRequest(w, fmt.Errorf("description is too long, please less than or equal %d", maxDescriptionLength))
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		httperror.BadRequest(w, errors.New("invalid file name in request body"))
		return
	}
	defer file.Close()

	filetype, err := request.DetectAttachmentType(file)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	attachment, err := repo.UploadFile(ctx, file, fileHeader.Filename, filetype, description)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(attachment); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}
