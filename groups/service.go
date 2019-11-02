package groups

import (
	"github.com/go-pg/pg/v9"
	"github.com/imtanmoy/authz/models"
)

type Service interface {
	List(organization *models.Organization) ([]*models.Group, error)
	Create(group *models.Group) (*models.Group, error)
	FindByName(organization *models.Organization, name string) (*models.Group, error)
}

type groupService struct {
	db         *pg.DB
	repository Repository
}

var _ Service = (*groupService)(nil)

func NewGroupService(db *pg.DB) Service {
	return &groupService{
		repository: NewGroupRepository(db),
		db:         db,
	}
}

func (g *groupService) List(organization *models.Organization) ([]*models.Group, error) {
	return g.repository.List(organization)
}

func (g *groupService) Create(group *models.Group) (*models.Group, error) {
	tx, _ := g.db.Begin()
	defer tx.Commit()
	return g.repository.Create(tx, group)
}

func (g *groupService) FindByName(organization *models.Organization, name string) (*models.Group, error) {
	return g.repository.FindByName(organization, name)
}
