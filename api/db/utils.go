package db

import (
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

	GetTFResourceType(e *TFResourceType, where *TFResourceType) error
	GetTFResourceAttribute(e *TFResourceAttribute, where *TFResourceAttribute) error
}

// Model a basic GoLang struct which includes the following fields: ID, CreatedAt, UpdatedAt, DeletedAt
// It may be embedded into your model or you may build your own model without it
//
//	type User struct {
//	  gorm.Model
//	}
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
