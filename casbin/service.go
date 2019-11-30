package casbin

import (
	"errors"
	"fmt"
	casbinerros "github.com/casbin/casbin/v2/errors"
	"github.com/go-pg/pg/v9"
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
}

type casbinService struct {
	db                   *pg.DB
	userRepository       users.Repository
	permissionRepository permissions.Repository
}

var _ Service = (*casbinService)(nil)

func NewCasbinService(db *pg.DB) Service {
	return &casbinService{
		db:                   db,
		userRepository:       users.NewUserRepository(db),
		permissionRepository: permissions.NewPermissionRepository(db),
	}
}

func (c *casbinService) AddPermissionsForGroup(id int32, permissions []*models.Permission) error {
	groupId := fmt.Sprintf("group::%d", id)
	for _, permission := range permissions {
		permissionID := fmt.Sprintf("permission::%d", permission.ID)
		_, err := Enforcer.AddPermissionForUser(groupId, permissionID, permission.Action)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *casbinService) GetPermissionsForGroup(id int32) ([]*models.Permission, error) {
	groupId := fmt.Sprintf("group::%d", id)

	permissionList, err := Enforcer.GetImplicitPermissionsForUser(groupId)
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

func (c *casbinService) RemovePermissionsForGroup(id int32, permissions []*models.Permission) error {
	groupId := fmt.Sprintf("group::%d", id)
	for _, permission := range permissions {
		permissionID := fmt.Sprintf("permission::%d", permission.ID)
		_, err := Enforcer.DeletePermissionForUser(groupId, permissionID, permission.Action)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *casbinService) AddUsersForGroup(id int32, users []*models.User) error {
	groupId := fmt.Sprintf("group::%d", id)
	for _, user := range users {
		userID := fmt.Sprintf("user::%d", user.ID)
		_, err := Enforcer.AddRoleForUser(userID, groupId)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *casbinService) GetUsersForGroup(id int32) ([]*models.User, error) {
	groupId := fmt.Sprintf("group::%d", id)

	userList, err := Enforcer.GetUsersForRole(groupId)
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

func (c *casbinService) RemoveUsersForGroup(id int32, users []*models.User) error {
	groupId := fmt.Sprintf("group::%d", id)
	for _, user := range users {
		userID := fmt.Sprintf("user::%d", user.ID)
		_, err := Enforcer.DeleteRoleForUser(userID, groupId)
		if err != nil {
			return err
		}
	}
	return nil
}
