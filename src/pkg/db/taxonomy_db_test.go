// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

//go:build dbtest
// +build dbtest

package db_test

import (
	"strings"
	"testing"

	"github.com/cldcvr/terrarium/src/pkg/db"
	"github.com/cldcvr/terrarium/src/pkg/pb/terrariumpb"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_gDB_CreateTaxonomy(t *testing.T) {
	tests := []struct {
		name     string
		taxonomy *db.Taxonomy
		wantErr  bool
	}{
		{
			name:     "first new insert",
			taxonomy: db.TaxonomyFromLevels("mocktest-l1", "mocktest-l2", "mocktest-l3", "mocktest-l4", "mocktest-l5", "mocktest-l6", "mocktest-l7"),
		},
		{
			name:     "redundant insert",
			taxonomy: db.TaxonomyFromLevels("mocktest-l1", "mocktest-l2", "mocktest-l3", "mocktest-l4", "mocktest-l5", "mocktest-l6", "mocktest-l7"),
		},
		{
			name:     "one field different",
			taxonomy: db.TaxonomyFromLevels("mocktest-l1", "mocktest-l2", "mocktest-l3", "mocktest-l4", "mocktest-l5", "mocktest-l6", "mocktest-l7-2"),
		},
		{
			name:     "new insert",
			taxonomy: db.TaxonomyFromLevels("mocktest-l1-3", "mocktest-l2-3", "mocktest-l3-3", "mocktest-l4-3", "mocktest-l5-3", "mocktest-l6-3", "mocktest-l7-3"),
		},
	}

	for dbName, connector := range getConnectorMap() {
		dbObj, err := db.AutoMigrate(connector(t))
		require.NoError(t, err)

		t.Run(dbName, func(t *testing.T) {
			moduleIDByNames := map[string]uuid.UUID{}
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					id, err := dbObj.CreateTaxonomy(tt.taxonomy)
					if tt.wantErr {
						assert.Error(t, err)
					} else {
						assert.NoError(t, err)
						uniqueFieldsJoined := strings.Join(tt.taxonomy.ToLevels(), "|")
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

func Test_gDB_QueryTaxonomies(t *testing.T) {
	tests := []struct {
		name       string
		filterOps  *terrariumpb.ListTaxonomyRequest
		validator  func(*testing.T, db.Taxonomies)
		wantResult []*terrariumpb.Taxonomy
		wantErr    bool
	}{
		{
			name: "get by common top two levels",
			filterOps: &terrariumpb.ListTaxonomyRequest{
				Taxonomy: "mockdata-l1/mockdata-l2",
				Page:     &terrariumpb.Page{Size: 5},
			},
			wantResult: []*terrariumpb.Taxonomy{
				{
					Id:     uuidTax1.String(),
					Levels: []string{"mockdata-l1", "mockdata-l2", "mockdata-l3", "mockdata-l4", "mockdata-l5", "mockdata-l6", "mockdata-l7"},
				},
				{
					Id:     uuidTax2.String(),
					Levels: []string{"mockdata-l1", "mockdata-l2", "mockdata-l3.2", "mockdata-l4.2", "mockdata-l5.2", "mockdata-l6.2", "mockdata-l7.2"},
				},
			},
		},
		{
			name: "get all",
			filterOps: &terrariumpb.ListTaxonomyRequest{
				Page: &terrariumpb.Page{Size: 5},
			},
			validator: func(t *testing.T, tax db.Taxonomies) {
				assert.GreaterOrEqual(t, len(tax), 2, "length of taxonomies array")
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
					gotResult, gotErr := dbObj.QueryTaxonomies(db.TaxonomyRequestToFilters(tt.filterOps)...)
					if tt.wantErr {
						assert.Error(t, gotErr)
						return
					}

					require.NoError(t, gotErr)
					if tt.validator != nil {
						tt.validator(t, gotResult)
						return
					}

					assert.EqualValues(t, tt.wantResult, gotResult.ToProto())
				})
			}
		})
	}
}
