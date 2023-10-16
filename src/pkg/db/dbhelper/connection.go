// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package dbhelper

import (
	"fmt"
	"path/filepath"

	"github.com/cldcvr/terrarium/src/pkg/utils"
	"github.com/rotisserie/eris"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DBDriver uint

const (
	DBDriverUndefined = iota
	DBDriverSQLite
	DBDriverPostgres
)

type ConfigPostgres struct {
	Host     string
	User     string
	Password string
	DBName   string
	Port     int
	SslMode  bool
}

type ConfigSQLite struct {
	ResolvePathCreateDir bool
	DSN                  string
	err                  error
}

type DialectorGetter interface {
	GetGormDialector() (gorm.Dialector, error)
}

type DialectorSwitcher struct {
	ConfigPostgres
	ConfigSQLite
}

// GetDSN creates the DSN (Data Source Name) string for the postgres database connection.
func (c ConfigPostgres) GetDSN() string {
	sslModeStr := "disable"
	if c.SslMode {
		sslModeStr = "enable"
	}

	passwordStr := ""
	if c.Password != "" {
		// omitting the password block when not set, allows the library to look at
		// other standard sources like `~/.pgpass`
		passwordStr = "password=" + c.Password
	}

	return fmt.Sprintf("host=%s user=%s %s dbname=%s port=%d sslmode=%s", c.Host, c.User, passwordStr, c.DBName, c.Port, sslModeStr)
}

func (c ConfigPostgres) GetGormDialector() (gorm.Dialector, error) {
	return postgres.Open(c.GetDSN()), nil
}

func (c ConfigSQLite) GetGormDialector() (gorm.Dialector, error) {
	resolvedDSN := c.DSN
	if c.ResolvePathCreateDir {
		dir := filepath.Dir(c.DSN)
		dir, err := utils.SetupDir(dir)
		if err != nil {
			return nil, err
		}
		resolvedDSN = filepath.Join(dir, filepath.Ext(c.DSN))
	}

	return sqlite.Open(resolvedDSN), nil
}

func DBDriverFromStr(dbType string) DBDriver {
	d, _ := (map[string]DBDriver{
		"postgres": DBDriverPostgres,
		"sqlite":   DBDriverSQLite,
	})[dbType]
	return d
}

func (dbType DBDriver) String() string {
	d, _ := (map[DBDriver]string{
		DBDriverPostgres: "postgres",
		DBDriverSQLite:   "sqlite",
	})[dbType]
	return d
}

func (conf DialectorSwitcher) Connect(dbType DBDriver, ops ...ConnOption) (*gorm.DB, error) {
	d, err := conf.Switch(dbType)
	if err != nil {
		return nil, err
	}

	return Connect(d, ops...)
}

func (conf DialectorSwitcher) Switch(dbType DBDriver) (gorm.Dialector, error) {
	d, ok := (map[DBDriver]DialectorGetter{
		DBDriverPostgres: conf.ConfigPostgres,
		DBDriverSQLite:   conf.ConfigSQLite,
	})[dbType]
	if !ok {
		return nil, eris.Errorf("invalid db driver: %d", dbType)
	}

	return d.GetGormDialector()
}
