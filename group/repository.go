package group

import (
	"context"
	"github.com/go-pg/pg/v9"
	"github.com/imtanmoy/authz/models"
)

type Repository interface {
	Fetch(ctx context.Context, organizationId int32) ([]*models.Group, error)
	Store(ctx context.Context, group *models.Group) (*models.Group, error)
	FindByName(ctx context.Context, organization *models.Organization, name string) (*models.Group, error)
	Find(ctx context.Context, ID int32) (*models.Group, error)
	Exists(ctx context.Context, ID int32) bool
	FindByIdAndOrganizationId(ctx context.Context, Id int32, Oid int32) (*models.Group, error)
	Delete(ctx context.Context, group *models.Group) error
	Update(ctx context.Context, tx *pg.Tx, group *models.Group) error
	FindAllByIdIn(ctx context.Context, ids []int32) []*models.Group
}
//
//type groupRepository struct {
//	db *pg.DB
//}
//
//var _ Repository = (*groupRepository)(nil)
//
//func NewGroupRepository(db *pg.DB) Repository {
//	return &groupRepository{
//		db,
//	}
//}
//
//func (g *groupRepository) List(ctx context.Context, organizationId int32) ([]*models.Group, error) {
//	var groups []*models.Group
//	err := g.db.Model(&groups).Where("organization_id = ?", organizationId).Relation("Organization").Select()
//	return groups, err
//}
//
//func (g *groupRepository) Create(ctx context.Context, group *models.Group) (*models.Group, error) {
//	_, err := g.db.Model(group).Returning("*").Insert()
//	return group, err
//}
//
//func (g *groupRepository) FindByName(ctx context.Context, organization *models.Organization, name string) (*models.Group, error) {
//	var group models.Group
//	err := g.db.Model(&group).
//		Where("name = ?", name).
//		Where("organization_id = ?", organization.ID).
//		First()
//	return &group, err
//}
//
//func (g *groupRepository) Find(ctx context.Context, ID int32) (*models.Group, error) {
//	if !g.Exists(ctx, ID) {
//		return nil, errors.New("group not found")
//	}
//	var group models.Group
//	err := g.db.Model(&group).Where("\"group\".id = ?", ID).Relation("Organization").Select()
//	return &group, err
//}
//
//func (g *groupRepository) Exists(ctx context.Context, ID int32) bool {
//	var num int32
//	_, err := g.db.Query(pg.Scan(&num), "SELECT id from groups where id = ?", ID)
//	if err != nil {
//		panic(err)
//	}
//	return num == ID
//}
//
//func (g *groupRepository) FindByIdAndOrganizationId(ctx context.Context, Id int32, Oid int32) (*models.Group, error) {
//	var group models.Group
//	err := g.db.Model(&group).
//		Where("\"group\".id = ?", Id).
//		Where("\"group\".organization_id = ?", Oid).
//		Relation("Organization").Select()
//	return &group, err
//}
//
//func (g *groupRepository) Update(ctx context.Context, tx *pg.Tx, group *models.Group) error {
//	err := tx.Update(group)
//	return err
//}
//
//func (g *groupRepository) Delete(ctx context.Context, group *models.Group) error {
//	err := g.db.Delete(group)
//	return err
//}
//
//func (g *groupRepository) FindAllByIdIn(ctx context.Context, ids []int32) []*models.Group {
//	var groups []*models.Group
//	_ = g.db.Model(&groups). // TODO err handling
//		Where("id in (?)", pg.In(ids)).
//		Select()
//	return groups
//}
