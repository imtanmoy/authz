package group

import (
	"context"
	"github.com/imtanmoy/authz/models"
)

// Usecase represent the groups's use cases
type Usecase interface {
	Fetch(ctx context.Context) ([]*models.Group, error)
	GetByID(ctx context.Context, id int64) (*models.Group, error)
	Update(ctx context.Context, gr *models.Group, users []*models.User, permissions []*models.Permission) error
	Store(ctx context.Context, gr *models.Group, users []*models.User, permissions []*models.Permission) error
	Delete(ctx context.Context, id int64) error
	Exists(ctx context.Context, ID int32) bool
	FindByName(ctx context.Context, organization *models.Organization, name string) (*models.Group, error)
	FindByIdAndOrganizationId(ctx context.Context, Id int32, Oid int32) (*models.Group, error)
}
