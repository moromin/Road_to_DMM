package statuses

import (
	"net/http"
	"yatter-backend-go/app/app"
	"yatter-backend-go/app/handler/auth"

	"github.com/go-chi/chi"
)

type handler struct {
	app *app.App
}

func NewRouter(app *app.App) http.Handler {
	r := chi.NewRouter()

	h := &handler{app: app}

	r.Route("/", func(r chi.Router) {
		r.Use(auth.Middleware(h.app))
		r.Post("/", h.Create)
	})

	return r
}
