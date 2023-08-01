package db

import (
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

//go:generate mockery --all

type DB interface {
	CreateTFProvider(e *TFProvider) (uuid.UUID, error)
	CreateTFResourceType(e *TFResourceType) (uuid.UUID, error)
	CreateTFResourceAttribute(e *TFResourceAttribute) (uuid.UUID, error)
	CreateTFResourceAttributesMapping(e *TFResourceAttributesMapping) (uuid.UUID, error)
	CreateTFModule(e *TFModule) (uuid.UUID, error)
	CreateTFModuleAttribute(e *TFModuleAttribute) (uuid.UUID, error)
	CreateTaxonomy(e *Taxonomy) (uuid.UUID, error)

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

	CreateDependencyInterface(e *Dependency) (uuid.UUID, error)
}

type FilterOption func(*gorm.DB) *gorm.DB

// Model a basic GoLang struct which includes the following fields: ID, CreatedAt, UpdatedAt, DeletedAt
type Model struct {
	ID        uuid.UUID `gorm:"type:uuid;primarykey;default:uuid_generate_v4()"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (m *Model) GetID() uuid.UUID {
	return m.ID
}

type entity interface {
	GetID() uuid.UUID
}

func createOrUpdate[T entity](g *gorm.DB, e T, uniqueFields []string) (uuid.UUID, error) {
	c := clause.OnConflict{
		Columns:   []clause.Column{},
		UpdateAll: true,
	}

	for _, f := range uniqueFields {
		c.Columns = append(c.Columns, clause.Column{Name: f})
	}

	err := g.Clauses(c).Create(e).Error
	if err != nil {
		fmt.Printf("Error in createOrUpdate: %v\n", err)
		return uuid.Nil, err
	}

	return e.GetID(), nil
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
