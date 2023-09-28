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

	TaxonomyID  uuid.UUID `gorm:"default:null"` // Given taxonomy's uncertain presence in YAML, setting TaxonomyID default as NULL accommodates potential absence of taxonomy data.
	InterfaceID string    `gorm:"unique"`
	Title       string    `gorm:"default:null"`
	Description string    `gorm:"default:null"`
	ExtendsID   string    `gorm:"-"` //This is yet to be finalized

	Attributes []DependencyAttribute `gorm:"foreignKey:DependencyID"`
	Taxonomy   *Taxonomy             `gorm:"foreignKey:TaxonomyID"`
}

type Dependencies []*Dependency

type DependencyOutput struct {
	Dependency
	Inputs  *jsonschema.Node `json:"inputs"`
	Outputs *jsonschema.Node `json:"outputs"`
}

type DependencyOutputs []DependencyOutput

func (do DependencyOutputs) ToProto() []*terrariumpb.Dependency {
	var res []*terrariumpb.Dependency
	for _, output := range do {
		res = append(res, output.ToProto())
	}
	return res
}

func (do *DependencyOutput) ToProto() *terrariumpb.Dependency {
	return &terrariumpb.Dependency{
		InterfaceId: do.InterfaceID,
		Title:       do.Title,
		Description: do.Description,
		Inputs:      JSONSchemaToProto(do.Inputs),
		Outputs:     JSONSchemaToProto(do.Outputs),
	}
}

// insert a row in DB or in case of conflict in unique fields, update the existing record and set the existing record ID in the given object
func (db *gDB) CreateDependencyInterface(e *Dependency) (uuid.UUID, error) {
	return createOrUpdate(db.g(), e, []string{"interface_id"})
}

func (dep *Dependency) GetCondition() entity {
	return &Dependency{
		InterfaceID: dep.InterfaceID,
	}
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
			protoSchema.Properties[key] = JSONSchemaToProto(prop)
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

func (db *gDB) QueryDependencies(filterOps ...FilterOption) (DependencyOutputs, error) {
	q := db.g().Preload("Taxonomy").Model(&Dependency{})
	q = q.Preload("Attributes")
	q = q.Joins("left join dependency_attributes on dependency_attributes.dependency_id = dependencies.id")

	// Apply each filter to the query
	for _, filter := range filterOps {
		q = filter(q)
	}

	var deps []Dependency
	err := q.Find(&deps).Error
	if err != nil {
		return nil, err
	}

	depOutputs := make([]DependencyOutput, len(deps))
	for i, dep := range deps {
		depOutput := DependencyOutput{Dependency: dep}
		for _, attr := range dep.Attributes {
			if attr.Computed {
				depOutput.Outputs = attr.Schema
			} else {
				depOutput.Inputs = attr.Schema
			}
		}
		depOutputs[i] = depOutput
	}

	return depOutputs, nil
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
