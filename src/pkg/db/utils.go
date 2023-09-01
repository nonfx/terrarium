package db

import (
	"math"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

//go:generate mockery --name DB

type DB interface {
	CreateTFProvider(e *TFProvider) (uuid.UUID, error)
	CreateTFResourceType(e *TFResourceType) (uuid.UUID, error)
	CreateTFResourceAttribute(e *TFResourceAttribute) (uuid.UUID, error)
	CreateTFResourceAttributesMapping(e *TFResourceAttributesMapping) (uuid.UUID, error)
	CreateTFModule(e *TFModule) (uuid.UUID, error)
	CreateTFModuleAttribute(e *TFModuleAttribute) (uuid.UUID, error)
	CreateTaxonomy(e *Taxonomy) (uuid.UUID, error)
	CreateDependencyInterface(e *Dependency) (uuid.UUID, error)
	GetTaxonomyByFieldName(fieldName string, fieldValue interface{}) (Taxonomy, error)

	// GetOrCreate finds and updates `e` and if the record doesn't exists, it creates a new record `e` and updates ID.
	GetOrCreateTFProvider(e *TFProvider) (id uuid.UUID, isNew bool, err error)

	GetTFProvider(e *TFProvider, where *TFProvider) error
	GetTFResourceType(e *TFResourceType, where *TFResourceType) error
	GetTFResourceAttribute(e *TFResourceAttribute, where *TFResourceAttribute) error

	// QueryTFModules list terraform modules
	QueryTFModules(filterOps ...FilterOption) (result TFModules, err error)
	// QueryTFModuleAttributes list terraform module attributes
	QueryTFModuleAttributes(filterOps ...FilterOption) (result TFModuleAttributes, err error)

	// FindOutputMappingsByModuleID DEPRECATED fetch the terraform module along with it's attribute and output mappings of the attribute.
	FindOutputMappingsByModuleID(ids ...uuid.UUID) (result TFModules, err error)
}

type FilterOption func(*gorm.DB) *gorm.DB

// Model a basic GoLang struct which includes the following fields: ID, CreatedAt, UpdatedAt, DeletedAt
type Model struct {
	ID        uuid.UUID `gorm:"type:uuid;primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (m *Model) GenerateID() {
	m.ID = uuid.New()
}

func (m *Model) GetID() uuid.UUID {
	return m.ID
}

func (m *Model) SetID(id uuid.UUID) {
	m.ID = id
}

type entity interface {
	GetID() uuid.UUID
	SetID(id uuid.UUID)
	GenerateID()
	GetCondition() entity
}

func createOrUpdate[T entity](g *gorm.DB, e T, uniqueFields []string) (uuid.UUID, error) {
	res := Model{}
	err := g.Model(e).Where(e.GetCondition()).First(&res).Error
	if err != nil && !IsNotFoundError(err) {
		return uuid.Nil, err
	}

	if res.GetID() != uuid.Nil {
		// update
		e.SetID(res.GetID())
		err = g.Save(e).Error
	} else {
		// create
		e.GenerateID()
		err = g.Create(e).Error
	}

	if err != nil {
		return uuid.Nil, err
	}
	return e.GetID(), nil
}

func (db *gDB) GetTaxonomyByFieldName(fieldName string, fieldValue interface{}) (Taxonomy, error) {
	var taxonomy Taxonomy
	uniqueFields := map[string]interface{}{
		fieldName: fieldValue,
	}
	err := db.g().Where(uniqueFields).First(&taxonomy).Error
	if err != nil {
		return Taxonomy{}, err
	}
	return taxonomy, nil
}

func get[T entity](g *gorm.DB, e T, where T) error {
	return g.First(e, where).Error
}

type gDB gorm.DB

func (db *gDB) g() *gorm.DB {
	return (*gorm.DB)(db)
}

func PaginateGlobalFilter(pageSize, pageIndex int32, totalPages *int32) FilterOption {
	offset := pageIndex * pageSize
	return func(g *gorm.DB) *gorm.DB {
		var count int64
		_ = g.Count(&count).Error
		(*totalPages) = int32(math.Ceil(float64(count) / float64(pageSize)))
		return g.Offset(int(offset)).Limit(int(pageSize))
	}
}

func NoOpFilter(g *gorm.DB) *gorm.DB {
	return g
}
