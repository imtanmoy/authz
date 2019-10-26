package organizations

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/go-pg/pg/v9"
	
	"github.com/imtanmoy/authz/models"
)

type Handler interface {
	List(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
}

type organizationHandler struct {
	service Service
	db      *pg.DB
}

func NewOrganizationHandler(db *pg.DB) Handler {
	return &organizationHandler{
		service: NewOrganizationService(db),
		db:      db,
	}
}

type UserResponse struct {
	ID    int32  `json:"id"`
	Email string `json:"email"`
}

func (u *UserResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func NewUserResponse(user *models.User) *UserResponse {
	resp := &UserResponse{ID: user.ID, Email: user.Email}
	return resp
}

type OrganizationPayload struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (o *OrganizationPayload) Bind(r *http.Request) error {
	return nil
}

type OrganizationResponse struct {
	ID    int32           `json:"id"`
	Name  string          `json:"name"`
	Users []*UserResponse `json:"users"`
}

func (o *OrganizationResponse) Render(w http.ResponseWriter, r *http.Request) error {
	if o.Users == nil {
		o.Users = make([]*UserResponse, 0)
	}
	return nil
}

func NewOrganizationResponse(organization *models.Organization) *OrganizationResponse {
	var list []*UserResponse
	for _, user := range organization.Users {
		list = append(list, NewUserResponse(user))
	}
	resp := &OrganizationResponse{
		ID:    organization.ID,
		Name:  organization.Name,
		Users: list,
	}
	return resp
}

func NewOrganizationListResponse(organizations []*models.Organization) []render.Renderer {
	var list []render.Renderer
	for _, organization := range organizations {
		list = append(list, NewOrganizationResponse(organization))
	}
	return list
}

func (o *organizationHandler) List(w http.ResponseWriter, r *http.Request) {
	organizations, err := o.service.List()
	if err != nil {
		_ = render.Render(w, r, ErrRender(err))
		return
	}

	if err := render.RenderList(w, r, NewOrganizationListResponse(organizations)); err != nil {
		_ = render.Render(w, r, ErrRender(err))
		return
	}
}

func (o *organizationHandler) Create(w http.ResponseWriter, r *http.Request) {

	data := &OrganizationPayload{}
	if err := render.Bind(r, data); err != nil {
		_ = render.Render(w, r, ErrRender(err))
		return
	}

	var organization models.Organization
	organization.ID = int32(data.ID)
	organization.Name = data.Name

	newOrganization, err := o.service.Create(&organization)
	if err != nil {
		_ = render.Render(w, r, ErrRender(err))
		return
	}

	render.Status(r, http.StatusCreated)
	_ = render.Render(w, r, NewOrganizationResponse(newOrganization))
}

func (o *organizationHandler) Get(w http.ResponseWriter, r *http.Request) {
	var err error
	if id := chi.URLParam(r, "id"); id != "" {
		var oid int
		oid, err = strconv.Atoi(id)
		if err != nil {
			_ = render.Render(w, r, ErrRender(err))
			return
		}
		organization, err := o.service.Find(oid)
		if err != nil {
			_ = render.Render(w, r, ErrRender(err))
			return
		}
		if err := render.Render(w, r, NewOrganizationResponse(organization)); err != nil {
			_ = render.Render(w, r, ErrRender(err))
			return
		}
	} else {
		_ = render.Render(w, r, ErrRender(errors.New("invalid argument")))
		return
	}
}

func (o *organizationHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id != "" {
		var oid int
		oid, err := strconv.Atoi(id)
		if err != nil {
			_ = render.Render(w, r, ErrRender(err))
			return
		}
		organization, err := o.service.Find(oid)
		if err != nil {
			_ = render.Render(w, r, ErrRender(err))
			return
		}
		data := &OrganizationPayload{}
		if err := render.Bind(r, data); err != nil {
			_ = render.Render(w, r, ErrRender(err))
			return
		}
		organization.Name = data.Name
		organization, err = o.service.Update(organization)
		if err != nil {
			_ = render.Render(w, r, ErrRender(err))
			return
		}
		if err := render.Render(w, r, NewOrganizationResponse(organization)); err != nil {
			_ = render.Render(w, r, ErrRender(err))
			return
		}
	} else {
		_ = render.Render(w, r, ErrRender(errors.New("invalid argument")))
		return
	}
}

func (o *organizationHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id != "" {
		var oid int
		oid, err := strconv.Atoi(id)
		if err != nil {
			_ = render.Render(w, r, ErrRender(err))
			return
		}
		organization, err := o.service.Find(oid)
		if err != nil {
			_ = render.Render(w, r, ErrRender(err))
			return
		}
		err = o.service.Delete(organization)
		if err != nil {
			_ = render.Render(w, r, ErrRender(err))
			return
		}
		render.NoContent(w, r)

	} else {
		_ = render.Render(w, r, ErrRender(errors.New("invalid argument")))
		return
	}
}

var _ Handler = (*organizationHandler)(nil)
