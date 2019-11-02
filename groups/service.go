package groups

import (
	"github.com/go-pg/pg/v9"
	"github.com/imtanmoy/authz/models"
)

type Service interface {
	List() ([]*models.Group, error)
	Create(group *models.Group) (*models.Group, error)
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

func (g *groupService) List() ([]*models.Group, error) {
	return g.repository.List()
}

func (g *groupService) Create(group *models.Group) (*models.Group, error) {
	tx, _ := g.db.Begin()
	defer tx.Commit()
	return g.repository.Create(tx, group)
}
