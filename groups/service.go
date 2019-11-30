package groups

import (
	"fmt"
	"github.com/go-pg/pg/v9"
	"github.com/imtanmoy/authz/casbin"
	"github.com/imtanmoy/authz/models"
	"github.com/imtanmoy/authz/permissions"
	"github.com/imtanmoy/authz/users"
	"github.com/imtanmoy/authz/utils"
)

type Service interface {
	List(organization *models.Organization) ([]*models.Group, error)
	Create(groupPayload *GroupPayload, organization *models.Organization, users []*models.User, permissions []*models.Permission) (*models.Group, error)
	FindByName(organization *models.Organization, name string) (*models.Group, error)
	Find(ID int32) (*models.Group, error)
	Exists(ID int32) bool
	FindByIdAndOrganizationId(Id int32, Oid int32) (*models.Group, error)
	Update(group *models.Group, users []*models.User, permissions []*models.Permission) error
	Delete(group *models.Group) error
}

type groupService struct {
	db                   *pg.DB
	repository           Repository
	userRepository       users.Repository
	permissionRepository permissions.Repository
	casbinService        casbin.Service
}

var _ Service = (*groupService)(nil)

func NewGroupService(db *pg.DB) Service {
	return &groupService{
		db:                   db,
		repository:           NewGroupRepository(db),
		userRepository:       users.NewUserRepository(db),
		permissionRepository: permissions.NewPermissionRepository(db),
		casbinService:        casbin.NewCasbinService(db),
	}
}

func (g *groupService) List(organization *models.Organization) ([]*models.Group, error) {
	groups, err := g.repository.List(organization.ID)
	if err != nil {
		return nil, err
	}
	for _, group := range groups {
		userList, err := g.casbinService.GetUsersForGroup(group.ID)
		if err != nil {
			return nil, err
		}
		group.Users = userList
	}

	for _, group := range groups {
		userList, err := g.casbinService.GetUsersForGroup(group.ID)
		if err != nil {
			return nil, err
		}
		group.Users = userList
	}
	return groups, nil
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
	// add permissions for group
	err = g.casbinService.AddPermissionsForGroup(group.ID, permissions)
	if err != nil {
		return nil, err
	}

	// add users for group
	err = g.casbinService.AddUsersForGroup(group.ID, users)
	if err != nil {
		return nil, err
	}
	newGroup.Users = users
	newGroup.Permissions = permissions
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
		existingPermissions = append(existingPermissions, utils.GetIntID(permission[1]))
	}
	newPermissions := make([]int32, 0)
	for _, permission := range permissions {
		newPermissions = append(newPermissions, permission.ID)
	}
	oldPermissions := utils.Intersection(existingPermissions, newPermissions)
	deletePermissions := utils.Minus(existingPermissions, oldPermissions)

	//create new permission with newPermissions
	for _, permission := range permissions {
		if utils.Exists(newPermissions, permission.ID) {
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
		existingUsers = append(existingUsers, utils.GetIntID(user))
	}
	newUsers := make([]int32, 0)
	for _, user := range users {
		newUsers = append(newUsers, user.ID)
	}
	oldUsers := utils.Intersection(existingUsers, newUsers)
	deleteUsers := utils.Minus(existingUsers, oldUsers)

	//add new users to group
	for _, user := range users {
		if utils.Exists(newUsers, user.ID) {
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

func (g *groupService) Delete(group *models.Group) error {
	//groupID := fmt.Sprintf("group::%d", group.ID)
	//gpermissions := casbin.Enforcer.GetPermissionsForUser(groupID)
	return g.repository.Delete(group)
}
