// Copyright (c) CloudCover
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

func Test_gDB_CreatePlatform(t *testing.T) {
	tests := []struct {
		name     string
		platform *db.Platform
		wantErr  bool
	}{
		{
			name:     "first new insert",
			platform: &db.Platform{Title: "test-1", RepoURL: "test-url", RepoDirectory: "test-dir", CommitSHA: "test-sha", RefLabel: "test-ref", LabelType: 1},
		},
	}

	for dbName, connector := range getConnectorMap() {
		dbObj, err := db.AutoMigrate(connector(t))
		require.NoError(t, err)

		t.Run(dbName, func(t *testing.T) {
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					id, err := dbObj.CreatePlatform(tt.platform)
					if tt.wantErr {
						assert.Error(t, err)
					} else {
						assert.NotEqual(t, uuid.Nil, id)
						assert.NoError(t, err)
					}
				})
			}
		})
	}
}

func Test_gDB_QueryPlatforms(t *testing.T) {
	tests := []struct {
		name         string
		filters      []db.FilterOption
		validator    func(*testing.T, db.Platforms)
		wantPlatform []*terrariumpb.Platform
		wantErr      bool
	}{
		{
			name: "pagination",
			filters: db.PlatformRequestToFilters(&terrariumpb.ListPlatformsRequest{
				Page: &terrariumpb.Page{Size: 2},
			}),
			validator: func(t *testing.T, p db.Platforms) {
				assert.Len(t, p, 2)
			},
		},
		{
			name: "search query",
			filters: db.PlatformRequestToFilters(&terrariumpb.ListPlatformsRequest{
				Search: "platform-1",
			}),
			wantPlatform: []*terrariumpb.Platform{
				{
					Id:         uuidPlat1.String(),
					Title:      "test-platform-1",
					RepoCommit: "2ed744403e50",
					Components: 2,
				},
			},
		},
		{
			name: "taxonomy query",
			filters: db.PlatformRequestToFilters(&terrariumpb.ListPlatformsRequest{
				Taxonomy: "mockdata-l1/mockdata-l2/mockdata-l3.2",
			}),
			wantPlatform: []*terrariumpb.Platform{
				{
					Id:         uuidPlat1.String(),
					Title:      "test-platform-1",
					RepoCommit: "2ed744403e50",
					Components: 2,
				},
			},
		},
		{
			name: "dependency query",
			filters: db.PlatformRequestToFilters(&terrariumpb.ListPlatformsRequest{
				InterfaceUuid: []string{uuidDep2.String()},
			}),
			wantPlatform: []*terrariumpb.Platform{
				{
					Id:         uuidPlat1.String(),
					Title:      "test-platform-1",
					RepoCommit: "2ed744403e50",
					Components: 2,
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
					gotResult, err := dbObj.QueryPlatforms(tt.filters...)
					if tt.wantErr {
						assert.Error(t, err)
						return
					}

					assert.NoError(t, err)

					if tt.validator != nil {
						tt.validator(t, gotResult)
					} else {
						assert.EqualValues(t, tt.wantPlatform, gotResult.ToProto())
					}
				})
			}
		})
	}
}
