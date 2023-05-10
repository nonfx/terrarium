package db

import (
	"github.com/cldcvr/terrarium/api/pkg/dbcon"
	"github.com/rotisserie/eris"
	"gorm.io/gorm"
)

func Connect() (*gorm.DB, error) {
	db, err := dbcon.ConnectFromEnv()
	if err != nil {
		return nil, eris.Wrap(err, "could not establish a connection to the database")
	}

	err = db.AutoMigrate(
		TFProvider{},
		TFResourceType{},
		TFResourceAttribute{},
		TFResourceAttributesMapping{},
		TFModule{},
		TFModuleAttribute{},
		Taxonomy{},
	)
	if err != nil {
		return nil, eris.Wrap(err, "failed to perform database migration")
	}

	return db, nil
}
