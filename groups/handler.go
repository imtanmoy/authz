package groups

import (
	"context"
	"github.com/go-chi/render"
	"github.com/go-pg/pg/v9"
	"github.com/imtanmoy/authz/models"
	"github.com/imtanmoy/authz/organizations"
	"github.com/imtanmoy/authz/permissions"
	"github.com/imtanmoy/authz/users"
	"github.com/imtanmoy/authz/utils/httputil"
	param "github.com/oceanicdev/chi-param"
	"net/http"
)

// Handler handles groups http method
type Handler interface {
	OrganizationCtx(next http.Handler) http.Handler
	GroupCtx(next http.Handler) http.Handler
	List(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
}

type groupHandler struct {
	service             Service
	organizationService organizations.Service
	userService         users.Service
	permissionService   permissions.Service
	db                  *pg.DB
}

var _ Handler = (*groupHandler)(nil)

// NewGroupHandler construct group handler
func NewGroupHandler(db *pg.DB) Handler {
	return &groupHandler{
		service:             NewGroupService(db),
		organizationService: organizations.NewOrganizationService(db),
		userService:         users.NewUserService(db),
		permissionService:   permissions.NewPermissionService(db),
		db:                  db,
	}
}

func (g *groupHandler) OrganizationCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		oid, err := param.Int32(r, "oid")
		if err != nil {
			_ = render.Render(w, r, httputil.NewAPIError(400, "Invalid request parameter", err))
			return
		}
		organization, err := g.organizationService.Find(oid)
		if err != nil {
			_ = render.Render(w, r, httputil.NewAPIError(404, "organization not found", err))
			return
		}
		ctx := context.WithValue(r.Context(), "organization", organization)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (g *groupHandler) GroupCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := param.Int32(r, "id")
		if err != nil {
			_ = render.Render(w, r, httputil.NewAPIError(400, "Invalid request parameter", err))
			return
		}
		ctx := r.Context()
		organization, ok := ctx.Value("organization").(*models.Organization)
		if !ok {
			_ = render.Render(w, r, httputil.NewAPIError(422, "Request Can not be processed"))
			return
		}
		group, err := g.service.FindByIdAndOrganizationId(id, organization.ID)
		if err != nil {
			_ = render.Render(w, r, httputil.NewAPIError(404, "group not found", err))
			return
		}
		ctx = context.WithValue(r.Context(), "group", group)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (g *groupHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	organization, ok := ctx.Value("organization").(*models.Organization)
	if !ok {
		_ = render.Render(w, r, httputil.NewAPIError(422, "Request Can not be processed"))
		return
	}
	groups, err := g.service.List(organization)
	if err != nil {
		_ = render.Render(w, r, httputil.NewAPIError(err))
		return
	}
	if err := render.RenderList(w, r, NewGroupListResponse(groups)); err != nil {
		_ = render.Render(w, r, httputil.NewAPIError(err))
		return
	}
}

func (g *groupHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	organization, ok := ctx.Value("organization").(*models.Organization)
	if !ok {
		_ = render.Render(w, r, httputil.NewAPIError(422, "Request Can not be processed"))
		return
	}
	data := &GroupPayload{}
	if err := render.Bind(r, data); err != nil {
		_ = render.Render(w, r, httputil.NewAPIError(422, "unable to decode the request content type"))
		return
	}

	// request validation
	validationErrors := data.validate()
	// check if users belongs to the organization
	userList, _ := g.organizationService.FindUsersByIds(organization, data.Users)
	if len(userList) != len(data.Users) {
		validationErrors.Add("users", "invalid user list")
	}
	// check if permissions belongs to the organization
	permissionList, _ := g.organizationService.FindPermissionsByIds(organization, data.Permissions)
	if len(permissionList) != len(data.Permissions) {
		validationErrors.Add("permissions", "invalid permission list")
	}

	if len(validationErrors) > 0 {
		_ = render.Render(w, r, httputil.NewAPIError(400, "Invalid request", validationErrors))
		return
	}

	// check if group with same name already exist
	existGroup, err := g.service.FindByName(organization, data.Name)
	if err == nil && existGroup.Name == data.Name {
		validationErrors.Add("name", "Group with same name already exits")
	}
	if len(validationErrors) > 0 {
		_ = render.Render(w, r, httputil.NewAPIError(400, "Invalid Request", validationErrors))
		return
	}

	newGroup, err := g.service.Create(data, organization, userList, permissionList)

	if err != nil {
		_ = render.Render(w, r, httputil.NewAPIError(err))
		return
	}

	render.Status(r, http.StatusCreated)
	_ = render.Render(w, r, NewGroupResponse(newGroup))
	return
}

func (g *groupHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	group, ok := ctx.Value("group").(*models.Group)
	if !ok {
		_ = render.Render(w, r, httputil.NewAPIError(422, "Request Can not be processed"))
		return
	}
	if err := render.Render(w, r, NewGroupResponse(group)); err != nil {
		_ = render.Render(w, r, httputil.NewAPIError(err))
		return
	}
}

func (g *groupHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	group, ok := ctx.Value("group").(*models.Group)
	if !ok {
		_ = render.Render(w, r, httputil.NewAPIError(422, "Request Can not be processed"))
		return
	}
	err := g.service.Delete(group)
	if err != nil {
		_ = render.Render(w, r, httputil.NewAPIError(err))
		return
	}
	render.NoContent(w, r)
}

func (g *groupHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	organization, ok := ctx.Value("organization").(*models.Organization)
	if !ok {
		_ = render.Render(w, r, httputil.NewAPIError(422, "Request Can not be processed"))
		return
	}
	group, ok := ctx.Value("group").(*models.Group)
	if !ok {
		_ = render.Render(w, r, httputil.NewAPIError(422, "Request Can not be processed"))
		return
	}

	data := &GroupPayload{}
	if err := render.Bind(r, data); err != nil {
		_ = render.Render(w, r, httputil.NewAPIError(422, "unable to decode the request content type"))
		return
	}

	// request validation
	validationErrors := data.validate()
	// check if users belongs to the organization
	userList, _ := g.organizationService.FindUsersByIds(organization, data.Users)
	if len(userList) != len(data.Users) {
		validationErrors.Add("users", "invalid user list")
	}
	// check if permissions belongs to the organization
	permissionList, _ := g.organizationService.FindPermissionsByIds(organization, data.Permissions)
	if len(permissionList) != len(data.Permissions) {
		validationErrors.Add("permissions", "invalid permission list")
	}

	if len(validationErrors) > 0 {
		_ = render.Render(w, r, httputil.NewAPIError(400, "Invalid request", validationErrors))
		return
	}
	// check if group with same name already exist
	existGroup, err := g.service.FindByName(organization, data.Name)
	if err == nil && existGroup.Name == data.Name && group.Name != data.Name {
		validationErrors.Add("name", "Group with same name already exits")
	}
	if len(validationErrors) > 0 {
		_ = render.Render(w, r, httputil.NewAPIError(400, "Invalid Request", validationErrors))
		return
	}

	//update group data
	group.Name = data.Name

	err = g.service.Update(group, userList, permissionList)

	if err != nil {
		_ = render.Render(w, r, httputil.NewAPIError(err))
		return
	}
	//group.Users = userList
	//group.Permissions = permissionList

	render.Status(r, http.StatusCreated)
	_ = render.Render(w, r, NewGroupResponse(group))
	return
}
