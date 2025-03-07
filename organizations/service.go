package organizations

import (
	"github.com/go-pg/pg/v9"
	"github.com/imtanmoy/authz/models"
)

type Service interface {
	List() ([]*models.Organization, error)
	Find(id int32) (*models.Organization, error)
	Create(organization *models.Organization) (*models.Organization, error)
	FirstOrCreate(organization *models.Organization) (*models.Organization, error)
	Update(organization *models.Organization) (*models.Organization, error)
	Delete(organization *models.Organization) error
	Exists(id int32) bool
	FindUsersByIds(organization *models.Organization, ids []int32) ([]*models.User, error)
	FindPermissionsByIds(organization *models.Organization, ids []int32) ([]*models.Permission, error)
}

type organizationService struct {
	db         *pg.DB
	repository Repository
}

var _ Service = (*organizationService)(nil)

func NewOrganizationService(db *pg.DB) Service {
	return &organizationService{
		repository: NewOrganizationRepository(db),
		db:         db,
	}
}

func (o *organizationService) List() ([]*models.Organization, error) {
	return o.repository.List()
}

func (o *organizationService) Exists(id int32) bool {
	return o.repository.Exists(id)
}

func (o *organizationService) Find(id int32) (*models.Organization, error) {
	return o.repository.Find(id)
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

func (o *organizationService) FindUsersByIds(organization *models.Organization, ids []int32) ([]*models.User, error) {
	return o.repository.FindUsersByIds(organization, ids)
}

func (o *organizationService) FindPermissionsByIds(organization *models.Organization, ids []int32) ([]*models.Permission, error) {
	return o.repository.FindPermissionsByIds(organization, ids)
}
