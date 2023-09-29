// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

//go:build dbtest
// +build dbtest

package db_test

import (
	"testing"

	"github.com/cldcvr/terrarium/src/pkg/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_gDB_CreatePlatform(t *testing.T) {
	tests := []struct {
		name     string
		platform *db.Platform
		wantErr  bool
	}{
		{
			name:     "first new insert",
			platform: &db.Platform{Name: "test-1", RepoURL: "test-url", RepoDirectory: "test-dir", CommitSHA: "test-sha", RefLabel: "test-ref", LabelType: 1},
		},
	}

	for dbName, connector := range getConnectorMap() {
		dbObj, err := db.AutoMigrate(connector(t))
		require.NoError(t, err)

		t.Run(dbName, func(t *testing.T) {
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					_, err := dbObj.CreatePlatform(tt.platform)
					if tt.wantErr {
						assert.Error(t, err)
					} else {
						assert.NoError(t, err)
					}
				})
			}
		})
	}
}
