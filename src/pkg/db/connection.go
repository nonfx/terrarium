// Copyright (c) Ollion
// SPDX-License-Identifier: Apache-2.0

package db

import (
	"github.com/google/uuid"
	"github.com/rotisserie/eris"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) (DB, error) {
	err := db.AutoMigrate(
		TFProvider{},
		TFResourceType{},
		TFResourceAttribute{},
		TFResourceAttributesMapping{},
		TFModule{},
		TFModuleAttribute{},
		Taxonomy{},
		Dependency{},
		DependencyAttribute{},
		DependencyAttributeMappings{},
		Platform{},
		PlatformComponent{},
		FarmRelease{},
	)
	if err != nil {
		return nil, eris.Wrap(err, "failed to perform database migration")
	}

	db.FirstOrCreate(&Taxonomy{Model: Model{ID: uuid.Nil}})

	return (*gDB)(db), nil
}
