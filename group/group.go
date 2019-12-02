package group

import (
	"github.com/go-chi/render"
	"github.com/imtanmoy/authz/models"
	"gopkg.in/thedevsaddam/govalidator.v1"
	"net/http"
	"net/url"
	"time"
)

type GroupPayload struct {
	Name        string  `json:"name"`
	Permissions []int32 `json:"permissions"`
	Users       []int32 `json:"users"`
}

func (g *GroupPayload) Bind(r *http.Request) error {
	return nil
}

func (g *GroupPayload) validate() url.Values {
	rules := govalidator.MapData{
		"name": []string{"required"},
	}
	opts := govalidator.Options{
		Data:  g,
		Rules: rules,
	}

	v := govalidator.New(opts)
	e := v.ValidateStruct()
	return e
}

type userResponse struct {
	ID    int32  `json:"id"`
	Email string `json:"email"`
}

func (u *userResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func NewUserResponse(user *models.User) *userResponse {
	resp := &userResponse{ID: user.ID, Email: user.Email}
	return resp
}

type permissionResponse struct {
	ID     int32  `json:"id"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	Action string `json:"action"`
}

func (u *permissionResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func NewPermissionResponse(permission *models.Permission) *permissionResponse {
	resp := &permissionResponse{
		ID:     permission.ID,
		Name:   permission.Name,
		Type:   permission.Type,
		Action: permission.Action,
	}
	return resp
}

type organizationResponse struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

type GroupResponse struct {
	ID           int32                 `json:"id"`
	Name         string                `json:"name"`
	CreatedAt    *time.Time            `json:"created_at"`
	UpdatedAt    *time.Time            `json:"updated_at"`
	Organization *organizationResponse `json:"organization"`
	Users        []*userResponse       `json:"users"`
	Permissions  []*permissionResponse `json:"permissions"`
}

func (g *GroupResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func NewGroupResponse(group *models.Group) *GroupResponse {
	users := make([]*userResponse, 0)
	for _, user := range group.Users {
		users = append(users, NewUserResponse(user))
	}
	permissions := make([]*permissionResponse, 0)
	for _, permission := range group.Permissions {
		permissions = append(permissions, NewPermissionResponse(permission))
	}

	return &GroupResponse{
		ID:          group.ID,
		Name:        group.Name,
		CreatedAt:   &group.CreatedAt,
		UpdatedAt:   &group.UpdatedAt,
		Users:       users,
		Permissions: permissions,
		Organization: &organizationResponse{
			ID:   group.Organization.ID,
			Name: group.Organization.Name,
		}}
}

func NewGroupListResponse(groups []*models.Group) []render.Renderer {
	list := make([]render.Renderer, 0)
	for _, group := range groups {
		list = append(list, NewGroupResponse(group))
	}
	return list
}
