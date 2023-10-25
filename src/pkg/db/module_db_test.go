// Copyright (c) Ollion
// SPDX-License-Identifier: Apache-2.0

//go:build dbtest
// +build dbtest

package db_test

import (
	"testing"

	"github.com/cldcvr/terrarium/src/pkg/db"
	"github.com/cldcvr/terrarium/src/pkg/pb/terrariumpb"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_gDB_CreateTFModule(t *testing.T) {
	tests := []struct {
		name    string
		module  *db.TFModule
		wantErr bool
	}{
		{
			name:   "first new insert",
			module: &db.TFModule{Source: "source-1", Version: "1.1", Namespace: "unit-test"},
		},
		{
			name:   "redundant insert",
			module: &db.TFModule{Source: "source-1", Version: "1.1", Namespace: "unit-test"},
		},
		{
			name:   "same source different version",
			module: &db.TFModule{Source: "source-1", Version: "1.2", Namespace: "unit-test"},
		},
		{
			name:   "same version different source",
			module: &db.TFModule{Source: "source-2", Version: "1.1", Namespace: "unit-test"},
		},
		{
			name:   "new insert",
			module: &db.TFModule{Source: "source-3", Version: "", Namespace: "unit-test"},
		},
	}

	for dbName, connector := range getConnectorMap() {
		dbObj, err := db.AutoMigrate(connector(t))
		require.NoError(t, err)

		t.Run(dbName, func(t *testing.T) {
			moduleIDByNames := map[string]uuid.UUID{}
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					id, err := dbObj.CreateTFModule(tt.module)
					if tt.wantErr {
						assert.Error(t, err)
					} else {
						assert.NoError(t, err)
						uniqueFieldsJoined := tt.module.Source + "|" + string(tt.module.Version)
						if wantID, ok := moduleIDByNames[uniqueFieldsJoined]; ok {
							assert.Equal(t, wantID, id)
						} else {
							assert.NotEqual(t, uuid.Nil, id)
							moduleIDByNames[uniqueFieldsJoined] = id
						}
					}
				})
			}
		})
	}
}

func Test_gDB_QueryTFModules(t *testing.T) {
	tests := []struct {
		name       string
		filters    []db.FilterOption
		validator  func(*testing.T, db.TFModules)
		wantModule []*terrariumpb.Module
		wantErr    bool
	}{
		{
			name: "query by id",
			filters: []db.FilterOption{
				db.ModuleByIDsFilter(modules[0].ID),
			},
			wantModule: []*terrariumpb.Module{
				{
					Id:              "f47ac10b-58cc-4372-a567-0e02b2c3d479",
					TaxonomyId:      "e6fb062d-74d6-4491-80bc-5d2c8e6d9ebb",
					ModuleName:      "module-1",
					Source:          "module-1-source-1",
					Version:         "1.1",
					Description:     "this is first test module",
					InputAttributes: []*terrariumpb.ModuleAttribute{},
					Namespace:       "unit-test",
				},
			},
		},
		{
			name: "search by name & namespace",
			filters: []db.FilterOption{
				db.ModuleSearchFilter(modules[0].ModuleName),
				db.ModuleNamespaceFilter([]string{modules[0].Namespace}),
			},
			wantModule: []*terrariumpb.Module{
				{
					Id:              "f47ac10b-58cc-4372-a567-0e02b2c3d479",
					TaxonomyId:      "e6fb062d-74d6-4491-80bc-5d2c8e6d9ebb",
					ModuleName:      "module-1",
					Source:          "module-1-source-1",
					Version:         "1.1",
					Description:     "this is first test module",
					InputAttributes: []*terrariumpb.ModuleAttribute{},
					Namespace:       "unit-test",
				},
			},
		},
		{
			name: "populate mappings",
			filters: []db.FilterOption{
				db.ModuleByIDsFilter(modules[1].ID),
				db.PopulateModuleMappingsFilter(true),
			},
			wantModule: []*terrariumpb.Module{
				{
					Id:          "d3c1d35c-47a9-4837-add4-0e30db0f2f3b",
					TaxonomyId:  "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
					ModuleName:  "module-2",
					Source:      "module-2-source-1",
					Version:     "1.1",
					Description: "this is second test module",
					InputAttributes: []*terrariumpb.ModuleAttribute{
						{
							Name:        "module-2-attr-1",
							Description: "first attribute of the second module",
							OutputModuleAttributes: []*terrariumpb.ModuleAttribute{
								{
									Name:        "module-1-attr-3",
									Description: "first output attribute of the first module",
									ParentModule: &terrariumpb.Module{
										Id:              "f47ac10b-58cc-4372-a567-0e02b2c3d479",
										TaxonomyId:      "e6fb062d-74d6-4491-80bc-5d2c8e6d9ebb",
										ModuleName:      "module-1",
										Source:          "module-1-source-1",
										Version:         "1.1",
										Description:     "this is first test module",
										InputAttributes: []*terrariumpb.ModuleAttribute{},
										Namespace:       "unit-test",
									},
								},
							},
						},
					},
					Namespace: "unit-test",
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
					gotResult, err := dbObj.QueryTFModules(tt.filters...)
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
