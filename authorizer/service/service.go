package service

import (
	"context"
	"errors"
	"fmt"
	casbinerros "github.com/casbin/casbin/v2/errors"
	"github.com/go-pg/pg/v9"
	"github.com/imtanmoy/authz/authorizer"
	"github.com/imtanmoy/authz/models"
	"github.com/imtanmoy/authz/permissions"
	"github.com/imtanmoy/authz/users"
	"github.com/imtanmoy/authz/utils"
	"golang.org/x/sync/errgroup"
)

type authorizerService struct {
	db                   *pg.DB
	userRepository       users.Repository
	permissionRepository permissions.Repository
}

var _ authorizer.Service = (*authorizerService)(nil)

func NewAuthorizerService(db *pg.DB) authorizer.Service {
	return &authorizerService{
		db:                   db,
		userRepository:       users.NewUserRepository(db),
		permissionRepository: permissions.NewPermissionRepository(db),
	}
}

func (c *authorizerService) AddPermissionsForGroup(ctx context.Context, id int32, permissions []*models.Permission) error {
	errs, _ := errgroup.WithContext(ctx)
	groupId := fmt.Sprintf("group::%d", id)
	for _, permission := range permissions {
		p := permission
		errs.Go(func() error {
			permissionID := fmt.Sprintf("permission::%d", p.ID)
			pAction := p.Action
			_, err := authorizer.Enforcer.AddPermissionForUser(groupId, permissionID, pAction)
			if err != nil {
				return err
			}
			return nil
		})
	}
	return errs.Wait()
}

func (c *authorizerService) GetPermissionsForGroup(id int32) ([]*models.Permission, error) {
	groupId := fmt.Sprintf("group::%d", id)

	permissionList, err := authorizer.Enforcer.GetImplicitPermissionsForUser(groupId)
	if errors.Is(err, casbinerros.ERR_NAME_NOT_FOUND) {
		return make([]*models.Permission, 0), nil
	}
	if err != nil {
		return nil, err
	}

	var pIds []int32
	for _, p := range permissionList {
		pIds = append(pIds, utils.GetIntID(p[1]))
	}
	return c.permissionRepository.FindAllByIdIn(pIds), nil
}

func (c *authorizerService) RemovePermissionsForGroup(ctx context.Context, id int32, permissions []*models.Permission) error {
	errs, _ := errgroup.WithContext(ctx)
	groupId := fmt.Sprintf("group::%d", id)
	for _, permission := range permissions {
		p := permission
		errs.Go(func() error {
			permissionID := fmt.Sprintf("permission::%d", p.ID)
			pAction := p.Action
			_, err := authorizer.Enforcer.DeletePermissionForUser(groupId, permissionID, pAction)
			if err != nil {
				return err
			}
			return nil
		})
	}
	return errs.Wait()
}

func (c *authorizerService) AddUsersForGroup(ctx context.Context, id int32, users []*models.User) error {
	errs, _ := errgroup.WithContext(ctx)
	groupId := fmt.Sprintf("group::%d", id)
	for _, u := range users {
		useR := u
		errs.Go(func() error {
			userID := fmt.Sprintf("user::%d", useR.ID)
			_, err := authorizer.Enforcer.AddRoleForUser(userID, groupId)
			if err != nil {
				return err
			}
			return nil
		})
	}
	return errs.Wait()
}

func (c *authorizerService) GetUsersForGroup(id int32) ([]*models.User, error) {
	groupId := fmt.Sprintf("group::%d", id)

	userList, err := authorizer.Enforcer.GetUsersForRole(groupId)
	if errors.Is(err, casbinerros.ERR_NAME_NOT_FOUND) {
		return make([]*models.User, 0), nil
	}
	if err != nil {
		return nil, err
	}
	var uIds []int32
	for _, user := range userList {
		uIds = append(uIds, utils.GetIntID(user))
	}
	return c.userRepository.FindAllByIdIn(uIds), nil
}

func (c *authorizerService) RemoveUsersForGroup(ctx context.Context, id int32, users []*models.User) error {
	errs, _ := errgroup.WithContext(ctx)
	groupId := fmt.Sprintf("group::%d", id)
	for _, user := range users {
		u := user
		errs.Go(func() error {
			userID := fmt.Sprintf("user::%d", u.ID)
			_, err := authorizer.Enforcer.DeleteRoleForUser(userID, groupId)
			if err != nil {
				return err
			}
			return nil
		})
	}
	return errs.Wait()
}

func (c *authorizerService) DeleteGroup(id int32) error {
	groupId := fmt.Sprintf("group::%d", id)
	_, err := authorizer.Enforcer.DeleteRole(groupId)
	return err
}

func (c *authorizerService) AddPermissionsForUser(id int32, permissions []*models.Permission) error {
	panic("implement me")
}

func (c *authorizerService) GetPermissionsForUser(id int32) ([]*models.Permission, error) {
	panic("implement me")
}

func (c *authorizerService) RemovePermissionsForUser(id int32, permissions []*models.Permission) error {
	panic("implement me")
}

func (c *authorizerService) AddGroupsForUser(id int32, groups []*models.Group) error {
	panic("implement me")
}

func (c *authorizerService) GetGroupsForUser(id int32) ([]*models.Group, error) {
	panic("implement me")
}

func (c *authorizerService) RemoveGroupsForUser(id int32, groups []*models.Group) error {
	panic("implement me")
}
