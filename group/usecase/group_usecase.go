package usecase

import (
	"context"
	"github.com/imtanmoy/authz/group"
	"github.com/imtanmoy/authz/models"
	"time"
)

type groupUsecase struct {
	groupRepo      group.Repository
	contextTimeout time.Duration
}

var _ group.Usecase = (*groupUsecase)(nil)

// NewGroupUsecase will create new an groupUsecase object representation of group.Usecase interface
func NewGroupUsecase(g group.Repository, timeout time.Duration) group.Usecase {
	return &groupUsecase{
		groupRepo:      g,
		contextTimeout: timeout,
	}
}

func (g *groupUsecase) Fetch(ctx context.Context) ([]*models.Group, error) {
	panic("implement me")
}

func (g *groupUsecase) GetByID(ctx context.Context, id int64) (*models.Group, error) {
	panic("implement me")
}

func (g *groupUsecase) Update(ctx context.Context, gr *models.Group, users []*models.User, permissions []*models.Permission) error {
	panic("implement me")
}

func (g *groupUsecase) Store(ctx context.Context, gr *models.Group, users []*models.User, permissions []*models.Permission) error {
	panic("implement me")
}

func (g *groupUsecase) Delete(ctx context.Context, id int64) error {
	panic("implement me")
}

func (g *groupUsecase) Exists(ctx context.Context, ID int32) bool {
	panic("implement me")
}

func (g *groupUsecase) FindByName(ctx context.Context, organization *models.Organization, name string) (*models.Group, error) {
	panic("implement me")
}

func (g *groupUsecase) FindByIdAndOrganizationId(ctx context.Context, Id int32, Oid int32) (*models.Group, error) {
	panic("implement me")
}
