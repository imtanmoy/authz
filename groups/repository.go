package groups

import (
	"github.com/go-pg/pg/v9"
	"github.com/imtanmoy/authz/models"
)

type Repository interface {
	List(organization *models.Organization) ([]*models.Group, error)
	Create(tx *pg.Tx, group *models.Group) (*models.Group, error)
	FindByName(organization *models.Organization, name string) (*models.Group, error)
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

func (g *groupRepository) List(organization *models.Organization) ([]*models.Group, error) {
	var groups []*models.Group
	err := g.db.Model(&groups).Where("organization_id = ?", organization.ID).Relation("Organization").Relation("Users").Select()
	return groups, err
}

func (g *groupRepository) Create(tx *pg.Tx, group *models.Group) (*models.Group, error) {
	_, err := tx.Model(group).Returning("*").Insert()
	return group, err
}

func (g *groupRepository) FindByName(organization *models.Organization, name string) (*models.Group, error) {
	var group models.Group
	err := g.db.Model(&group).
		Where("name = ?", name).
		Where("organization_id = ?", organization.ID).
		First()
	return &group, err
}
