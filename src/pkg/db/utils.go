// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package db

import (
	"errors"
	"math"
	"time"

	"github.com/cldcvr/terrarium/src/pkg/utils"
	"github.com/google/uuid"
	"github.com/rotisserie/eris"
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
	CreateDependencyAttribute(e *DependencyAttribute) (uuid.UUID, error)
	CreatePlatform(p *Platform) (uuid.UUID, error)
	CreatePlatformComponents(p *PlatformComponents) (uuid.UUID, error)

	// GetOrCreateTFProvider finds and updates `e` and if the record doesn't exists, it creates a new record `e` and updates ID.
	GetOrCreateTFProvider(e *TFProvider) (id uuid.UUID, isNew bool, err error)

	GetTFProvider(e *TFProvider, where *TFProvider) error
	GetTFResourceType(e *TFResourceType, where *TFResourceType) error
	GetTFResourceAttribute(e *TFResourceAttribute, where *TFResourceAttribute) error

	// QueryTFModules list terraform modules
	QueryTFModules(filterOps ...FilterOption) (result TFModules, err error)
	// QueryTFModuleAttributes list terraform module attributes
	QueryTFModuleAttributes(filterOps ...FilterOption) (result TFModuleAttributes, err error)

	QueryDependencies(filterOps ...FilterOption) (result Dependencies, err error)

	ExecuteSQLStatement(string) error
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

type entity[E any] interface {
	*E

	GetID() uuid.UUID
	SetID(id uuid.UUID)
	GenerateID()
}

type entityEq[E any] interface {
	// IsEq checks if the entity is equal to the given entity.
	// When an entity implements this, we skip update if the
	// their is no change in the entity row.
	IsEq(*E) bool
}

// createOrGetOrUpdate if the record does not exist, then create new and return id,
// if the record exists, and has no change, then return id,
// if the record exists, and may have change (has change or non deterministic), then update and return id.
func createOrGetOrUpdate[T entity[E], E any](g *gorm.DB, e T, uniqueFields []string) (id uuid.UUID, isNew bool, isUpdated bool, err error) {
	var dbObj E // create new empty object

	err = g.Where(e, utils.ToIfaceArr(uniqueFields)...).First(&dbObj).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil // ignore not found error
	} else if err != nil {
		err = eris.Wrap(err, "db get failed")
		return
	}

	if T(&dbObj).GetID() == uuid.Nil {
		// create
		isNew = true
		e.GenerateID()
		err = g.Create(e).Error
		if err != nil {
			err = eris.Wrap(err, "db create failed")
			return
		}
	} else if eEq, ok := ((interface{})(e)).(entityEq[E]); ok && eEq.IsEq(&dbObj) {
		// no update when the object in DB vs given objects are equal.
		e.SetID(T(&dbObj).GetID())
	} else {
		// update
		isUpdated = true
		e.SetID(T(&dbObj).GetID())
		err = g.Save(e).Error
		if err != nil {
			err = eris.Wrap(err, "db update failed")
			return
		}
	}

	id = e.GetID()

	return
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

func (db *gDB) ExecuteSQLStatement(statement string) error {
	return db.g().Exec(statement).Error
}
