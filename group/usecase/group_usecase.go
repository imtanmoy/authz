package usecase

import (
	"context"
	"github.com/imtanmoy/authz/authorizer"
	"github.com/imtanmoy/authz/group"
	"github.com/imtanmoy/authz/models"
	"golang.org/x/sync/errgroup"
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

// populates each group with users
func (g *groupUsecase) populateGroupsWithUsers(ctx context.Context, groups []*models.Group) error {
	errs, _ := errgroup.WithContext(ctx)
	for _, grp := range groups {
		grpR := grp
		errs.Go(func() error {
			userList, err := g.authorizerService.GetUsersForGroup(grpR.ID)
			if err != nil {
				return err
			}
			grpR.Users = userList
			return nil
		})
	}
	return errs.Wait()
}

// populates each group with permissions
func (g *groupUsecase) populateGroupsWithPermissions(ctx context.Context, groups []*models.Group) error {
	errs, _ := errgroup.WithContext(ctx)
	for _, grp := range groups {
		grpR := grp
		errs.Go(func() error {
			permissionList, err := g.authorizerService.GetPermissionsForGroup(grpR.ID)
			if err != nil {
				return err
			}
			grpR.Permissions = permissionList
			return nil
		})
	}
	return errs.Wait()
}

func (g *groupUsecase) Fetch(ctx context.Context, organizationId int32) ([]*models.Group, error) {
	groups, err := g.groupRepo.Fetch(ctx, organizationId)
	if err != nil {
		return nil, err
	}

	err = g.populateGroupsWithUsers(ctx, groups)
	if err != nil {
		return nil, err
	}

	err = g.populateGroupsWithPermissions(ctx, groups)
	if err != nil {
		return nil, err
	}
	return groups, nil
}

func (g *groupUsecase) Store(ctx context.Context, gr *models.Group, users []*models.User, permissions []*models.Permission) error {
	grp, err := g.groupRepo.Store(ctx, gr)
	if err != nil {
		return err
	}
	// add permissions for group
	err = g.authorizerService.AddPermissionsForGroup(ctx, grp.ID, permissions)
	if err != nil {
		return err
	}
	// add users for group
	err = g.authorizerService.AddUsersForGroup(ctx, grp.ID, users)
	if err != nil {
		return err
	}
	grp.Users = users
	grp.Permissions = permissions
	return nil
}

func (g *groupUsecase) GetByID(ctx context.Context, id int32) (*models.Group, error) {
	panic("implement me")
}

func (g *groupUsecase) Update(ctx context.Context, gr *models.Group, users []*models.User, permissions []*models.Permission) error {
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
