package usecase

import (
	"context"
	"github.com/imtanmoy/authz/authorizer"
	"github.com/imtanmoy/authz/group"
	"github.com/imtanmoy/authz/models"
	"github.com/imtanmoy/authz/utils"
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
	err := g.groupRepo.Update(ctx, gr)
	if err != nil {
		return err
	}
	err = g.updatePermissionsForGroup(ctx, gr, permissions)
	if err != nil {
		return nil
	}
	err = g.updateUsersForGroup(ctx, gr, users)
	if err != nil {
		return nil
	}
	return nil
}

func (g *groupUsecase) Delete(ctx context.Context, gr *models.Group) error {
	err := g.authorizerService.DeleteGroup(gr.ID) // it will delete all permissions and users
	if err != nil {
		return err
	}
	return g.groupRepo.Delete(ctx, gr)
}

func (g *groupUsecase) Exists(ctx context.Context, ID int32) bool {
	panic("implement me")
}

func (g *groupUsecase) FindByNameAndOrganizationId(ctx context.Context, name string, oid int32) (*models.Group, error) {
	return g.groupRepo.FindByNameAndOrganizationId(ctx, name, oid)
}

func (g *groupUsecase) FindByIdAndOrganizationId(ctx context.Context, Id int32, Oid int32) (*models.Group, error) {
	grp, err := g.groupRepo.FindByIdAndOrganizationId(ctx, Id, Oid)
	if err != nil {
		return nil, err
	}

	// get user list
	userList, err := g.authorizerService.GetUsersForGroup(grp.ID)
	if err != nil {
		return nil, err
	}
	grp.Users = userList

	// get permission list
	permissionList, err := g.authorizerService.GetPermissionsForGroup(grp.ID)
	if err != nil {
		return nil, err
	}
	grp.Permissions = permissionList
	return grp, nil
}

// update a group's permissions
func (g *groupUsecase) updatePermissionsForGroup(ctx context.Context, grp *models.Group, permissions []*models.Permission) error {
	// getting groups existing permissions
	permissionList, err := g.authorizerService.GetPermissionsForGroup(grp.ID)
	if err != nil {
		return err
	}
	existingPermissions := make([]int32, 0)
	for _, permission := range permissionList {
		existingPermissions = append(existingPermissions, permission.ID)
	}
	newPermissions := make([]int32, 0)
	for _, permission := range permissions {
		newPermissions = append(newPermissions, permission.ID)
	}
	oldPermissions := utils.Intersection(existingPermissions, newPermissions)
	deletePermissions := utils.Minus(existingPermissions, oldPermissions)

	//create new permission with newPermissions
	willBeAddedPermissions := utils.Minus(newPermissions, oldPermissions)
	//willBeAddedPermissionModels := g.permissionRepository.FindAllByIdIn(willBeAddedPermissions)
	willBeAddedPermissionModels := getPermissionModels(willBeAddedPermissions, permissions)
	err = g.authorizerService.AddPermissionsForGroup(ctx, grp.ID, willBeAddedPermissionModels) //TODO dont send the model
	if err != nil {
		return err
	}

	//delete permissions with deletePermissions
	//deletePermissionModels := g.permissionRepository.FindAllByIdIn(deletePermissions)
	deletePermissionModels := getPermissionModels(deletePermissions, permissionList)
	err = g.authorizerService.RemovePermissionsForGroup(ctx, grp.ID, deletePermissionModels) //TODO dont send the model
	if err != nil {
		return err
	}
	grp.Permissions = permissions
	return nil
}

// update a group's users
func (g *groupUsecase) updateUsersForGroup(ctx context.Context, grp *models.Group, users []*models.User) error {
	// user update
	userList, err := g.authorizerService.GetUsersForGroup(grp.ID)
	if err != nil {
		return err
	}
	// group user update
	existingUsers := make([]int32, 0)
	for _, user := range userList {
		existingUsers = append(existingUsers, user.ID)
	}
	newUsers := make([]int32, 0)
	for _, user := range users {
		newUsers = append(newUsers, user.ID)
	}
	oldUsers := utils.Intersection(existingUsers, newUsers)
	deleteUsers := utils.Minus(existingUsers, oldUsers)

	// add users for group
	willBeAddedUsers := utils.Minus(newUsers, oldUsers)
	//willBeAddedUserModels := g.userRepository.FindAllByIdIn(willBeAddedUsers)
	willBeAddedUserModels := getUserModels(willBeAddedUsers, users)
	err = g.authorizerService.AddUsersForGroup(ctx, grp.ID, willBeAddedUserModels)
	if err != nil {
		return err
	}

	//delete users from group
	//deleteUsersModels := g.userRepository.FindAllByIdIn(deleteUsers)
	deleteUsersModels := getUserModels(deleteUsers, userList)
	err = g.authorizerService.RemoveUsersForGroup(ctx, grp.ID, deleteUsersModels)
	if err != nil {
		return err
	}
	grp.Users = users
	return nil
}

func getPermissionModels(ids []int32, permissions []*models.Permission) []*models.Permission {
	list := make([]*models.Permission, 0)
	for _, permission := range permissions {
		if utils.Exists(ids, permission.ID) {
			list = append(list, permission)
		}
	}
	return list
}

func getUserModels(ids []int32, users []*models.User) []*models.User {
	list := make([]*models.User, 0)
	for _, user := range users {
		if utils.Exists(ids, user.ID) {
			list = append(list, user)
		}
	}
	return list
}
