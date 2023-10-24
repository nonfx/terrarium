// Copyright (c) Ollion
// SPDX-License-Identifier: Apache-2.0

package db

import "github.com/google/uuid"

type TFResourceAttributesMapping struct {
	Model

	InputAttributeID  uuid.UUID `gorm:"uniqueIndex:resource_attribute_mapping_unique"`
	OutputAttributeID uuid.UUID `gorm:"uniqueIndex:resource_attribute_mapping_unique"`

	InputAttribute  TFResourceAttribute `gorm:"foreignKey:InputAttributeID"`  // Resource input-attribute object
	OutputAttribute TFResourceAttribute `gorm:"foreignKey:OutputAttributeID"` // Resource attribute object that provides the input-attribute
}

// insert a row in DB or in case of conflict in unique fields, update the existing record and set existing record ID in the given object
func (db *gDB) CreateTFResourceAttributesMapping(e *TFResourceAttributesMapping) (uuid.UUID, error) {
	id, _, _, err := createOrGetOrUpdate(db.g(), e, []string{"input_attribute_id", "output_attribute_id"})
	return id, err
}
