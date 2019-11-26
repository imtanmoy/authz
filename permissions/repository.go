package permissions

import (
	"github.com/go-pg/pg/v9"
	"github.com/imtanmoy/authz/models"
)

type Repository interface {
	FindAllByIdIn(ids []int32) []*models.Permission
}

type permissionRepository struct {
	db *pg.DB
}

var _ Repository = (*permissionRepository)(nil)

func NewPermissionRepository(db *pg.DB) Repository {
	return &permissionRepository{
		db,
	}
}

func (p *permissionRepository) FindAllByIdIn(ids []int32) []*models.Permission {
	var permissions []*models.Permission
	_ = p.db.Model(&permissions). // TODO err handling
		Where("id in (?)", pg.In(ids)).
		Select()
	return permissions
}
