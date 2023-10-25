// Copyright (c) Ollion
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

func Test_gDB_QueryPlatformComponents(t *testing.T) {
	emptyJsonSchema := &terrariumpb.JSONSchema{
		Type:       gojsonschema.TYPE_OBJECT,
		Properties: map[string]*terrariumpb.JSONSchema{},
	}
	tests := []struct {
		name       string
		filters    []db.FilterOption
		validator  func(*testing.T, db.PlatformComponents)
		wantModule []*terrariumpb.Component //DependencyOutputs
		wantErr    bool
	}{
		{
			name: "query by InterfaceID",
			filters: db.ComponentRequestToFilters(&terrariumpb.ListComponentsRequest{
				Search: "ependency-2-int",
			}),
			wantModule: []*terrariumpb.Component{
				{
					Id:            uuidPlat1Comp2.String(),
					InterfaceUuid: uuidDep2.String(),
					InterfaceId:   "dependency-2-interface",
					Title:         "dependency-2",
					Description:   "this is second test dependency",
					Taxonomy:      []string{"mockdata-l1", "mockdata-l2", "mockdata-l3.2", "mockdata-l4.2", "mockdata-l5.2", "mockdata-l6.2", "mockdata-l7.2"},
					Inputs:        emptyJsonSchema,
					Outputs:       emptyJsonSchema,
				},
			},
		},
		{
			name: "query by taxonomy",
			filters: db.ComponentRequestToFilters(&terrariumpb.ListComponentsRequest{
				Taxonomy: "mockdata-l1/mockdata-l2/mockdata-l3.2",
			}),
			wantModule: []*terrariumpb.Component{
				{
					Id:            uuidPlat1Comp2.String(),
					InterfaceUuid: uuidDep2.String(),
					InterfaceId:   "dependency-2-interface",
					Title:         "dependency-2",
					Description:   "this is second test dependency",
					Taxonomy:      []string{"mockdata-l1", "mockdata-l2", "mockdata-l3.2", "mockdata-l4.2", "mockdata-l5.2", "mockdata-l6.2", "mockdata-l7.2"},
					Inputs:        emptyJsonSchema,
					Outputs:       emptyJsonSchema,
				},
			},
		},
		{
			name: "query by platform",
			filters: db.ComponentRequestToFilters(&terrariumpb.ListComponentsRequest{
				PlatformId: uuidPlat2.String(),
			}),
			wantModule: []*terrariumpb.Component{
				{
					Id:            uuidPlat2Comp1.String(),
					InterfaceUuid: uuidDep1.String(),
					InterfaceId:   "dependency-1-interface",
					Title:         "dependency-1",
					Description:   "this is first test dependency",
					Taxonomy:      []string{"mockdata-l1", "mockdata-l2", "mockdata-l3", "mockdata-l4", "mockdata-l5", "mockdata-l6", "mockdata-l7"},
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
			filters: db.ComponentRequestToFilters(&terrariumpb.ListComponentsRequest{
				Page: &terrariumpb.Page{Size: 2},
			}),
			validator: func(t *testing.T, d db.PlatformComponents) {
				assert.Equal(t, len(d), 2, "length of returned results")
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
					gotResult, err := dbObj.QueryPlatformComponents(tt.filters...)
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
