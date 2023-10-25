// Copyright (c) Ollion
// SPDX-License-Identifier: Apache-2.0

package platform

import (
	"testing"

	"github.com/cldcvr/terraform-config-inspect/tfconfig"
	"github.com/stretchr/testify/assert"
)

func TestProfiles_Parse(t *testing.T) {
	type args struct {
		platformModule *tfconfig.Module
	}
	tests := []struct {
		name              string
		initialCollection *Profiles
		args              args
		wantProfiles      Profiles
	}{
		{
			name:              "load to empty collection",
			initialCollection: &Profiles{},
			args: args{
				platformModule: &tfconfig.Module{
					Path: "./test-component",
				},
			},
			wantProfiles: Profiles{
				Profile{
					ID:          "dev",
					Title:       "development profile",
					Description: "Development configuration profile.",
				},
				Profile{
					ID:          "prod",
					Title:       "",
					Description: "",
				},
			},
		},
		{
			name: "load to existing collection",
			initialCollection: &Profiles{
				Profile{
					ID:          "dev",
					Title:       "",
					Description: "",
				},
				Profile{
					ID:          "prod",
					Title:       "",
					Description: "",
				},
			},
			args: args{
				platformModule: &tfconfig.Module{
					Path: "./test-component",
				},
			},
			wantProfiles: Profiles{
				Profile{
					ID:          "dev",
					Title:       "development profile",
					Description: "Development configuration profile.",
				},
				Profile{
					ID:          "prod",
					Title:       "",
					Description: "",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.initialCollection.Parse(tt.args.platformModule)
			assert.ElementsMatch(t, tt.wantProfiles, *tt.initialCollection)
		})
	}
}
