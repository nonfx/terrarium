// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package db

import (
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
	)
	if err != nil {
		return nil, eris.Wrap(err, "failed to perform database migration")
	}

	db.Exec("insert into taxonomies(id) values('00000000-0000-0000-0000-000000000000')")

	return (*gDB)(db), nil
}
