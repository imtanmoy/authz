package users

import (
	"context"
	"github.com/imtanmoy/authz/authorizer"
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-pg/pg/v9"
	"github.com/imtanmoy/authz/organizations"
	"github.com/imtanmoy/authz/utils/httputil"
	param "github.com/oceanicdev/chi-param"

	"github.com/imtanmoy/authz/models"
)

type Handler interface {
	UserCtx(next http.Handler) http.Handler
	List(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)

	GetGroups(w http.ResponseWriter, r *http.Request)
	GetPermissions(w http.ResponseWriter, r *http.Request)
}

type userHandler struct {
	db                  *pg.DB
	service             Service
	organizationService organizations.Service
}

var _ Handler = (*userHandler)(nil)

func NewUserHandler(db *pg.DB, authorizationService authorizer.Service) Handler {
	return &userHandler{
		db:                  db,
		service:             NewUserService(db, authorizationService),
		organizationService: organizations.NewOrganizationService(db),
	}
}

func (u *userHandler) UserCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := param.Int32(r, "id")
		if err != nil {
			_ = render.Render(w, r, httputil.NewAPIError(400, "Invalid request parameter", err))
			return
		}
		user, err := u.service.Find(id)
		if err != nil {
			_ = render.Render(w, r, httputil.NewAPIError(404, "user not found", err))
			return
		}
		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (u *userHandler) List(w http.ResponseWriter, r *http.Request) {
	users, err := u.service.List()
	if err != nil {
		_ = render.Render(w, r, httputil.NewAPIError(err))
		return
	}
	if err := render.RenderList(w, r, NewUserListResponse(users)); err != nil {
		_ = render.Render(w, r, httputil.NewAPIError(err))
		return
	}
}

func (u *userHandler) Create(w http.ResponseWriter, r *http.Request) {

	data := &UserPayload{}
	if err := render.Bind(r, data); err != nil {
		_ = render.Render(w, r, httputil.NewAPIError(err))
		return
	}

	validationErrors := data.validate()

	if len(validationErrors) > 0 {
		_ = render.Render(w, r, httputil.NewAPIError(400, "Invalid request", validationErrors))
		return
	}
	exist := u.service.Exists(data.ID)
	orgExist := u.organizationService.Exists(data.OrganizationID)
	existErr := make(map[string][]string)
	if exist {
		existErr = map[string][]string{
			"id": {"User with same id already exits"},
		}
	}
	if !orgExist {
		existErr = map[string][]string{
			"organization_id": {"organization does not exist"},
		}
	}
	if len(existErr) > 0 {
		_ = render.Render(w, r, httputil.NewAPIError(400, "Invalid request", existErr))
		return
	}

	organization, err := u.organizationService.Find(data.OrganizationID)
	if err != nil {
		_ = render.Render(w, r, httputil.NewAPIError(err))
		return
	}

	var user models.User
	user.ID = data.ID
	user.Email = data.Email
	user.Organization = organization
	user.OrganizationID = organization.ID

	newUser, err := u.service.Create(&user)

	if err != nil {
		_ = render.Render(w, r, httputil.NewAPIError(err))
		return
	}

	render.Status(r, http.StatusCreated)
	_ = render.Render(w, r, NewUserResponse(newUser))
	return
}

func (u *userHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := param.Int32(r, "id")
	if err != nil {
		_ = render.Render(w, r, httputil.NewAPIError(400, "Invalid request parameter", err))
		return
	}
	user, err := u.service.Find(id)
	if err != nil {
		_ = render.Render(w, r, httputil.NewAPIError(404, "user not found", err))
		return
	}
	if err := render.Render(w, r, NewUserResponse(user)); err != nil {
		_ = render.Render(w, r, httputil.NewAPIError(err))
		return
	}
}

func (u *userHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := param.Int32(r, "id")
	if err != nil {
		_ = render.Render(w, r, httputil.NewAPIError(400, "Invalid request parameter", err))
		return
	}

	data := &UserPayload{}
	if err := render.Bind(r, data); err != nil {
		_ = render.Render(w, r, httputil.NewAPIError(err))
		return
	}

	validationErrors := data.validate()

	if len(validationErrors) > 0 {
		_ = render.Render(w, r, httputil.NewAPIError(400, "Invalid request", validationErrors))
		return
	}

	user, err := u.service.Find(id)
	if err != nil {
		_ = render.Render(w, r, httputil.NewAPIError(err))
		return
	}

	user.Email = data.Email
	user, err = u.service.Update(user)
	if err != nil {
		_ = render.Render(w, r, httputil.NewAPIError(err))
		return
	}
	if err := render.Render(w, r, NewUserResponse(user)); err != nil {
		_ = render.Render(w, r, httputil.NewAPIError(err))
		return
	}
}

func (u *userHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := param.Int32(r, "id")
	if err != nil {
		_ = render.Render(w, r, httputil.NewAPIError(400, "Invalid request parameter", err))
		return
	}
	user, err := u.service.Find(id)
	if err != nil {
		_ = render.Render(w, r, httputil.NewAPIError(err))
		return
	}
	err = u.service.Delete(user)
	if err != nil {
		_ = render.Render(w, r, httputil.NewAPIError(err))
		return
	}
	render.NoContent(w, r)
}

func (u *userHandler) GetGroups(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, ok := ctx.Value("user").(*models.User)
	if !ok {
		_ = render.Render(w, r, httputil.NewAPIError(422, "Request Can not be processed"))
		return
	}
	_, err := u.service.GetGroups(user)
	if err != nil {
		_ = render.Render(w, r, httputil.NewAPIError(err))
		return
	}
	if err := render.Render(w, r, NewUserResponse(user)); err != nil {
		_ = render.Render(w, r, httputil.NewAPIError(err))
		return
	}
}

func (u *userHandler) GetPermissions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, ok := ctx.Value("user").(*models.User)
	if !ok {
		_ = render.Render(w, r, httputil.NewAPIError(422, "Request Can not be processed"))
		return
	}
	_, err := u.service.GetPermissions(user)
	if err != nil {
		_ = render.Render(w, r, httputil.NewAPIError(err))
		return
	}
	if err := render.Render(w, r, NewUserResponse(user)); err != nil {
		_ = render.Render(w, r, httputil.NewAPIError(err))
		return
	}
}
