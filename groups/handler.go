package groups

import (
	"github.com/go-chi/render"
	"github.com/go-pg/pg/v9"
	"github.com/imtanmoy/authz/models"
	"github.com/imtanmoy/authz/organizations"
	"github.com/imtanmoy/authz/utils/errutil"
	"net/http"
)

type Handler interface {
	List(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
}

type groupHandler struct {
	service             Service
	organizationService organizations.Service
	db                  *pg.DB
}

var _ Handler = (*groupHandler)(nil)

func NewGroupHandler(db *pg.DB) Handler {
	return &groupHandler{
		service:             NewGroupService(db),
		organizationService: organizations.NewOrganizationService(db),
		db:                  db,
	}
}

func (g *groupHandler) List(w http.ResponseWriter, r *http.Request) {
	groups, err := g.service.List()
	if err != nil {
		_ = render.Render(w, r, errutil.ErrRender(err))
		return
	}
	if err := render.RenderList(w, r, NewGroupListResponse(groups)); err != nil {
		_ = render.Render(w, r, errutil.ErrRender(err))
		return
	}
}

func (g *groupHandler) Create(w http.ResponseWriter, r *http.Request) {
	data := &GroupPayload{}
	if err := render.Bind(r, data); err != nil {
		_ = render.Render(w, r, errutil.ErrRender(err))
		return
	}

	validationErrors := data.validate()

	if len(validationErrors) > 0 {
		_ = render.Render(w, r, errutil.ErrInvalidRequest(validationErrors))
		return
	}
	orgExist := g.organizationService.Exists(data.OrganizationID)
	existErr := make(map[string][]string)
	//if exist {
	//	existErr = map[string][]string{
	//		"id": {"User with same id already exits"},
	//	}
	//}
	if !orgExist {
		existErr = map[string][]string{
			"organization_id": {"organization does not exist"},
		}
	}
	if len(existErr) > 0 {
		_ = render.Render(w, r, errutil.ErrInvalidRequest(existErr))
		return
	}

	organization, err := g.organizationService.Find(data.OrganizationID)
	if err != nil {
		_ = render.Render(w, r, errutil.ErrRender(err))
		return
	}

	var group models.Group
	group.Name = data.Name
	group.Organization = organization
	group.OrganizationID = organization.ID

	newGroup, err := g.service.Create(&group)

	if err != nil {
		_ = render.Render(w, r, errutil.ErrRender(err))
		return
	}

	render.Status(r, http.StatusCreated)
	_ = render.Render(w, r, NewGroupResponse(newGroup))
	return
}
