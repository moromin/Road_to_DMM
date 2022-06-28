package accounts

import (
	"net/http"

	"yatter-backend-go/app/app"
	"yatter-backend-go/app/handler/auth"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
)

// Implementation of handler
type handler struct {
	app       *app.App
	validator *validator.Validate
}

// Create Handler for `/v1/accounts/`
func NewRouter(app *app.App, validator *validator.Validate) http.Handler {
	r := chi.NewRouter()

	h := &handler{
		app:       app,
		validator: validator,
	}

	r.Route("/", func(r chi.Router) {
		r.Use(auth.BasicAuth(h.app))
		r.Post("/{username}/follow", h.Follow)
		r.Post("/{username}/unfollow", h.Unfollow)
	})
	r.With(auth.BasicAuth(h.app)).Get("/relationships", h.Relationships)
	r.With(auth.BasicAuth(h.app)).Post("/update_credentials", h.UpdateCredentials)

	r.Post("/", h.Create)
	r.Get("/{username}", h.Get)
	r.Get("/{username}/following", h.Following)
	r.Get("/{username}/followers", h.Followers)

	return r
}
