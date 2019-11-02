package groups

import (
	"github.com/go-pg/pg/v9"
	"github.com/imtanmoy/authz/models"
)

type Repository interface {
	List() ([]*models.Group, error)
	Create(tx *pg.Tx, group *models.Group) (*models.Group, error)
}

type groupRepository struct {
	db *pg.DB
}

var _ Repository = (*groupRepository)(nil)

func NewGroupRepository(db *pg.DB) Repository {
	return &groupRepository{
		db,
	}
}

func (g *groupRepository) List() ([]*models.Group, error) {
	var groups []*models.Group
	err := g.db.Model(&groups).Relation("Organization").Select()
	return groups, err
}

func (g *groupRepository) Create(tx *pg.Tx, group *models.Group) (*models.Group, error) {
	_, err := g.db.Model(group).Returning("*").Insert()
	return group, err
}
