package db

import (
	"github.com/rotisserie/eris"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) (DB, error) {
	err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`).Error
	if err != nil {
		return nil, eris.Wrap(err, "failed to perform database migration")
	}

	err = db.AutoMigrate(
		TFProvider{},
		TFResourceType{},
		TFResourceAttribute{},
		TFResourceAttributesMapping{},
		TFModule{},
		TFModuleAttribute{},
		Taxonomy{},
		Dependency{},
	)
	if err != nil {
		return nil, eris.Wrap(err, "failed to perform database migration")
	}

	return (*gDB)(db), nil
}
