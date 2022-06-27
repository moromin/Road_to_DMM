package media

import (
	"net/http"
	"os"
	"path/filepath"
	"yatter-backend-go/app/app"

	"github.com/go-chi/chi"
)

type handler struct {
	app *app.App
}

const FilesDir = "files"

func NewRouter(app *app.App) http.Handler {
	r := chi.NewRouter()

	h := &handler{app: app}

	r.Post("/", h.Upload)

	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, FilesDir))
	FileServer(r, "/files", filesDir)

	return r
}
