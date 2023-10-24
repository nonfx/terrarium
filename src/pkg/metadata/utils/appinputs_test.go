// Copyright (c) Ollion
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"testing"

	"github.com/cldcvr/terrarium/src/pkg/jsonschema"
	"github.com/cldcvr/terrarium/src/pkg/metadata/app"
	"github.com/cldcvr/terrarium/src/pkg/metadata/platform"
	"github.com/stretchr/testify/assert"
	"github.com/xeipuuv/gojsonschema"
)

func TestMatchAppAndPlatform(t *testing.T) {
	tests := []struct {
		name     string
		pm       *platform.PlatformMetadata
		apps     app.Apps
		wantErr  bool
		errMsg   string
		wantApps app.Apps
	}{
		{
			name: "Component not found in platform",
			pm: &platform.PlatformMetadata{
				Components: platform.Components{},
			},
			apps: app.Apps{
				app.App{
					ID: "testApp",
					Dependencies: app.Dependencies{
						app.Dependency{
							ID:  "testDep",
							Use: "missingComp",
						},
					},
				},
			},
			wantErr: true,
			errMsg:  "component 'testDep.missingComp' is not implemented in the platform",
		},
		{
			name: "Invalid input",
			pm: &platform.PlatformMetadata{
				Components: platform.Components{
					{
						ID: "comp1",
						Inputs: &jsonschema.Node{
							Type: gojsonschema.TYPE_OBJECT,
							Properties: map[string]*jsonschema.Node{
								"input1": {
									Type: gojsonschema.TYPE_NUMBER,
								},
							},
						},
					},
				},
			},
			apps: app.Apps{
				app.App{
					ID: "testApp",
					Dependencies: app.Dependencies{
						app.Dependency{
							ID:  "testDep",
							Use: "comp1",
							Inputs: map[string]interface{}{
								"input1": "not number",
							},
						},
					},
				},
			},
			wantErr: true,
			errMsg:  "component 'testDep.comp1' does not contain a valid set of inputs: validation failed with following errors: \n\tinput1: Invalid type. Expected: number, given: string",
		},
		{
			name: "Success set defaults",
			pm: &platform.PlatformMetadata{
				Components: platform.Components{
					{
						ID: "comp1",
						Inputs: &jsonschema.Node{
							Type: gojsonschema.TYPE_OBJECT,
							Properties: map[string]*jsonschema.Node{
								"input1": {
									Type:    gojsonschema.TYPE_NUMBER,
									Default: 10,
								},
								"input2": {
									Type:    gojsonschema.TYPE_STRING,
									Default: "val2",
								},
							},
						},
					},
				},
			},
			apps: app.Apps{
				app.App{
					ID: "testApp",
					Compute: app.Dependency{
						ID:     "testComp",
						Use:    "comp1",
						Inputs: map[string]interface{}{},
					},
					Dependencies: app.Dependencies{
						app.Dependency{
							ID:  "testDep",
							Use: "comp1",
							Inputs: map[string]interface{}{
								"input1": 20,
							},
						},
					},
				},
			},
			wantApps: app.Apps{
				app.App{
					ID: "testApp",
					Compute: app.Dependency{
						ID:  "testComp",
						Use: "comp1",
						Inputs: map[string]interface{}{
							"input1": 10,
							"input2": "val2",
						},
					},
					Dependencies: app.Dependencies{
						app.Dependency{
							ID:  "testDep",
							Use: "comp1",
							Inputs: map[string]interface{}{
								"input1": 20,
								"input2": "val2",
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := MatchAppAndPlatform(tt.pm, tt.apps)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantApps, tt.apps)
			}
		})
	}
}
