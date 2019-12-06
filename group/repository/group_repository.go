package repository

import (
	"context"
	"errors"
	"github.com/go-pg/pg/v9"
	"github.com/imtanmoy/authz/group"
	"github.com/imtanmoy/authz/models"
)

type groupRepository struct {
	db *pg.DB
}

var _ group.Repository = (*groupRepository)(nil)

// NewRepository will create an object that represent the group.Repository interface
func NewRepository(db *pg.DB) group.Repository {
	return &groupRepository{db}
}

func (g *groupRepository) Fetch(ctx context.Context, organizationId int32) ([]*models.Group, error) {
	db := g.db.WithContext(ctx)
	var groups []*models.Group
	err := db.Model(&groups).Where("organization_id = ?", organizationId).Relation("Organization").Select()
	return groups, err
}

func (g *groupRepository) Store(ctx context.Context, group *models.Group) (*models.Group, error) {
	db := g.db.WithContext(ctx)
	_, err := db.Model(group).Returning("*").Insert()
	return group, err
}

func (g *groupRepository) Find(ctx context.Context, ID int32) (*models.Group, error) {
	db := g.db.WithContext(ctx)
	if !g.Exists(ctx, ID) {
		return nil, errors.New("group not found")
	}
	var grp models.Group
	err := db.Model(&grp).Where("\"group\".id = ?", ID).Relation("Organization").Select()
	return &grp, err
}

func (g *groupRepository) Exists(ctx context.Context, ID int32) bool {
	db := g.db.WithContext(ctx)
	var num int32
	_, err := db.Query(pg.Scan(&num), "SELECT id from groups where id = ?", ID)
	if err != nil {
		panic(err)
	}
	return num == ID
}

func (g *groupRepository) FindByIdAndOrganizationId(ctx context.Context, Id int32, Oid int32) (*models.Group, error) {
	db := g.db.WithContext(ctx)
	var grp models.Group
	err := db.Model(&grp).
		Where("\"group\".id = ?", Id).
		Where("\"group\".organization_id = ?", Oid).
		Relation("Organization").Select()
	return &grp, err
}

func (g *groupRepository) Update(ctx context.Context, group *models.Group) error {
	db := g.db.WithContext(ctx)
	err := db.Update(group)
	return err
}

func (g *groupRepository) Delete(ctx context.Context, group *models.Group) error {
	db := g.db.WithContext(ctx)
	err := db.Delete(group)
	return err
}

func (g *groupRepository) FindAllByIdIn(ctx context.Context, ids []int32) []*models.Group {
	db := g.db.WithContext(ctx)
	var groups []*models.Group
	_ = db.Model(&groups). // TODO err handling
		Where("id in (?)", pg.In(ids)).
		Select()
	return groups
}

func (g *groupRepository) FindByNameAndOrganizationId(ctx context.Context, name string, oid int32) (*models.Group, error) {
	db := g.db.WithContext(ctx)
	var grp models.Group
	err := db.Model(&grp).
		Where("name = ?", name).
		Where("organization_id = ?", oid).
		First()
	return &grp, err
}
