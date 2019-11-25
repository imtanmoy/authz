package users

import (
	"net/http"
	"net/url"
	"time"

	"github.com/go-chi/render"
	"gopkg.in/thedevsaddam/govalidator.v1"

	"github.com/imtanmoy/authz/models"
)

type organizationResponse struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

type groupResponse struct {
	ID        int32     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type UserResponse struct {
	ID           int32                 `json:"id"`
	Email        string                `json:"email"`
	Organization *organizationResponse `json:"organization"`
	Groups       []*groupResponse      `json:"groups"`
}

func (u *UserResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func newGroupsResponse(group *models.Group) *groupResponse {
	return &groupResponse{
		ID:        group.ID,
		Name:      group.Name,
		CreatedAt: group.CreatedAt,
	}
}

func NewUserResponse(user *models.User) *UserResponse {
	var groups []*groupResponse
	//for _, group := range user.Groups {
	//	groups = append(groups, newGroupsResponse(group))
	//}
	return &UserResponse{
		ID:     user.ID,
		Email:  user.Email,
		Groups: groups,
		Organization: &organizationResponse{
			ID:   user.Organization.ID,
			Name: user.Organization.Name,
		}}
}

type UserPayload struct {
	ID             int32  `json:"id"`
	Email          string `json:"email"`
	OrganizationID int32  `json:"organization_id"`
}

func (u *UserPayload) Bind(r *http.Request) error {
	return nil
}

func (u *UserPayload) validate() url.Values {
	rules := govalidator.MapData{
		"id":              []string{"required"},
		"email":           []string{"required", "email"},
		"organization_id": []string{"required"},
	}
	opts := govalidator.Options{
		Data:  u,
		Rules: rules,
	}

	v := govalidator.New(opts)
	e := v.ValidateStruct()
	return e
}

func NewUserListResponse(users []*models.User) []render.Renderer {
	var list []render.Renderer
	for _, user := range users {
		list = append(list, NewUserResponse(user))
	}
	return list
}
