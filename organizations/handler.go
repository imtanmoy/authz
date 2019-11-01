package organizations

import (
	"github.com/go-chi/render"
	"github.com/go-pg/pg/v9"
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

	validationErrors := data.validate()

	if len(validationErrors) > 0 {
		_ = render.Render(w, r, ErrInvalidRequest(validationErrors))
		return
	}
	exist := o.service.Exists(data.ID)
	if exist {
		existErr := map[string][]string{
			"id": {"Organization with same id already exits"},
		}
		_ = render.Render(w, r, ErrInvalidRequest(existErr))
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
	return
}

func (o *organizationHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := param.Int(r, "id")
	if err != nil {
		_ = render.Render(w, r, ErrInvalidRequestParam())
		return
	}
	organization, err := o.service.Find(id)
	if err != nil {
		_ = render.Render(w, r, ErrRender(err))
		return
	}
	if err := render.Render(w, r, NewOrganizationResponse(organization)); err != nil {
		_ = render.Render(w, r, ErrRender(err))
		return
	}
}

func (o *organizationHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := param.Int(r, "id")
	if err != nil {
		_ = render.Render(w, r, ErrInvalidRequestParam())
		return
	}

	data := &OrganizationPayload{}
	if err := render.Bind(r, data); err != nil {
		_ = render.Render(w, r, ErrRender(err))
		return
	}

	validationErrors := data.validate()

	if len(validationErrors) > 0 {
		_ = render.Render(w, r, ErrInvalidRequest(validationErrors))
		return
	}

	organization, err := o.service.Find(id)
	if err != nil {
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
}

func (o *organizationHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := param.Int(r, "id")
	if err != nil {
		_ = render.Render(w, r, ErrInvalidRequestParam())
		return
	}
	organization, err := o.service.Find(id)
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
}

var _ Handler = (*organizationHandler)(nil)
