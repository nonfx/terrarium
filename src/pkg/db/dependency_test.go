// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

//go:build dbtest
// +build dbtest

package db_test

import (
	"testing"

	"github.com/cldcvr/terrarium/src/pkg/db"
	"github.com/cldcvr/terrarium/src/pkg/pb/terrariumpb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_gDB_QueryDependencies(t *testing.T) {
	tests := []struct {
		name       string
		filters    []db.FilterOption
		validator  func(*testing.T, db.Dependencies)
		wantModule []*terrariumpb.Dependency //DependencyOutputs
		wantErr    bool
	}{
		{
			name: "query by InterfaceID",
			filters: []db.FilterOption{
				db.DependencySearchFilter("dependency-1-interface"),
			},
			wantModule: []*terrariumpb.Dependency{
				{
					InterfaceId: "dependency-1-interface",
					Title:       "dependency-1",
					Description: "this is first test dependency",
				},
			},
		},
	}
	for dbName, connector := range getConnectorMap() {
		g := connector(t)
		dbObj, err := db.AutoMigrate(g)
		require.NoError(t, err)
		saveTestData(t, g)

		t.Run(dbName, func(t *testing.T) {
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					gotResult, err := dbObj.QueryDependencies(tt.filters...)
					if tt.wantErr {
						assert.Error(t, err)
					} else {
						assert.NoError(t, err)
						assert.EqualValues(t, tt.wantModule, gotResult.ToProto())
					}
				})
			}
		})
	}
}
