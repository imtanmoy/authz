package repository

import (
	"errors"
	"github.com/go-pg/pg/v9"
	"github.com/imtanmoy/authz/group"
	"github.com/imtanmoy/authz/models"
)

type pgGroupRepository struct {
	db *pg.DB
}

var _ group.Repository = (*pgGroupRepository)(nil)

// NewPgGroupRepository will create an object that represent the group.Repository interface
func NewPgGroupRepository(db *pg.DB) group.Repository {
	return &pgGroupRepository{db}
}

func (g *pgGroupRepository) List(organizationId int32) ([]*models.Group, error) {
	var groups []*models.Group
	err := g.db.Model(&groups).Where("organization_id = ?", organizationId).Relation("Organization").Select()
	return groups, err
}

func (g *pgGroupRepository) Create(tx *pg.Tx, group *models.Group) (*models.Group, error) {
	_, err := tx.Model(group).Returning("*").Insert()
	return group, err
}

func (g *pgGroupRepository) FindByName(organization *models.Organization, name string) (*models.Group, error) {
	var grp models.Group
	err := g.db.Model(&grp).
		Where("name = ?", name).
		Where("organization_id = ?", organization.ID).
		First()
	return &grp, err
}

func (g *pgGroupRepository) Find(ID int32) (*models.Group, error) {
	if !g.Exists(ID) {
		return nil, errors.New("group not found")
	}
	var grp models.Group
	err := g.db.Model(&grp).Where("\"group\".id = ?", ID).Relation("Organization").Select()
	return &grp, err
}

func (g *pgGroupRepository) Exists(ID int32) bool {
	var num int32
	_, err := g.db.Query(pg.Scan(&num), "SELECT id from groups where id = ?", ID)
	if err != nil {
		panic(err)
	}
	return num == ID
}

func (g *pgGroupRepository) FindByIdAndOrganizationId(Id int32, Oid int32) (*models.Group, error) {
	var grp models.Group
	err := g.db.Model(&grp).
		Where("\"group\".id = ?", Id).
		Where("\"group\".organization_id = ?", Oid).
		Relation("Organization").Select()
	return &grp, err
}

func (g *pgGroupRepository) Update(tx *pg.Tx, group *models.Group) error {
	err := tx.Update(group)
	return err
}

func (g *pgGroupRepository) Delete(group *models.Group) error {
	err := g.db.Delete(group)
	return err
}

func (g *pgGroupRepository) FindAllByIdIn(ids []int32) []*models.Group {
	var groups []*models.Group
	_ = g.db.Model(&groups). // TODO err handling
		Where("id in (?)", pg.In(ids)).
		Select()
	return groups
}
