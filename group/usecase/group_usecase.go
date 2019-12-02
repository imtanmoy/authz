package usecase

import (
	"context"
	"github.com/imtanmoy/authz/authorizer"
	"github.com/imtanmoy/authz/group"
	"github.com/imtanmoy/authz/models"
	"time"
)

type groupUsecase struct {
	groupRepo         group.Repository
	contextTimeout    time.Duration
	authorizerService authorizer.Service
}

var _ group.Usecase = (*groupUsecase)(nil)

// NewGroupUsecase will create new an groupUsecase object representation of group.Usecase interface
func NewGroupUsecase(g group.Repository, timeout time.Duration, authorizerService authorizer.Service) group.Usecase {
	return &groupUsecase{
		groupRepo:         g,
		contextTimeout:    timeout,
		authorizerService: authorizerService,
	}
}

func (g *groupUsecase) Fetch(ctx context.Context, organizationId int32) ([]*models.Group, error) {
	groups, err := g.groupRepo.Fetch(ctx, organizationId)
	if err != nil {
		return nil, err
	}
	for _, gr := range groups {
		userList, err := g.authorizerService.GetUsersForGroup(gr.ID)
		if err != nil {
			return nil, err
		}
		gr.Users = userList
	}

	for _, gr := range groups {
		permissionList, err := g.authorizerService.GetPermissionsForGroup(gr.ID)
		if err != nil {
			return nil, err
		}
		gr.Permissions = permissionList
	}
	return groups, nil
}

func (g *groupUsecase) GetByID(ctx context.Context, id int32) (*models.Group, error) {
	panic("implement me")
}

func (g *groupUsecase) Update(ctx context.Context, gr *models.Group, users []*models.User, permissions []*models.Permission) error {
	panic("implement me")
}

func (g *groupUsecase) Store(ctx context.Context, gr *models.Group, users []*models.User, permissions []*models.Permission) error {
	panic("implement me")
}

func (g *groupUsecase) Delete(ctx context.Context, id int32) error {
	panic("implement me")
}

func (g *groupUsecase) Exists(ctx context.Context, ID int32) bool {
	panic("implement me")
}

func (g *groupUsecase) FindByName(ctx context.Context, organization *models.Organization, name string) (*models.Group, error) {
	return g.groupRepo.FindByName(ctx, organization, name)
}

func (g *groupUsecase) FindByIdAndOrganizationId(ctx context.Context, Id int32, Oid int32) (*models.Group, error) {
	panic("implement me")
}
