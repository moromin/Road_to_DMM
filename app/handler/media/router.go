package media

import (
	"log"
	"net/http"
	"os"
	"yatter-backend-go/app/app"

	"github.com/go-chi/chi"
)

type handler struct {
	app *app.App
}

const FilesDir = "./files"

func NewRouter(app *app.App) http.Handler {
	r := chi.NewRouter()

	h := &handler{app: app}

	r.Post("/", h.Upload)

	workDir, _ := os.Getwd()
	filePath := workDir + "/files"
	log.Println(filePath)

	return r
}
