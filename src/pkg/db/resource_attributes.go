// Copyright (c) Ollion
// SPDX-License-Identifier: Apache-2.0

package db

import "github.com/google/uuid"

type TFResourceAttribute struct {
	Model

	ResourceTypeID uuid.UUID `gorm:"uniqueIndex:resource_attribute_unique"`
	ProviderID     uuid.UUID `gorm:"uniqueIndex:resource_attribute_unique"`
	AttributePath  string    `gorm:"uniqueIndex:resource_attribute_unique"`
	DataType       string
	Description    string
	Optional       bool
	Computed       bool

	ResourceType       TFResourceType                `gorm:"foreignKey:ResourceTypeID"`
	Provider           TFProvider                    `gorm:"foreignKey:ProviderID"`
	RelatedModuleAttrs []TFModuleAttribute           `gorm:"foreignKey:RelatedResourceTypeAttributeID"` // Module attributes that relates to this resource attribute
	OutputMappings     []TFResourceAttributesMapping `gorm:"foreignKey:InputAttributeID"`               // Mappings with another resources's output attribute
	InputMappings      []TFResourceAttributesMapping `gorm:"foreignKey:OutputAttributeID"`              // Mappings with another resources's input attribute
}

func (a TFResourceAttribute) GetConnectedModuleOutputs() TFModuleAttributes {
	resp := TFModuleAttributes{}
	for _, om := range a.OutputMappings {
		resp = append(resp, om.OutputAttribute.RelatedModuleAttrs...)
	}
	return resp
}

// insert a row in DB or in case of conflict in unique fields, update the existing record and set existing record ID in the given object
func (db *gDB) CreateTFResourceAttribute(e *TFResourceAttribute) (uuid.UUID, error) {
	id, _, _, err := createOrGetOrUpdate(db.g(), e, []string{"resource_type_id", "provider_id", "attribute_path"})
	return id, err
}

func (db *gDB) GetTFResourceAttribute(e *TFResourceAttribute, where *TFResourceAttribute) error {
	return db.g().First(e, where).Error
}
