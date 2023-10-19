// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

//go:build dbtest
// +build dbtest

package dbhelper

import (
	"path/filepath"
	"testing"

	"github.com/cldcvr/terrarium/src/pkg/confighelper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDialectorSwitcherConnect(t *testing.T) {
	testsqliteloc := filepath.Join(t.TempDir(), "t8db", "testdb.sqlite")

	confighelper.LoadDefaults(map[string]interface{}{
		"db": map[string]interface{}{
			"host":     "localhost",
			"user":     "postgres",
			"password": "",
			"name":     "cc_terrarium",
			"port":     5432,
			"ssl_mode": false,
		},
	}, "tr")

	tests := []struct {
		name         string
		conf         DialectorSwitcher
		dbType       DBDriver
		wantErr      string
		wantDialName string
	}{
		{
			name: "sqlite fail no file",
			conf: DialectorSwitcher{
				ConfigSQLite: ConfigSQLite{
					DSN: testsqliteloc,
				},
			},
			dbType:  DBDriverSQLite,
			wantErr: "unable to open database file: no such file or directory",
		},
		{
			name: "sqlite success create file",
			conf: DialectorSwitcher{
				ConfigSQLite: ConfigSQLite{
					DSN:                  testsqliteloc,
					ResolvePathCreateDir: true,
				},
			},
			dbType:       DBDriverSQLite,
			wantDialName: "sqlite",
		},
		{
			name: "sqlite success use existing file",
			conf: DialectorSwitcher{
				ConfigSQLite: ConfigSQLite{
					DSN: testsqliteloc,
				},
			},
			dbType:       DBDriverSQLite,
			wantDialName: "sqlite",
		},
		{
			name: "pg fail config issue",
			conf: DialectorSwitcher{
				ConfigPostgres: ConfigPostgres{},
			},
			dbType:  DBDriverPostgres,
			wantErr: "hostname resolving error",
		},
		{
			name: "pg success",
			conf: DialectorSwitcher{
				ConfigPostgres: ConfigPostgres{
					Host:     confighelper.MustGetString("db.host"),
					User:     confighelper.MustGetString("db.user"),
					Password: confighelper.MustGetString("db.password"),
					DBName:   confighelper.MustGetString("db.name"),
					Port:     confighelper.MustGetInt("db.port"),
					SslMode:  confighelper.MustGetBool("db.ssl_mode"),
				},
			},
			dbType:       DBDriverPostgres,
			wantDialName: "postgres",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g, err := tt.conf.Connect(tt.dbType, WithRetries(0, 0, 0))
			if tt.wantErr != "" {
				assert.ErrorContains(t, err, tt.wantErr)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.wantDialName, g.Dialector.Name())
		})
	}

	assert.FileExists(t, testsqliteloc)
}
