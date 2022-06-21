package timelines

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

	r.Get("/public", h.Public)

	r.Route("/", func(r chi.Router) {
		r.Use(auth.BasicAuth(h.app))
		r.Get("/home", h.Home)
	})

	return r
}
