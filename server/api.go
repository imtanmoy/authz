package server

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	
	"github.com/imtanmoy/authz/db"
	"github.com/imtanmoy/authz/organizations"
)

// New configures application resources and routes.
func New() (*chi.Mux, error) {

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.DefaultCompress)
	r.Use(middleware.Timeout(15 * time.Second))

	//r.Use(logging.NewStructuredLogger(logger))
	r.Use(render.SetContentType(render.ContentTypeJSON))

	//r.Use(corsConfig().Handler)

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("pong"))
	})

	//routes.Routes(r)
	r.Mount("/organizations", organizationRouter())

	return r, nil
}

func organizationRouter() http.Handler {
	r := chi.NewRouter()
	organizationHandler := organizations.NewOrganizationHandler(db.DB)

	r.Group(func(r chi.Router) {
		r.Get("/", organizationHandler.List)
		r.Post("/", organizationHandler.Create)
		r.Get("/{id}", organizationHandler.Get)
		r.Put("/{id}", organizationHandler.Update)
		r.Delete("/{id}", organizationHandler.Delete)
	})

	return r
}
