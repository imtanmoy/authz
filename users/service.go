package users

import (
	"github.com/go-pg/pg/v9"
	"github.com/imtanmoy/authz/models"
)

type Service interface {
	List() ([]*models.User, error)
	Find(ID int32) (*models.User, error)
	Create(organization *models.User) (*models.User, error)
	FirstOrCreate(organization *models.User) (*models.User, error)
	Update(organization *models.User) (*models.User, error)
	Delete(organization *models.User) error
	Exists(ID int32) bool
	FindAllByIdIn(ids []int32) []*models.User

	GetGroups(user *models.User) ([]*models.Group, error)
	GetPermissions(user *models.User) ([]*models.Permission, error)
}

type userService struct {
	db                *pg.DB
	repository        Repository
}

var _ Service = (*userService)(nil)

func NewUserService(db *pg.DB) Service {
	return &userService{
		repository: NewUserRepository(db),
		db:         db,
	}
}

func (u *userService) List() ([]*models.User, error) {
	return u.repository.List()
}

func (u *userService) Exists(ID int32) bool {
	return u.repository.Exists(ID)
}

func (u *userService) Find(ID int32) (*models.User, error) {
	return u.repository.Find(ID)
}

func (u *userService) Create(user *models.User) (*models.User, error) {
	tx, _ := u.db.Begin()
	defer tx.Commit()
	return u.repository.Create(tx, user)
}

func (u *userService) FirstOrCreate(organization *models.User) (*models.User, error) {
	panic("implement me")
}

func (u *userService) Update(user *models.User) (*models.User, error) {
	tx, _ := u.db.Begin()
	defer tx.Commit()
	return u.repository.Update(tx, user)
}

func (u *userService) Delete(user *models.User) error {
	tx, _ := u.db.Begin()
	defer tx.Commit()
	return u.repository.Delete(tx, user)
}

func (u *userService) FindAllByIdIn(ids []int32) []*models.User {
	return u.repository.FindAllByIdIn(ids)
}

func (u *userService) GetGroups(user *models.User) ([]*models.Group, error) {
	panic("implement me")
}

func (u *userService) GetPermissions(user *models.User) ([]*models.Permission, error) {
	panic("implement me")
}
