package accounts

import (
	"net/http"

	"yatter-backend-go/app/app"
	"yatter-backend-go/app/handler/auth"

	"github.com/go-chi/chi"
)

// Implementation of handler
type handler struct {
	app *app.App
}

// Create Handler for `/v1/accounts/`
func NewRouter(app *app.App) http.Handler {
	r := chi.NewRouter()

	h := &handler{app: app}

	r.Post("/", h.Create)
	r.Get("/{username}", h.Get)
	r.Get("/{username}/following", h.Following)

	r.Route("/", func(r chi.Router) {
		r.Use(auth.BasicAuth(h.app))
		r.Post("/{username}/follow", h.Follow)
	})

	return r
}
