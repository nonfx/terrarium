package db

import (
	"github.com/cldcvr/terrarium/src/pkg/jsonschema"
	"github.com/cldcvr/terrarium/src/pkg/pb/terrariumpb"
	"github.com/google/uuid"
)

type DependencyAttribute struct {
	Model

	DependnecyID uuid.UUID `gorm:"uniqueIndex:dependency_attribute_unique"`
	Name         string    `gorm:"uniqueIndex:dependency_attribute_unique"`
	Schema       *jsonschema.Node
	Computed     bool // true means output, false means input

	Dependency *Dependency `gorm:"foreignKey:DependnecyID"`
}

type DependencyAttributeMappings struct {
	Model

	DependencyAttributeID uuid.UUID `gorm:"uniqueIndex:dependency_attribute_mapping_unique"`
	ResourceAttributeID   uuid.UUID `gorm:"uniqueIndex:dependency_attribute_mapping_unique"`
}

type DependencyAttributes []*DependencyAttribute

func (dbAttr DependencyAttribute) ToProto() *terrariumpb.DependencyInputsAndOutputs {
	resp := &terrariumpb.DependencyInputsAndOutputs{
		Title:       dbAttr.Name,
		Description: dbAttr.Schema.Description,
		Type:        dbAttr.Schema.Type,
	}

	return resp
}

func (dbAttrs DependencyAttributes) ToProto() []*terrariumpb.DependencyInputsAndOutputs {
	resp := make([]*terrariumpb.DependencyInputsAndOutputs, len(dbAttrs))

	for i, dbAttr := range dbAttrs {
		resp[i] = dbAttr.ToProto()
	}

	return resp
}
