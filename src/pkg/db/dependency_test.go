// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

//go:build dbtest
// +build dbtest

package db_test

// func Test_gDB_QueryDependencies(t *testing.T) {
// 	tests := []struct {
// 		name     string
// 		filters  []db.FilterOption
// 		wantDeps []*terrariumpb.Dependency
// 		wantErr  bool
// 		errMsg   string
// 	}{
// 		{
// 			name: "query by InterfaceID",
// 			filters: []db.FilterOption{
// 				db.DependencySearchFilter("server_web"),
// 			},
// 			wantDeps: []*terrariumpb.Dependency{
// 				{
// 					InterfaceId: "server_web",
// 					Title:       "Web Server",
// 					Description: "A server that hosts web applications and handles HTTP requests.",
// 					Inputs: &terrariumpb.JSONSchema{
// 						Type: "object",
// 						Properties: map[string]*terrariumpb.JSONSchema{
// 							"port": {
// 								Title:       "Port",
// 								Description: "The port number on which the server should listen.",
// 								Type:        "number",
// 								Default:     structpb.NewStringValue("80"),
// 							},
// 						},
// 					},
// 					Outputs: &terrariumpb.JSONSchema{
// 						Properties: map[string]*terrariumpb.JSONSchema{
// 							"host": {
// 								Title:       "Host",
// 								Description: "The host address of the web server.",
// 								Type:        "string",
// 							},
// 						},
// 					},
// 				},
// 			},
// 		},
// 	}

// 	for dbName, connector := range getConnectorMap() {
// 		g := connector(t)
// 		dbObj, err := db.AutoMigrate(g)
// 		require.NoError(t, err)

// 		t.Run(dbName, func(t *testing.T) {
// 			for _, tt := range tests {
// 				t.Run(tt.name, func(t *testing.T) {
// 					gotDeps, err := dbObj.QueryDependencies(tt.filters...)
// 					if tt.wantErr {
// 						assert.Error(t, err)

// 						if tt.errMsg != "" {
// 							assert.EqualError(t, err, tt.errMsg)
// 						}
// 					} else {
// 						assert.NoError(t, err)
// 						assert.EqualValues(t, tt.wantDeps, gotDeps.ToProto())
// 					}
// 				})
// 			}
// 		})
// 	}
// }

// func mustNewValue(v interface{}) *structpb.Value {
// 	value, err := structpb.NewValue(v)
// 	if err != nil {
// 		panic(fmt.Sprintf("Failed to create proto value: %v", err))
// 	}
// 	return value
// }
