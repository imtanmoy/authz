package groups

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

type organizationResponse struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

type GroupResponse struct {
	ID           int32                 `json:"id"`
	Name         string                `json:"name"`
	CreatedAt    time.Time             `json:"created_at"`
	Organization *organizationResponse `json:"organization"`
	Users        []*userResponse       `json:"users"`
}

func (g *GroupResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func NewGroupResponse(group *models.Group) *GroupResponse {
	var users []*userResponse
	for _, user := range group.Users {
		users = append(users, NewUserResponse(user))
	}
	return &GroupResponse{
		ID:        group.ID,
		Name:      group.Name,
		CreatedAt: group.CreatedAt,
		Users:     users,
		Organization: &organizationResponse{
			ID:   group.Organization.ID,
			Name: group.Organization.Name,
		}}
}

func NewGroupListResponse(groups []*models.Group) []render.Renderer {
	var list []render.Renderer
	for _, group := range groups {
		list = append(list, NewGroupResponse(group))
	}
	return list
}
