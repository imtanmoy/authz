package authorizer

import (
	"context"
	"github.com/imtanmoy/authz/models"
)

type Service interface {
	AddPermissionsForGroup(ctx context.Context, id int32, permissions []*models.Permission) error
	GetPermissionsForGroup(id int32) ([]*models.Permission, error)
	RemovePermissionsForGroup(id int32, permissions []*models.Permission) error

	AddUsersForGroup(ctx context.Context, id int32, users []*models.User) error
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
