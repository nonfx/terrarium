// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package db

import "github.com/google/uuid"

type PlatformComponent struct {
	Model

	PlatformID   uuid.UUID `gorm:"uniqueIndex:platform_components_unique"`
	DependencyID uuid.UUID `gorm:"uniqueIndex:platform_components_unique"`

	Platform   Platform   `gorm:"foreignKey:PlatformID"`
	Dependency Dependency `gorm:"foreignKey:DependencyID"`
}

// insert a row in DB or in case of conflict in unique fields, update the existing record and set existing record ID in the given object
func (db *gDB) CreatePlatformComponents(p *PlatformComponent) (uuid.UUID, error) {
	id, _, _, err := createOrGetOrUpdate(db.g(), p, []string{"platform_id", "dependency_id"})
	return id, err
}
