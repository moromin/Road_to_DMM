package media

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"unicode/utf8"
	"yatter-backend-go/app/handler/httperror"
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

	filetype, err := detectContentType(file)
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

func detectContentType(file io.ReadSeeker) (string, error) {
	buff := make([]byte, 512)
	if _, err := file.Read(buff); err != nil {
		return "", err
	}

	filetype := http.DetectContentType(buff)

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return "", err
	}

	return classifyFiletype(filetype), nil
}

var imageType = regexp.MustCompile("image/*")
var videoType = regexp.MustCompile("video/*")

func classifyFiletype(filetype string) string {
	if filetype == "image/gif" {
		return "gifv"
	} else if imageType.MatchString(filetype) {
		return "image"
	} else if videoType.MatchString(filetype) {
		return "video"
	}
	return "unknown"
}
