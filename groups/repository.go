package groups

import (
	"errors"
	"fmt"
	"github.com/go-pg/pg/v9"
	"github.com/imtanmoy/authz/models"
)

type Repository interface {
	List(organization *models.Organization) ([]*models.Group, error)
	Create(tx *pg.Tx, group *models.Group) (*models.Group, error)
	FindByName(organization *models.Organization, name string) (*models.Group, error)
	Find(ID int32) (*models.Group, error)
	Exists(ID int32) bool
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
	err := g.db.Model(&groups).Where("organization_id = ?", organization.ID).Relation("Organization").Select()
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

func (g *groupRepository) Find(ID int32) (*models.Group, error) {
	if !g.Exists(ID) {
		return nil, errors.New("group does not exist")
	}
	var group models.Group
	err := g.db.Model(&group).Where("groups.id = ?", ID).Relation("Organization").Select()
	fmt.Println("repo")
	fmt.Println(ID)
	fmt.Println(group)
	fmt.Println(err)
	return &group, err
}

func (g *groupRepository) Exists(ID int32) bool {
	var num int32
	_, err := g.db.Query(pg.Scan(&num), "SELECT id from groups where id = ?", ID)
	if err != nil {
		panic(err)
	}
	return num == ID
}
