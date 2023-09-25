// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

//go:build dbtest
// +build dbtest

package db_test

import (
	"testing"

	"github.com/cldcvr/terrarium/src/pkg/db"
	"github.com/cldcvr/terrarium/src/pkg/jsonschema"
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

func Test_gDB_QueryDependencyByInterfaceID(t *testing.T) {
	tests := []struct {
		name        string
		interfaceID string
		filters     []db.FilterOption
		validator   func(*testing.T, *db.Dependency)
		wantErr     bool
	}{
		{
			name:        "query by InterfaceID",
			interfaceID: "dependency-1-interface",
			filters: []db.FilterOption{
				db.DependencySearchFilter("dependency-1-interface"),
			},
			validator: func(t *testing.T, result *db.Dependency) {
				require.NotNil(t, result)
				assert.Equal(t, "dependency-1-interface", result.InterfaceID)
				assert.Equal(t, "dependency-1", result.Title)
				assert.Equal(t, "this is first test dependency", result.Description)
			},
			wantErr: false,
		},
		{
			name:        "non-existent InterfaceID",
			interfaceID: "non-existent",
			filters:     []db.FilterOption{},
			validator: func(t *testing.T, result *db.Dependency) {
				assert.Nil(t, result)
			},
			wantErr: true,
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
					result, err := dbObj.QueryDependencyByInterfaceID(tt.interfaceID, tt.filters...)
					if tt.wantErr {
						assert.Error(t, err)
						assert.Nil(t, result)
					} else {
						assert.NoError(t, err)
						tt.validator(t, result)
					}
				})
			}
		})
	}
}

func TestJSONSchemaToProto(t *testing.T) {
	tests := []struct {
		name      string
		jsn       *jsonschema.Node
		validator func(*testing.T, *terrariumpb.JSONSchema)
		wantErr   bool
	}{
		{
			name: "valid JSONSchema",
			jsn: &jsonschema.Node{
				Title:       "Test Title",
				Description: "Test Description",
				Type:        "object",
				Properties: map[string]*jsonschema.Node{
					"property1": {
						Title: "Property 1",
						Type:  "string",
					},
					"property2": {
						Title: "Property 2",
						Type:  "integer",
					},
				},
			},
			validator: func(t *testing.T, result *terrariumpb.JSONSchema) {
				require.NotNil(t, result)
				assert.Equal(t, "Test Title", result.Title)
				assert.Equal(t, "Test Description", result.Description)
				assert.Equal(t, "object", result.Type)
				require.Len(t, result.Properties, 2)
				assert.NotNil(t, result.Properties["property1"])
				assert.NotNil(t, result.Properties["property2"])
			},
			wantErr: false,
		},
		{
			name: "nil JSONSchema",
			jsn:  nil,
			validator: func(t *testing.T, result *terrariumpb.JSONSchema) {
				require.Nil(t, result)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := db.JSONSchemaToProto(tt.jsn)
			tt.validator(t, result)
		})
	}
}

func TestToProto(t *testing.T) {
	tests := []struct {
		name      string
		jsn       *jsonschema.Node
		validator func(*testing.T, *terrariumpb.JSONSchema)
	}{
		{
			name: "nil JSONSchema",
			jsn:  nil,
			validator: func(t *testing.T, result *terrariumpb.JSONSchema) {
				require.Nil(t, result)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := db.ToProto(tt.jsn)
			tt.validator(t, result)
		})
	}
}
