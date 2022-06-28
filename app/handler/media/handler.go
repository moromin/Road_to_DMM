package media

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"unicode/utf8"
	"yatter-backend-go/app/handler/httperror"
	"yatter-backend-go/app/handler/request"

	"github.com/go-chi/chi"
)

const maxDescriptionLength = 420

// Handle request for `POST /v1/media`
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

	attachment, err := repo.UploadFile(ctx, file, FilesDir, fileHeader.Filename, filetype, description)
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

// Handle request for `GET /v1/media/files/{id}`
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		log.Println("FileServer does not permit URL parameters.")
		return
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}

	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	}))
}
