package users

import (
	"errors"
	"github.com/go-pg/pg/v9"
	"github.com/imtanmoy/authz/models"
)

type Repository interface {
	List() ([]*models.User, error)
	Find(ID int32) (*models.User, error)
	Create(tx *pg.Tx, user *models.User) (*models.User, error)
	FirstOrCreate(tx *pg.Tx, user *models.User) (*models.User, error)
	Update(tx *pg.Tx, user *models.User) (*models.User, error)
	Delete(tx *pg.Tx, user *models.User) error
	Exists(ID int32) bool
}

type userRepository struct {
	db *pg.DB
}

func NewUserRepository(db *pg.DB) Repository {
	return &userRepository{
		db,
	}
}

func (u *userRepository) List() ([]*models.User, error) {
	var users []*models.User
	err := u.db.Model(&users).Relation("Organization").Relation("Groups").Select()
	return users, err
}

func (u *userRepository) Find(ID int32) (*models.User, error) {
	if u.Exists(ID) {
		return nil, errors.New("user does not exists")
	}
	user := new(models.User)
	err := u.db.Model(user).Where("id = ?", ID).Relation("Organization").Relation("Groups").Select()
	return user, err
}

func (u *userRepository) Create(tx *pg.Tx, user *models.User) (*models.User, error) {
	_, err := u.db.Model(user).Returning("*").Insert()
	return user, err
}

func (u *userRepository) FirstOrCreate(tx *pg.Tx, user *models.User) (*models.User, error) {
	panic("implement me")
}

func (u *userRepository) Update(tx *pg.Tx, user *models.User) (*models.User, error) {
	err := u.db.Update(user)
	return user, err
}

func (u *userRepository) Delete(tx *pg.Tx, user *models.User) error {
	err := u.db.Delete(user)
	return err
}

func (u *userRepository) Exists(ID int32) bool {
	var num int32
	_, err := u.db.Query(pg.Scan(&num), "SELECT id from users where id = ?", ID)
	if err != nil {
		panic(err)
	}
	return num == ID
}

var _ Repository = (*userRepository)(nil)
