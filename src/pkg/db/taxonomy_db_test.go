// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

//go:build dbtest
// +build dbtest

package db_test

import (
	"strings"
	"testing"

	"github.com/cldcvr/terrarium/src/pkg/db"
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
			taxonomy: db.TaxonomyFromLevels("l1", "l2", "l3", "l4", "l5", "l6", "l7"),
		},
		{
			name:     "redundant insert",
			taxonomy: db.TaxonomyFromLevels("l1", "l2", "l3", "l4", "l5", "l6", "l7"),
		},
		{
			name:     "one field different",
			taxonomy: db.TaxonomyFromLevels("l1", "l2", "l3", "l4", "l5", "l6", "l7-2"),
		},
		{
			name:     "new insert",
			taxonomy: db.TaxonomyFromLevels("l1-3", "l2-3", "l3-3", "l4-3", "l5-3", "l6-3", "l7-3"),
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
