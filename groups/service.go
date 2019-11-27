package groups

import (
	"fmt"
	"github.com/go-pg/pg/v9"
	"github.com/imtanmoy/authz/casbin"
	"github.com/imtanmoy/authz/models"
	"github.com/imtanmoy/authz/permissions"
	"github.com/imtanmoy/authz/users"
)

type Service interface {
	List(organization *models.Organization) ([]*models.Group, error)
	Create(groupPayload *GroupPayload, organization *models.Organization, users []*models.User, permissions []*models.Permission) (*models.Group, error)
	FindByName(organization *models.Organization, name string) (*models.Group, error)
	Find(ID int32) (*models.Group, error)
	Exists(ID int32) bool
	FindByIdAndOrganizationId(Id int32, Oid int32) (*models.Group, error)
	Update(group *models.Group, users []*models.User, permissions []*models.Permission) error
}

type groupService struct {
	db                   *pg.DB
	repository           Repository
	userRepository       users.Repository
	permissionRepository permissions.Repository
}

var _ Service = (*groupService)(nil)

func NewGroupService(db *pg.DB) Service {
	return &groupService{
		db:                   db,
		repository:           NewGroupRepository(db),
		userRepository:       users.NewUserRepository(db),
		permissionRepository: permissions.NewPermissionRepository(db),
	}
}

func (g *groupService) List(organization *models.Organization) ([]*models.Group, error) {
	return g.repository.List(organization)
}

func (g *groupService) Create(
	groupPayload *GroupPayload,
	organization *models.Organization,
	users []*models.User,
	permissions []*models.Permission,
) (*models.Group, error) {
	var group models.Group

	tx, err := g.db.Begin()
	if err != nil {
		return nil, err
	}

	group.Name = groupPayload.Name
	group.Organization = organization
	group.OrganizationID = organization.ID

	newGroup, err := g.repository.Create(tx, &group)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}
	err = tx.Commit()
	groupID := fmt.Sprintf("group::%d", newGroup.ID)

	// add permissions for group
	for _, permission := range permissions {
		permissionID := fmt.Sprintf("permission::%d", permission.ID)
		_, err = casbin.Enforcer.AddPermissionForUser(groupID, permissionID, permission.Action)
		if err != nil {
			return nil, err
		}
	}
	for _, user := range users {
		userID := fmt.Sprintf("user::%d", user.ID)
		_, err = casbin.Enforcer.AddRoleForUser(userID, groupID)
		if err != nil {
			return nil, err
		}
	}
	return newGroup, err
}

func (g *groupService) FindByName(organization *models.Organization, name string) (*models.Group, error) {
	return g.repository.FindByName(organization, name)
}

func (g *groupService) Find(ID int32) (*models.Group, error) {
	return g.repository.Find(ID)
}

func (g *groupService) Exists(ID int32) bool {
	return g.repository.Exists(ID)
}

func (g *groupService) FindByIdAndOrganizationId(Id int32, Oid int32) (*models.Group, error) {
	return g.repository.FindByIdAndOrganizationId(Id, Oid)
}

func (g *groupService) Update(group *models.Group, users []*models.User, permissions []*models.Permission) error {
	tx, err := g.db.Begin()
	if err != nil {
		return err
	}

	err = g.repository.Update(tx, group)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	groupID := fmt.Sprintf("group::%d", group.ID)
	gpermissions := casbin.Enforcer.GetPermissionsForUser(groupID)
	existingPermissions := make([]int32, 0)
	for _, permission := range gpermissions {
		existingPermissions = append(existingPermissions, GetIntID(permission[1]))
	}
	newPermissions := make([]int32, 0)
	for _, permission := range permissions {
		newPermissions = append(newPermissions, permission.ID)
	}
	oldPermissions := Intersection(existingPermissions, newPermissions)
	deletePermissions := Minus(existingPermissions, oldPermissions)

	//create new permission with newPermissions
	for _, permission := range permissions {
		if Exists(newPermissions, permission.ID) {
			permissionID := fmt.Sprintf("permission::%d", permission.ID)
			params := []string{groupID, permissionID, permission.Action}
			_, err = casbin.Enforcer.AddPolicy(params)
			if err != nil {
				return err
			}
		}
	}

	//delete permissions with deletePermissions
	deletePermissionModels := g.permissionRepository.FindAllByIdIn(deletePermissions)
	for _, permission := range deletePermissionModels {
		permissionID := fmt.Sprintf("permission::%d", permission.ID)
		params := []string{groupID, permissionID, permission.Action}
		_, err = casbin.Enforcer.RemovePolicy(params)
		if err != nil {
			return err
		}
	}

	group.Permissions = permissions

	gUsers, err := casbin.Enforcer.GetUsersForRole(groupID)
	if err != nil {
		return err
	}
	// group user update
	existingUsers := make([]int32, 0)
	for _, user := range gUsers {
		existingUsers = append(existingUsers, GetIntID(user))
	}
	newUsers := make([]int32, 0)
	for _, user := range users {
		newUsers = append(newUsers, user.ID)
	}
	oldUsers := Intersection(existingUsers, newUsers)
	deleteUsers := Minus(existingUsers, oldUsers)

	//add new users to group
	for _, user := range users {
		if Exists(newUsers, user.ID) {
			userID := fmt.Sprintf("user::%d", user.ID)
			_, err = casbin.Enforcer.AddRoleForUser(userID, groupID)
			if err != nil {
				return err
			}
		}
	}

	//delete users from group
	deleteUsersModels := g.userRepository.FindAllByIdIn(deleteUsers)
	for _, user := range deleteUsersModels {
		fmt.Println(user)
		userID := fmt.Sprintf("user::%d", user.ID)
		_, err = casbin.Enforcer.DeleteRoleForUser(userID, groupID)
		if err != nil {
			return err
		}
	}

	group.Users = users

	return nil
}
