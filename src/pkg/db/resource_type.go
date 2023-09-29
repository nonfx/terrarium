// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package db

import "github.com/google/uuid"

type TFResourceType struct {
	Model

	ProviderID   uuid.UUID `gorm:"uniqueIndex:resource_type_unique"`
	ResourceType string    `gorm:"uniqueIndex:resource_type_unique"`
	TaxonomyID   uuid.UUID `gorm:"default:null"`

	Provider TFProvider `gorm:"foreignKey:ProviderID"`
	Taxonomy *Taxonomy  `gorm:"foreignKey:TaxonomyID"`
}

// insert a row in DB or in case of conflict in unique fields, update the existing record and set existing record ID in the given object
func (db *gDB) CreateTFResourceType(e *TFResourceType) (uuid.UUID, error) {
	id, _, _, err := createOrGetOrUpdate(db.g(), e, []string{"provider_id", "resource_type"})
	return id, err
}

func (db *gDB) GetTFResourceType(e *TFResourceType, where *TFResourceType) error {
	return db.g().First(e, where).Error
}
