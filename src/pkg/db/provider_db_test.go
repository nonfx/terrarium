// Copyright (c) Ollion
// SPDX-License-Identifier: Apache-2.0

//go:build dbtest
// +build dbtest

package db_test

import (
	"testing"

	"github.com/cldcvr/terrarium/src/pkg/db"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_gDB_CreateTFProvider(t *testing.T) {
	tests := []struct {
		name    string
		e       *db.TFProvider
		wantErr bool
	}{
		{
			name: "first new insert",
			e:    &db.TFProvider{Name: "mock-aws"},
		},
		{
			name: "redundant insert",
			e:    &db.TFProvider{Name: "mock-aws"},
		},
	}

	for dbName, connector := range getConnectorMap() {
		dbObj, err := db.AutoMigrate(connector(t))
		require.NoError(t, err)

		t.Run(dbName, func(t *testing.T) {
			moduleIDByNames := map[string]uuid.UUID{}
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					id, err := dbObj.CreateTFProvider(tt.e)
					if tt.wantErr {
						assert.Error(t, err)
					} else {
						assert.NoError(t, err)
						uniqueFieldsJoined := tt.e.Name
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

func Test_gDB_GetOrCreateTFProvider(t *testing.T) {
	uniqueName1 := uuid.New().String()
	uniqueName2 := uuid.New().String()
	tests := []struct {
		name    string
		e       *db.TFProvider
		wantErr bool
	}{
		{
			name: "first new insert",
			e:    &db.TFProvider{Name: uniqueName1},
		},
		{
			name: "redundant insert",
			e:    &db.TFProvider{Name: uniqueName1},
		},
		{
			name: "second new insert",
			e:    &db.TFProvider{Name: uniqueName2},
		},
	}

	for dbName, connector := range getConnectorMap() {
		dbObj, err := db.AutoMigrate(connector(t))
		require.NoError(t, err)

		t.Run(dbName, func(t *testing.T) {
			moduleIDByNames := map[string]uuid.UUID{}
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					id, isNew, err := dbObj.GetOrCreateTFProvider(tt.e)
					if tt.wantErr {
						assert.Error(t, err)
					} else {
						assert.NoError(t, err)
						uniqueFieldsJoined := tt.e.Name
						if wantID, ok := moduleIDByNames[uniqueFieldsJoined]; ok {
							assert.Equal(t, wantID, id)
							assert.False(t, isNew, "isNew should be false")
						} else {
							assert.NotEqual(t, uuid.Nil, id)
							assert.True(t, isNew, "isNew should be true")
							moduleIDByNames[uniqueFieldsJoined] = id
						}
					}
				})
			}
		})
	}
}
