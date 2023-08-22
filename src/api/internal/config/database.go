package config

import (
	"github.com/cldcvr/terrarium/src/pkg/confighelper"
	"github.com/cldcvr/terrarium/src/pkg/db"
	"github.com/cldcvr/terrarium/src/pkg/db/dbhelper"
	"github.com/rotisserie/eris"
)

// DBHost returns the database host
func DBHost() string {
	return confighelper.MustGetString("db.host")
}

// DBUser returns the database user
func DBUser() string {
	return confighelper.MustGetString("db.user")
}

// DBPassword returns the database password
func DBPassword() string {
	return confighelper.MustGetString("db.password")
}

// DBName returns the database name
func DBName() string {
	return confighelper.MustGetString("db.name")
}

// DBPort returns the database port
func DBPort() int {
	return confighelper.MustGetInt("db.port")
}

// DBSSLMode returns the database SSL mode.
func DBSSLMode() bool {
	return confighelper.MustGetBool("db.ssl_mode")
}

// DBConnect establishes a connection to the database using the connection parameters from the environment.
func DBConnect() (db.DB, error) {
	g, err := dbhelper.ConnectPG(
		DBHost(),
		DBUser(),
		DBPassword(),
		DBName(),
		DBPort(),
		DBSSLMode(),
		dbhelper.WithRetries(20, 3, 3),
	)
	if err != nil {
		return nil, eris.Wrap(err, "could not establish a connection to the database")
	}

	return db.AutoMigrate(g)
}
