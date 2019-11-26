package groups

import (
	"fmt"
	"github.com/go-pg/pg/v9"
	"github.com/imtanmoy/authz/casbin"
	"github.com/imtanmoy/authz/models"
	"github.com/imtanmoy/authz/users"
)

type Service interface {
	List(organization *models.Organization) ([]*models.Group, error)
	Create(groupPayload *GroupPayload, organization *models.Organization, users []*models.User, permissions []*models.Permission) (*models.Group, error)
	FindByName(organization *models.Organization, name string) (*models.Group, error)
}

type groupService struct {
	db             *pg.DB
	repository     Repository
	userRepository users.Repository
}

var _ Service = (*groupService)(nil)

func NewGroupService(db *pg.DB) Service {
	return &groupService{
		db:             db,
		repository:     NewGroupRepository(db),
		userRepository: users.NewUserRepository(db),
	}
}

func (g *groupService) List(organization *models.Organization) ([]*models.Group, error) {
	return g.repository.List(organization)
}

func (g *groupService) Create(
	groupPayload *GroupPayload,
	organization *models.Organization,
	users []*models.User,
	permissions []*models.Permission,
) (*models.Group, error) {
	var group models.Group

	tx, err := g.db.Begin()
	if err != nil {
		return nil, err
	}

	group.Name = groupPayload.Name
	group.Organization = organization
	group.OrganizationID = organization.ID

	newGroup, err := g.repository.Create(tx, &group)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}
	err = tx.Commit()

	for _, permission := range permissions {
		groupID := fmt.Sprintf("group::%d", newGroup.ID)
		permissionID := fmt.Sprintf("permission::%d", permission.ID)
		params := []string{groupID, permissionID, permission.Action}
		_, err = casbin.Enforcer.AddPolicy(params)
		if err != nil {
			return nil, err
		}
	}
	fmt.Println(newGroup.UpdatedAt)
	for _, user := range users {
		userID := fmt.Sprintf("user::%d", user.ID)
		groupID := fmt.Sprintf("group::%d", newGroup.ID)
		params := []string{userID, groupID}
		_, err = casbin.Enforcer.AddGroupingPolicy(params)
		if err != nil {
			return nil, err
		}
	}
	return newGroup, err
}

func (g *groupService) FindByName(organization *models.Organization, name string) (*models.Group, error) {
	return g.repository.FindByName(organization, name)
}
