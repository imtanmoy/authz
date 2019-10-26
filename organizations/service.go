package organizations

import (
	"github.com/go-pg/pg/v9"
	"github.com/imtanmoy/authz/models"
)

type Service interface {
	List() ([]*models.Organization, error)
	Find(ID int) (*models.Organization, error)
	Create(organization *models.Organization) (*models.Organization, error)
	FirstOrCreate(organization *models.Organization) (*models.Organization, error)
	Update(organization *models.Organization) (*models.Organization, error)
	Delete(organization *models.Organization) error
}

type organizationService struct {
	db         *pg.DB
	repository Repository
}

func NewOrganizationService(db *pg.DB) Service {
	return &organizationService{
		repository: NewOrganizationRepository(db),
		db:         db,
	}
}

func (o *organizationService) List() ([]*models.Organization, error) {
	return o.repository.List()
}

func (o *organizationService) Find(ID int) (*models.Organization, error) {
	return o.repository.Find(ID)
}

func (o *organizationService) Create(organization *models.Organization) (*models.Organization, error) {
	tx, _ := o.db.Begin()
	defer tx.Commit()
	return o.repository.Create(tx, organization)
}

func (o *organizationService) FirstOrCreate(organization *models.Organization) (*models.Organization, error) {
	panic("implement me")
}

func (o *organizationService) Update(organization *models.Organization) (*models.Organization, error) {
	tx, _ := o.db.Begin()
	defer tx.Commit()
	return o.repository.Update(tx, organization)
}

func (o *organizationService) Delete(organization *models.Organization) error {
	tx, _ := o.db.Begin()
	defer tx.Commit()
	return o.repository.Delete(tx, organization)
}

var _ Service = (*organizationService)(nil)
