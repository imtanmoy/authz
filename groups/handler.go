package groups

import (
	"context"
	"github.com/go-chi/render"
	"github.com/go-pg/pg/v9"
	"github.com/imtanmoy/authz/models"
	"github.com/imtanmoy/authz/organizations"
	"github.com/imtanmoy/authz/utils/errutil"
	param "github.com/oceanicdev/chi-param"
	"net/http"
	"strings"
)

type Handler interface {
	OrganizationCtx(next http.Handler) http.Handler
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

func (g *groupHandler) OrganizationCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		oid, err := param.Int32(r, "oid")
		if err != nil {
			_ = render.Render(w, r, errutil.ErrInvalidRequestParam())
			return
		}
		organization, err := g.organizationService.Find(oid)
		if err != nil {
			_ = render.Render(w, r, errutil.ErrNotFound("organization not found"))
			return
		}
		ctx := context.WithValue(r.Context(), "organization", organization)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (g *groupHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	organization, ok := ctx.Value("organization").(*models.Organization)
	if !ok {
		_ = render.Render(w, r, errutil.ErrUnprocessableEntity())
		return
	}
	groups, err := g.service.List(organization)
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
	ctx := r.Context()
	organization, ok := ctx.Value("organization").(*models.Organization)
	if !ok {
		_ = render.Render(w, r, errutil.ErrUnprocessableEntity())
		return
	}
	data := &GroupPayload{}
	if err := render.Bind(r, data); err != nil {
		_ = render.Render(w, r, errutil.ErrRender(err))
		return
	}
	// request validation
	validationErrors := data.validate()
	if len(validationErrors) > 0 {
		_ = render.Render(w, r, errutil.ErrInvalidRequest(validationErrors))
		return
	}

	// check if users belongs to the organization
	users, _ := g.organizationService.FindUsersByIds(organization, data.Users)
	if len(users) != len(data.Users) {
		_ = render.Render(w, r, errutil.ErrInvalidRequestParamWithMsg("invalid users list"))
		return
	}

	// check if group with same name already exist
	existGroup, err := g.service.FindByName(organization, data.Name)
	existErr := make(map[string][]string)
	if err == nil && existGroup.Name == data.Name {
		existErr = map[string][]string{
			"name": {"Group with same name already exits"},
		}
	}
	if len(existErr) > 0 {
		_ = render.Render(w, r, errutil.ErrInvalidRequest(existErr))
		return
	}
	if err != nil && strings.Contains(err.Error(), "no rows in result set") {
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
	} else {
		_ = render.Render(w, r, errutil.ErrRender(err))
		return
	}
}
