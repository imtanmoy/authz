package http

import (
	"context"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/imtanmoy/authz/group"
	"github.com/imtanmoy/authz/models"
	"github.com/imtanmoy/authz/organizations"
	"github.com/imtanmoy/authz/utils/httputil"
	param "github.com/oceanicdev/chi-param"
	"net/http"
)

// GroupHandler  represent the http handler for group
type GroupHandler struct {
	GUsecase            group.Usecase
	organizationService organizations.Service
}

// NewGroupHandler will initialize the articles/ resources endpoint
func NewGroupHandler(r *chi.Mux, us group.Usecase, organizationService organizations.Service) {
	handler := &GroupHandler{
		GUsecase:            us,
		organizationService: organizationService,
	}
	r.Group(func(r chi.Router) {
		r.Use(handler.OrganizationCtx)
		r.Get("/{oid}/groups", handler.Fetch)
		r.Post("/{oid}/groups", handler.Store)
		r.Get("/{oid}/groups/{id}", handler.GetByID)
		r.Put("/{oid}/groups/{id}", handler.Update)
		r.Delete("/{oid}/groups/{id}", handler.Delete)
	})
}

func (gh *GroupHandler) OrganizationCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		oid, err := param.Int32(r, "oid")
		if err != nil {
			_ = render.Render(w, r, httputil.NewAPIError(400, "Invalid request parameter", err))
			return
		}
		organization, err := gh.organizationService.Find(oid)
		if err != nil {
			_ = render.Render(w, r, httputil.NewAPIError(404, "organization not found", err))
			return
		}
		ctx := context.WithValue(r.Context(), "organization", organization)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (gh *GroupHandler) Fetch(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if ctx == nil {
		//ctx = context.Background()
		_ = render.Render(w, r, httputil.NewAPIError(500, "Something went wrong"))
		return
	}
	organization, ok := ctx.Value("organization").(*models.Organization)
	if !ok {
		_ = render.Render(w, r, httputil.NewAPIError(422, "Request Can not be processed"))
		return
	}
	groups, err := gh.GUsecase.Fetch(ctx, organization.ID)
	if err != nil {
		_ = render.Render(w, r, httputil.NewAPIError(err))
		return
	}
	if err := render.RenderList(w, r, group.NewGroupListResponse(groups)); err != nil {
		_ = render.Render(w, r, httputil.NewAPIError(err))
		return
	}
}

func (gh *GroupHandler) Store(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if ctx == nil {
		//ctx = context.Background()
		_ = render.Render(w, r, httputil.NewAPIError(500, "Something went wrong"))
		return
	}
	organization, ok := ctx.Value("organization").(*models.Organization)
	if !ok {
		_ = render.Render(w, r, httputil.NewAPIError(422, "Request Can not be processed"))
		return
	}

	data := &group.GroupPayload{}
	if err := render.Bind(r, data); err != nil {
		_ = render.Render(w, r, httputil.NewAPIError(422, "unable to decode the request content type"))
		return
	}

	// request validation
	validationErrors := data.Validate()
	// check if users belongs to the organization
	userList, _ := gh.organizationService.FindUsersByIds(organization, data.Users)
	if len(userList) != len(data.Users) {
		validationErrors.Add("users", "invalid user list")
	}
	// check if permissions belongs to the organization
	permissionList, _ := gh.organizationService.FindPermissionsByIds(organization, data.Permissions)
	if len(permissionList) != len(data.Permissions) {
		validationErrors.Add("permissions", "invalid permission list")
	}

	if len(validationErrors) > 0 {
		_ = render.Render(w, r, httputil.NewAPIError(400, "Invalid request", validationErrors))
		return
	}

	// check if group with same name already exist
	existGroup, err := gh.GUsecase.FindByName(ctx, organization, data.Name)
	if err == nil && existGroup.Name == data.Name {
		validationErrors.Add("name", "Group with same name already exits")
	}
	if len(validationErrors) > 0 {
		_ = render.Render(w, r, httputil.NewAPIError(400, "Invalid Request", validationErrors))
		return
	}
}

func (gh *GroupHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}
}

func (gh *GroupHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}
}

func (gh *GroupHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}
}
