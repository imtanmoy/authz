package users

import (
	"github.com/go-chi/render"
	"github.com/go-pg/pg/v9"
	"github.com/imtanmoy/authz/organizations"
	"github.com/imtanmoy/authz/utils/httputil"
	"github.com/oceanicdev/chi-param"
	"net/http"

	"github.com/imtanmoy/authz/models"
)

type Handler interface {
	List(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
}

type userHandler struct {
	service             Service
	organizationService organizations.Service
	db                  *pg.DB
}

func NewUserHandler(db *pg.DB) Handler {
	return &userHandler{
		service:             NewUserService(db),
		organizationService: organizations.NewOrganizationService(db),
		db:                  db,
	}
}

func (u *userHandler) List(w http.ResponseWriter, r *http.Request) {
	users, err := u.service.List()
	if err != nil {
		_ = render.Render(w, r, httputil.ErrRender(err))
		return
	}
	if err := render.RenderList(w, r, NewUserListResponse(users)); err != nil {
		_ = render.Render(w, r, httputil.ErrRender(err))
		return
	}
}

func (u *userHandler) Create(w http.ResponseWriter, r *http.Request) {

	data := &UserPayload{}
	if err := render.Bind(r, data); err != nil {
		_ = render.Render(w, r, httputil.ErrRender(err))
		return
	}

	validationErrors := data.validate()

	if len(validationErrors) > 0 {
		_ = render.Render(w, r, httputil.ErrInvalidRequest(validationErrors))
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
		_ = render.Render(w, r, httputil.ErrInvalidRequest(existErr))
		return
	}

	organization, err := u.organizationService.Find(data.OrganizationID)
	if err != nil {
		_ = render.Render(w, r, httputil.ErrRender(err))
		return
	}

	var user models.User
	user.ID = data.ID
	user.Email = data.Email
	user.Organization = organization
	user.OrganizationID = organization.ID

	newUser, err := u.service.Create(&user)

	if err != nil {
		_ = render.Render(w, r, httputil.ErrRender(err))
		return
	}

	render.Status(r, http.StatusCreated)
	_ = render.Render(w, r, NewUserResponse(newUser))
	return
}

func (u *userHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := param.Int32(r, "id")
	if err != nil {
		_ = render.Render(w, r, httputil.ErrInvalidRequestParam())
		return
	}
	user, err := u.service.Find(id)
	if err != nil {
		_ = render.Render(w, r, httputil.ErrNotFound("user not found"))
		return
	}
	if err := render.Render(w, r, NewUserResponse(user)); err != nil {
		_ = render.Render(w, r, httputil.ErrRender(err))
		return
	}
}

func (u *userHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := param.Int32(r, "id")
	if err != nil {
		_ = render.Render(w, r, httputil.ErrInvalidRequestParam())
		return
	}

	data := &UserPayload{}
	if err := render.Bind(r, data); err != nil {
		_ = render.Render(w, r, httputil.ErrRender(err))
		return
	}

	validationErrors := data.validate()

	if len(validationErrors) > 0 {
		_ = render.Render(w, r, httputil.ErrInvalidRequest(validationErrors))
		return
	}

	user, err := u.service.Find(id)
	if err != nil {
		_ = render.Render(w, r, httputil.ErrRender(err))
		return
	}

	user.Email = data.Email
	user, err = u.service.Update(user)
	if err != nil {
		_ = render.Render(w, r, httputil.ErrRender(err))
		return
	}
	if err := render.Render(w, r, NewUserResponse(user)); err != nil {
		_ = render.Render(w, r, httputil.ErrRender(err))
		return
	}
}

func (u *userHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := param.Int32(r, "id")
	if err != nil {
		_ = render.Render(w, r, httputil.ErrInvalidRequestParam())
		return
	}
	user, err := u.service.Find(id)
	if err != nil {
		_ = render.Render(w, r, httputil.ErrRender(err))
		return
	}
	err = u.service.Delete(user)
	if err != nil {
		_ = render.Render(w, r, httputil.ErrRender(err))
		return
	}
	render.NoContent(w, r)
}

var _ Handler = (*userHandler)(nil)
