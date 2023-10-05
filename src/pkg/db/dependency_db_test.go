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
	"github.com/xeipuuv/gojsonschema"
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
				db.DependencySearchFilter("ependency-1-int"),
			},
			wantModule: []*terrariumpb.Dependency{
				{
					Id:          uuidDep1.String(),
					InterfaceId: "dependency-1-interface",
					Title:       "dependency-1",
					Description: "this is first test dependency",
					Taxonomy:    []string{"mockdata-l1", "mockdata-l2", "mockdata-l3", "mockdata-l4", "mockdata-l5", "mockdata-l6", "mockdata-l7"},
					Inputs: &terrariumpb.JSONSchema{
						Type: gojsonschema.TYPE_OBJECT,
						Properties: map[string]*terrariumpb.JSONSchema{
							"dep-1-attr-1": {
								Title:       "Attr 1",
								Description: "attribute 1",
								Type:        gojsonschema.TYPE_NUMBER,
							},
						},
					},
					Outputs: &terrariumpb.JSONSchema{
						Type: gojsonschema.TYPE_OBJECT,
						Properties: map[string]*terrariumpb.JSONSchema{
							"dep-1-attr-2": {
								Title:       "Attr 2",
								Description: "attribute 2",
								Type:        gojsonschema.TYPE_NUMBER,
							},
							"dep-1-attr-3": {
								Title:       "Attr 3",
								Description: "attribute 3",
								Type:        gojsonschema.TYPE_NUMBER,
							},
						},
					},
				},
			},
		},
		{
			name: "empty query return everything",
			filters: []db.FilterOption{
				db.DependencySearchFilter(""),
			},
			validator: func(t *testing.T, d db.Dependencies) {
				assert.GreaterOrEqual(t, len(d), 1, "length of returned results")
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
					} else if tt.validator != nil {
						assert.NoError(t, err)
						tt.validator(t, gotResult)
					} else {
						assert.NoError(t, err)
						assert.EqualValues(t, tt.wantModule, gotResult.ToProto())
					}
				})
			}
		})
	}
}
