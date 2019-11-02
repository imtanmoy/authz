package organizations

import (
	"net/http"
	"net/url"

	"github.com/go-chi/render"
	"gopkg.in/thedevsaddam/govalidator.v1"

	"github.com/imtanmoy/authz/models"
)

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

type OrganizationPayload struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

func (o *OrganizationPayload) Bind(r *http.Request) error {
	return nil
}

func (o *OrganizationPayload) validate() url.Values {
	rules := govalidator.MapData{
		"id":   []string{"required"},
		"name": []string{"required", "min:4", "max:20"},
	}
	opts := govalidator.Options{
		Data:  o,
		Rules: rules,
	}

	v := govalidator.New(opts)
	e := v.ValidateStruct()
	return e
}

type OrganizationResponse struct {
	ID    int32           `json:"id"`
	Name  string          `json:"name"`
	Users []*userResponse `json:"users"`
}

func (o *OrganizationResponse) Render(w http.ResponseWriter, r *http.Request) error {
	if o.Users == nil {
		o.Users = make([]*userResponse, 0)
	}
	return nil
}

func NewOrganizationResponse(organization *models.Organization) *OrganizationResponse {
	var list []*userResponse
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
