package permissions

import (
	"github.com/go-pg/pg/v9"
	"github.com/imtanmoy/authz/models"
)

type Service interface {
	FindAllByIdIn(ids []int32) []*models.Permission
}

type permissionService struct {
	db         *pg.DB
	repository Repository
}

var _ Service = (*permissionService)(nil)

func NewPermissionService(db *pg.DB) Service {
	return &permissionService{
		repository: NewPermissionRepository(db),
		db:         db,
	}
}

func (p *permissionService) FindAllByIdIn(ids []int32) []*models.Permission {
	return p.repository.FindAllByIdIn(ids)
}
