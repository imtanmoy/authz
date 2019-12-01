package authorizer

import (
	"errors"
	"fmt"
	casbinerros "github.com/casbin/casbin/v2/errors"
	"github.com/go-pg/pg/v9"
	"github.com/imtanmoy/authz/groups"
	"github.com/imtanmoy/authz/models"
	"github.com/imtanmoy/authz/permissions"
	"github.com/imtanmoy/authz/users"
	"github.com/imtanmoy/authz/utils"
)

type Service interface {
	AddPermissionsForGroup(id int32, permissions []*models.Permission) error
	GetPermissionsForGroup(id int32) ([]*models.Permission, error)
	RemovePermissionsForGroup(id int32, permissions []*models.Permission) error

	AddUsersForGroup(id int32, users []*models.User) error
	GetUsersForGroup(id int32) ([]*models.User, error)
	RemoveUsersForGroup(id int32, users []*models.User) error

	DeleteGroup(id int32) error

	AddPermissionsForUser(id int32, permissions []*models.Permission) error
	GetPermissionsForUser(id int32) ([]*models.Permission, error)
	RemovePermissionsForUser(id int32, permissions []*models.Permission) error

	AddGroupsForUser(id int32, groups []*models.Group) error
	GetGroupsForUser(id int32) ([]*models.Group, error)
	RemoveGroupsForUser(id int32, groups []*models.Group) error
}

type authorizerService struct {
	db                   *pg.DB
	userRepository       users.Repository
	permissionRepository permissions.Repository
	groupRepository      groups.Repository
}

var _ Service = (*authorizerService)(nil)

func NewAuthorizerService(db *pg.DB) Service {
	return &authorizerService{
		db:                   db,
		userRepository:       users.NewUserRepository(db),
		permissionRepository: permissions.NewPermissionRepository(db),
		groupRepository:      groups.NewGroupRepository(db),
	}
}

func (c *authorizerService) AddPermissionsForGroup(id int32, permissions []*models.Permission) error {
	groupId := fmt.Sprintf("group::%d", id)
	for _, permission := range permissions {
		permissionID := fmt.Sprintf("permission::%d", permission.ID)
		_, err := enforcer.AddPermissionForUser(groupId, permissionID, permission.Action)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *authorizerService) GetPermissionsForGroup(id int32) ([]*models.Permission, error) {
	groupId := fmt.Sprintf("group::%d", id)

	permissionList, err := enforcer.GetImplicitPermissionsForUser(groupId)
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

func (c *authorizerService) RemovePermissionsForGroup(id int32, permissions []*models.Permission) error {
	groupId := fmt.Sprintf("group::%d", id)
	for _, permission := range permissions {
		permissionID := fmt.Sprintf("permission::%d", permission.ID)
		_, err := enforcer.DeletePermissionForUser(groupId, permissionID, permission.Action)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *authorizerService) AddUsersForGroup(id int32, users []*models.User) error {
	groupId := fmt.Sprintf("group::%d", id)
	for _, user := range users {
		userID := fmt.Sprintf("user::%d", user.ID)
		_, err := enforcer.AddRoleForUser(userID, groupId)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *authorizerService) GetUsersForGroup(id int32) ([]*models.User, error) {
	groupId := fmt.Sprintf("group::%d", id)

	userList, err := enforcer.GetUsersForRole(groupId)
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

func (c *authorizerService) RemoveUsersForGroup(id int32, users []*models.User) error {
	groupId := fmt.Sprintf("group::%d", id)
	for _, user := range users {
		userID := fmt.Sprintf("user::%d", user.ID)
		_, err := enforcer.DeleteRoleForUser(userID, groupId)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *authorizerService) DeleteGroup(id int32) error {
	groupId := fmt.Sprintf("group::%d", id)
	_, err := enforcer.DeleteRole(groupId)
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
	userId := fmt.Sprintf("user::%d", id)
	groupList, err := enforcer.GetRolesForUser(userId)
	if errors.Is(err, casbinerros.ERR_NAME_NOT_FOUND) {
		return make([]*models.Group, 0), nil
	}
	if err != nil {
		return nil, err
	}
	var gIds []int32
	for _, group := range groupList {
		gIds = append(gIds, utils.GetIntID(group))
	}
	return c.groupRepository.FindAllByIdIn(gIds), nil
}

func (c *authorizerService) RemoveGroupsForUser(id int32, groups []*models.Group) error {
	panic("implement me")
}
