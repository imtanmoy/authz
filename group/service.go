package group

import (
	"github.com/go-pg/pg/v9"
	"github.com/imtanmoy/authz/authorizer"
	"github.com/imtanmoy/authz/models"
	"github.com/imtanmoy/authz/permissions"
	"github.com/imtanmoy/authz/users"
	"github.com/imtanmoy/authz/utils"
)

type Service interface {
	List(organization *models.Organization) ([]*models.Group, error)
	Create(groupPayload *GroupPayload, organization *models.Organization, users []*models.User, permissions []*models.Permission) (*models.Group, error)
	Find(ID int32) (*models.Group, error)
	Update(group *models.Group, users []*models.User, permissions []*models.Permission) error
	Delete(group *models.Group) error
	Exists(ID int32) bool
	FindByName(organization *models.Organization, name string) (*models.Group, error)
	FindByIdAndOrganizationId(Id int32, Oid int32) (*models.Group, error)
}

type groupService struct {
	db                   *pg.DB
	repository           Repository
	userRepository       users.Repository
	permissionRepository permissions.Repository
	authorizerService    authorizer.Service
}

var _ Service = (*groupService)(nil)

func NewGroupService(db *pg.DB) Service {
	return &groupService{
		db:                   db,
		repository:           NewGroupRepository(db),
		userRepository:       users.NewUserRepository(db),
		permissionRepository: permissions.NewPermissionRepository(db),
		authorizerService:    authorizer.NewAuthorizerService(db),
	}
}

func (g *groupService) List(organization *models.Organization) ([]*models.Group, error) {
	groups, err := g.repository.List(organization.ID)
	if err != nil {
		return nil, err
	}
	for _, group := range groups {
		userList, err := g.authorizerService.GetUsersForGroup(group.ID)
		if err != nil {
			return nil, err
		}
		group.Users = userList
	}

	for _, group := range groups {
		permissionList, err := g.authorizerService.GetPermissionsForGroup(group.ID)
		if err != nil {
			return nil, err
		}
		group.Permissions = permissionList
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
	err = g.authorizerService.AddPermissionsForGroup(group.ID, permissions)
	if err != nil {
		return nil, err
	}

	// add users for group
	err = g.authorizerService.AddUsersForGroup(group.ID, users)
	if err != nil {
		return nil, err
	}
	newGroup.Users = users
	newGroup.Permissions = permissions
	return newGroup, err
}

func (g *groupService) Find(ID int32) (*models.Group, error) {
	group, err := g.repository.Find(ID)
	if err != nil {
		return nil, err
	}
	// get user list
	userList, err := g.authorizerService.GetUsersForGroup(group.ID)
	if err != nil {
		return nil, err
	}
	group.Users = userList

	// get permission list
	permissionList, err := g.authorizerService.GetPermissionsForGroup(group.ID)
	if err != nil {
		return nil, err
	}
	group.Permissions = permissionList
	return group, nil
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

	// permission update
	permissionList, err := g.authorizerService.GetPermissionsForGroup(group.ID)
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
	err = g.authorizerService.AddPermissionsForGroup(group.ID, willBeAddedPermissionModels)
	if err != nil {
		return err
	}

	//delete permissions with deletePermissions
	//deletePermissionModels := g.permissionRepository.FindAllByIdIn(deletePermissions)
	deletePermissionModels := getPermissionModels(deletePermissions, permissionList)
	err = g.authorizerService.RemovePermissionsForGroup(group.ID, deletePermissionModels)
	if err != nil {
		return nil
	}
	group.Permissions = permissions

	// user update
	userList, err := g.authorizerService.GetUsersForGroup(group.ID)
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
	err = g.authorizerService.AddUsersForGroup(group.ID, willBeAddedUserModels)
	if err != nil {
		return err
	}

	//delete users from group
	//deleteUsersModels := g.userRepository.FindAllByIdIn(deleteUsers)
	deleteUsersModels := getUserModels(deleteUsers, userList)
	err = g.authorizerService.RemoveUsersForGroup(group.ID, deleteUsersModels)
	if err != nil {
		return err
	}

	group.Users = users

	return nil
}

func (g *groupService) Delete(group *models.Group) error {
	err := g.authorizerService.DeleteGroup(group.ID) // it will delete all permissions and users
	if err != nil {
		return err
	}
	return g.repository.Delete(group)
}

func (g *groupService) Exists(ID int32) bool {
	return g.repository.Exists(ID)
}

func (g *groupService) FindByName(organization *models.Organization, name string) (*models.Group, error) {
	return g.repository.FindByName(organization, name)
}

func (g *groupService) FindByIdAndOrganizationId(Id int32, Oid int32) (*models.Group, error) {
	group, err := g.repository.FindByIdAndOrganizationId(Id, Oid)
	if err != nil {
		return nil, err
	}

	// get user list
	userList, err := g.authorizerService.GetUsersForGroup(group.ID)
	if err != nil {
		return nil, err
	}
	group.Users = userList

	// get permission list
	permissionList, err := g.authorizerService.GetPermissionsForGroup(group.ID)
	if err != nil {
		return nil, err
	}
	group.Permissions = permissionList
	return group, nil
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
