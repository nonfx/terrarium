// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package db

import (
	"github.com/cldcvr/terrarium/src/pkg/pb/terrariumpb"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Dependency struct {
	Model

	TaxonomyID  uuid.UUID `gorm:"default:null"` // Given taxonomy's uncertain presence in YAML, setting TaxonomyID default as NULL accommodates potential absence of taxonomy data.
	InterfaceID string    `gorm:"unique"`
	Title       string    `gorm:"default:null"`
	Description string    `gorm:"default:null"`
	ExtendsID   string    `gorm:"-"` //This is yet to be finalized

	Attributes DependencyAttributes `gorm:"foreignKey:DependencyID"`
	Taxonomy   *Taxonomy            `gorm:"foreignKey:TaxonomyID"`
}

type Dependencies []Dependency

func (depArr Dependencies) ToProto() []*terrariumpb.Dependency {
	var res []*terrariumpb.Dependency
	for _, output := range depArr {
		res = append(res, output.ToProto())
	}
	return res
}

func (d Dependency) ToProto() *terrariumpb.Dependency {
	protoDep := &terrariumpb.Dependency{
		Id:          d.ID.String(),
		InterfaceId: d.InterfaceID,
		Title:       d.Title,
		Description: d.Description,
	}

	if d.Attributes != nil {
		protoDep.Inputs = d.Attributes.GetByCompute(false).ToJSONSchema().ToProto()
		protoDep.Outputs = d.Attributes.GetByCompute(true).ToJSONSchema().ToProto()
	}

	if d.Taxonomy != nil {
		protoDep.Taxonomy = d.Taxonomy.ToLevels()
	}

	return protoDep
}

// insert a row in DB or in case of conflict in unique fields, update the existing record and set the existing record ID in the given object
func (db *gDB) CreateDependencyInterface(e *Dependency) (uuid.UUID, error) {
	id, _, _, err := createOrGetOrUpdate(db.g(), e, []string{"interface_id"})
	return id, err
}

func (db *gDB) QueryDependencies(filterOps ...FilterOption) (Dependencies, error) {
	q := db.g().Preload("Taxonomy").Model(&Dependency{})
	q = q.Preload("Attributes")

	// Apply each filter to the query
	for _, filter := range filterOps {
		q = filter(q)
	}

	var deps Dependencies
	err := q.Find(&deps).Error
	if err != nil {
		return nil, err
	}

	return deps, nil
}

func DependencySearchFilter(query string) FilterOption {
	if query == "" {
		return NoOpFilter
	}

	return func(g *gorm.DB) *gorm.DB {
		q := "%" + query + "%"
		return g.Where("interface_id LIKE ? OR title LIKE ?", q, q)
	}
}

var results []DependencyResult

func (db *gDB) Fetchdeps() []DependencyResult {
	db.g().Model(&Dependency{}).
		Select("dependencies.id AS DependencyID, dependencies.interface_id AS InterfaceID").
		Joins("INNER JOIN dependency_attributes ON dependencies.id = dependency_attributes.dependency_id").
		Scan(&results)

	return results
}
