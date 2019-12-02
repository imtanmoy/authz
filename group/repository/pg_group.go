package repository

import (
	"context"
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

func (g *pgGroupRepository) Fetch(ctx context.Context, organizationId int32) ([]*models.Group, error) {
	db := g.db.WithContext(ctx)
	var groups []*models.Group
	err := db.Model(&groups).Where("organization_id = ?", organizationId).Relation("Organization").Select()
	return groups, err
}

func (g *pgGroupRepository) Store(ctx context.Context, group *models.Group) (*models.Group, error) {
	db := g.db.WithContext(ctx)
	_, err := db.Model(group).Returning("*").Insert()
	return group, err
}

func (g *pgGroupRepository) FindByName(ctx context.Context, organization *models.Organization, name string) (*models.Group, error) {
	db := g.db.WithContext(ctx)
	var grp models.Group
	err := db.Model(&grp).
		Where("name = ?", name).
		Where("organization_id = ?", organization.ID).
		First()
	return &grp, err
}

func (g *pgGroupRepository) Find(ctx context.Context, ID int32) (*models.Group, error) {
	db := g.db.WithContext(ctx)
	if !g.Exists(ctx, ID) {
		return nil, errors.New("group not found")
	}
	var grp models.Group
	err := db.Model(&grp).Where("\"group\".id = ?", ID).Relation("Organization").Select()
	return &grp, err
}

func (g *pgGroupRepository) Exists(ctx context.Context, ID int32) bool {
	db := g.db.WithContext(ctx)
	var num int32
	_, err := db.Query(pg.Scan(&num), "SELECT id from groups where id = ?", ID)
	if err != nil {
		panic(err)
	}
	return num == ID
}

func (g *pgGroupRepository) FindByIdAndOrganizationId(ctx context.Context, Id int32, Oid int32) (*models.Group, error) {
	db := g.db.WithContext(ctx)
	var grp models.Group
	err := db.Model(&grp).
		Where("\"group\".id = ?", Id).
		Where("\"group\".organization_id = ?", Oid).
		Relation("Organization").Select()
	return &grp, err
}

func (g *pgGroupRepository) Update(ctx context.Context, tx *pg.Tx, group *models.Group) error {
	err := tx.Update(group)
	return err
}

func (g *pgGroupRepository) Delete(ctx context.Context, group *models.Group) error {
	db := g.db.WithContext(ctx)
	err := db.Delete(group)
	return err
}

func (g *pgGroupRepository) FindAllByIdIn(ctx context.Context, ids []int32) []*models.Group {
	db := g.db.WithContext(ctx)
	var groups []*models.Group
	_ = db.Model(&groups). // TODO err handling
		Where("id in (?)", pg.In(ids)).
		Select()
	return groups
}
