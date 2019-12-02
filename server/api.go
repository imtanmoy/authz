package server

import (
	"github.com/imtanmoy/authz/group"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"

	"github.com/imtanmoy/authz/db"
	"github.com/imtanmoy/authz/organizations"
	"github.com/imtanmoy/authz/users"
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
	r.Mount("/users", userRouter())
	r.Mount("/{oid}/groups", groupRouter())

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

func userRouter() http.Handler {
	r := chi.NewRouter()
	userHandler := users.NewUserHandler(db.DB)

	r.Group(func(r chi.Router) {
		r.Get("/", userHandler.List)
		r.Post("/", userHandler.Create)
		r.Get("/{id}", userHandler.Get)
		r.Put("/{id}", userHandler.Update)
		r.Delete("/{id}", userHandler.Delete)

		r.Group(func(r chi.Router) {
			r.Use(userHandler.UserCtx)
			r.Get("/{id}/groups", userHandler.GetGroups)
			r.Get("/{id}/permissions", userHandler.GetPermissions)
		})
	})

	return r
}

func groupRouter() http.Handler {
	r := chi.NewRouter()
	groupHandler := group.NewGroupHandler(db.DB)
	r.Use(groupHandler.OrganizationCtx)

	r.Group(func(r chi.Router) {
		r.Get("/", groupHandler.List)
		r.Post("/", groupHandler.Create)
		r.Group(func(r chi.Router) {
			r.Use(groupHandler.GroupCtx)
			r.Get("/{id}", groupHandler.Get)
			r.Put("/{id}", groupHandler.Update)
			r.Delete("/{id}", groupHandler.Delete)
		})
	})

	return r
}