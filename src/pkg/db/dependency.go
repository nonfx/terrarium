// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package db

import (
	"github.com/cldcvr/terrarium/src/pkg/jsonschema"
	"github.com/cldcvr/terrarium/src/pkg/pb/terrariumpb"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Dependency struct {
	Model

	TaxonomyID  uuid.UUID        `gorm:"default:null"` // Given taxonomy's uncertain presence in YAML, setting TaxonomyID default as NULL accommodates potential absence of taxonomy data.
	InterfaceID string           `gorm:"unique"`
	Title       string           `gorm:"default:null"`
	Description string           `gorm:"default:null"`
	Inputs      *jsonschema.Node `gorm:"type:json"`
	Outputs     *jsonschema.Node `gorm:"type:json"`
	ExtendsID   string           `gorm:"-"` //This is yet to be finalized

	Taxonomy *Taxonomy `gorm:"foreignKey:TaxonomyID"`
}

type Dependencies []*Dependency

// insert a row in DB or in case of conflict in unique fields, update the existing record and set the existing record ID in the given object
func (db *gDB) CreateDependencyInterface(e *Dependency) (uuid.UUID, error) {
	return createOrUpdate(db.g(), e, []string{"interface_id"})
}

func (dep *Dependency) GetCondition() entity {
	return &Dependency{
		InterfaceID: dep.InterfaceID,
	}
}

func (d *Dependency) ToProto() *terrariumpb.Dependency {
	return &terrariumpb.Dependency{
		InterfaceId: d.InterfaceID,
		Title:       d.Title,
		Description: d.Description,
		Inputs:      JSONSchemaToProto(d.Inputs),
		Outputs:     JSONSchemaToProto(d.Outputs),
	}
}

func (dArr Dependencies) ToProto() []*terrariumpb.Dependency {
	res := make([]*terrariumpb.Dependency, len(dArr))
	for i, m := range dArr {
		res[i] = m.ToProto()
	}

	return res
}

func JSONSchemaToProto(jsn *jsonschema.Node) *terrariumpb.JSONSchema {
	if jsn == nil {
		return nil
	}

	// Create the base proto representation
	protoSchema := &terrariumpb.JSONSchema{
		Title:       jsn.Title,
		Description: jsn.Description,
		Type:        jsn.Type,
	}

	// If properties exist OR the type is an "object",
	// then we can convert each property
	if jsn.Properties != nil || jsn.Type == "object" {
		protoSchema.Properties = make(map[string]*terrariumpb.JSONSchema)

		for key, prop := range jsn.Properties {
			// Convert each property to its proto representation
			protoProperty := &terrariumpb.JSONSchema{
				Title:       prop.Title,
				Description: prop.Description,
				Type:        prop.Type,
			}

			// Set the converted property in the Proto Schema
			protoSchema.Properties[key] = protoProperty
		}
	}

	return protoSchema
}

func ToProto(jsn *jsonschema.Node) *terrariumpb.JSONSchema {
	if jsn == nil {
		return nil
	}
	return &terrariumpb.JSONSchema{
		Title:       jsn.Title,
		Description: jsn.Description,
		Type:        jsn.Type,
	}
}

func (db *gDB) QueryDependencyByInterfaceID(interfaceID string, filterOps ...FilterOption) (*Dependency, error) {
	q := db.g().Preload("Taxonomy").Model(&Dependency{}).Where("interface_id = ?", interfaceID)

	// Apply each filter to the query
	for _, filter := range filterOps {
		q = filter(q)
	}

	var dep Dependency
	err := q.First(&dep).Error
	if err != nil {
		return nil, err
	}

	return &dep, nil
}

func (db *gDB) QueryDependencies(filterOps ...FilterOption) (Dependencies, error) {
	q := db.g().Preload("Taxonomy").Model(&Dependency{})

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
