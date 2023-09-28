// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package db

import "github.com/google/uuid"

type DependencyAttributeMappings struct {
	Model

	DependencyAttributeID uuid.UUID `gorm:"uniqueIndex:dependency_attribute_mapping_unique"`
	ResourceAttributeID   uuid.UUID `gorm:"uniqueIndex:dependency_attribute_mapping_unique"`
}
