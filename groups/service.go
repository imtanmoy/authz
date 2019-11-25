package groups

import (
	"github.com/go-pg/pg/v9"
	"github.com/imtanmoy/authz/casbin"
	"github.com/imtanmoy/authz/models"
	"github.com/imtanmoy/authz/users"
	"strconv"
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
		_, err = casbin.Enforcer.AddPolicy(strconv.Itoa(int(newGroup.ID)), strconv.Itoa(int(permission.ID)), permission.Action)
		if err != nil {
			return nil, err
		}
	}
	for _, user := range users {
		_, err = casbin.Enforcer.AddGroupingPolicy(strconv.Itoa(int(user.ID)), strconv.Itoa(int(newGroup.ID)))
		if err != nil {
			return nil, err
		}
	}
	return newGroup, err
}

func (g *groupService) FindByName(organization *models.Organization, name string) (*models.Group, error) {
	return g.repository.FindByName(organization, name)
}
