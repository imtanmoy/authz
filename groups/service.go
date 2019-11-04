package groups

import (
	"github.com/go-pg/pg/v9"
	"github.com/imtanmoy/authz/models"
	"github.com/imtanmoy/authz/users"
)

type Service interface {
	List(organization *models.Organization) ([]*models.Group, error)
	Create(groupPayload *GroupPayload, organization *models.Organization) (*models.Group, error)
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

func (g *groupService) Create(groupPayload *GroupPayload, organization *models.Organization) (*models.Group, error) {
	var group models.Group

	tx, err := g.db.Begin()
	if err != nil {
		return nil, err
	}

	group.Name = groupPayload.Name
	group.Organization = organization
	group.OrganizationID = organization.ID
	userList := g.userRepository.FindAllByIdIn(groupPayload.Users)

	newGroup, err := g.repository.Create(tx, &group)
	if err != nil {
		return nil, err
	}
	for _, user := range userList {
		_, err = tx.Model(&models.UserGroup{
			UserId:  user.ID,
			User:    user,
			GroupId: newGroup.ID,
			Group:   newGroup,
		}).Insert()
		if err != nil {
			_ = tx.Rollback()
			return nil, err
		}
	}
	newGroup.Users = userList
	tx.Commit()
	return newGroup, err
}

func (g *groupService) FindByName(organization *models.Organization, name string) (*models.Group, error) {
	return g.repository.FindByName(organization, name)
}
