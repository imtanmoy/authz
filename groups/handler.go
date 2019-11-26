package groups

import (
	"context"
	"fmt"
	"github.com/go-chi/render"
	"github.com/go-pg/pg/v9"
	"github.com/imtanmoy/authz/casbin"
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
	List(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
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
	for _, group := range groups {
		userList, err := casbin.Enforcer.GetUsersForRole(fmt.Sprintf("group::%d", group.ID))
		var uIDs []int32

		if err == nil {
			for _, user := range userList {
				uIDs = append(uIDs, GetIntID(user))
			}
			group.Users = g.userService.FindAllByIdIn(uIDs)
		}
		pList, err := casbin.Enforcer.GetImplicitPermissionsForUser(fmt.Sprintf("group::%d", group.ID))
		var pIDs []int32
		if err == nil {
			for _, p := range pList {
				pIDs = append(pIDs, GetIntID(p[1]))
			}
			group.Permissions = g.permissionService.FindAllByIdIn(pIDs)
		}
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
	newGroup.Users = userList
	newGroup.Permissions = permissionList

	render.Status(r, http.StatusCreated)
	_ = render.Render(w, r, NewGroupResponse(newGroup))
	return
}
