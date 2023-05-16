package db

import (
	"github.com/cldcvr/terrarium/api/internal/config"
	"github.com/rotisserie/eris"
)

func Connect() (DB, error) {
	db, err := config.ConnectFromEnv()
	if err != nil {
		return nil, eris.Wrap(err, "could not establish a connection to the database")
	}

	err = db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`).Error
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
	)
	if err != nil {
		return nil, eris.Wrap(err, "failed to perform database migration")
	}

	return (*gDB)(db), nil
}
