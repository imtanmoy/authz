package organizations

import (
	"errors"
	"github.com/go-pg/pg/v9"
	"github.com/imtanmoy/authz/models"
)

type Repository interface {
	List() ([]*models.Organization, error)
	Find(id int32) (*models.Organization, error)
	Create(tx *pg.Tx, organization *models.Organization) (*models.Organization, error)
	FirstOrCreate(tx *pg.Tx, organization *models.Organization) (*models.Organization, error)
	Update(tx *pg.Tx, organization *models.Organization) (*models.Organization, error)
	Delete(tx *pg.Tx, organization *models.Organization) error
	Exists(ID int32) bool
	FindUsersByIds(organization *models.Organization, ids []int32) ([]*models.User, error)
	FindPermissionsByIds(organization *models.Organization, ids []int32) ([]*models.Permission, error)
}

type organizationRepository struct {
	db *pg.DB
}

var _ Repository = (*organizationRepository)(nil)

func NewOrganizationRepository(db *pg.DB) Repository {
	return &organizationRepository{
		db,
	}
}

func (o *organizationRepository) List() ([]*models.Organization, error) {
	var organizations []*models.Organization
	err := o.db.Model(&organizations).Relation("Users").Select()
	return organizations, err
}

func (o *organizationRepository) Find(id int32) (*models.Organization, error) {
	if !o.Exists(id) {
		return nil, errors.New("organization does not exist")
	}
	organization := new(models.Organization)
	err := o.db.Model(organization).Where("id = ?", id).Relation("Users").Select()
	return organization, err
}

func (o *organizationRepository) Create(tx *pg.Tx, organization *models.Organization) (*models.Organization, error) {
	err := o.db.Insert(organization)
	return organization, err
	//tx.Insert(organization)
}

func (o *organizationRepository) FirstOrCreate(tx *pg.Tx, organization *models.Organization) (*models.Organization, error) {
	panic("implement me")
}

func (o *organizationRepository) Update(tx *pg.Tx, organization *models.Organization) (*models.Organization, error) {
	err := o.db.Update(organization)
	return organization, err
}

func (o *organizationRepository) Delete(tx *pg.Tx, organization *models.Organization) error {
	err := o.db.Delete(organization)
	return err
}

func (o *organizationRepository) Exists(id int32) bool {
	var num int32
	_, err := o.db.Query(pg.Scan(&num), "SELECT id from organizations where id = ?", id)
	if err != nil {
		panic(err)
	}
	return num == id
}

func (o *organizationRepository) FindUsersByIds(organization *models.Organization, ids []int32) ([]*models.User, error) {
	var users []*models.User
	err := o.db.Model(&users).
		Where("id in (?)", pg.In(ids)).
		Where("organization_id = ? ", organization.ID).
		Select()
	return users, err
}

func (o *organizationRepository) FindPermissionsByIds(organization *models.Organization, ids []int32) ([]*models.Permission, error) {
	var permissions []*models.Permission
	err := o.db.Model(&permissions).
		Where("id in (?)", pg.In(ids)).
		Where("organization_id = ? ", organization.ID).
		Select()
	return permissions, err
}
