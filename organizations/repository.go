package organizations

import (
	"github.com/go-pg/pg/v9"
	"github.com/imtanmoy/authz/models"
)

type Repository interface {
	List() ([]*models.Organization, error)
	Find(ID int) (*models.Organization, error)
	Create(tx *pg.Tx, organization *models.Organization) (*models.Organization, error)
	FirstOrCreate(tx *pg.Tx, organization *models.Organization) (*models.Organization, error)
	Update(tx *pg.Tx, organization *models.Organization) (*models.Organization, error)
	Delete(tx *pg.Tx, organization *models.Organization) error
}

type organizationRepository struct {
	db *pg.DB
}

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

func (o *organizationRepository) Find(ID int) (*models.Organization, error) {
	organization := new(models.Organization)
	err := o.db.Model(organization).Where("id = ?", ID).Relation("Users").Select()
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

var _ Repository = (*organizationRepository)(nil)
