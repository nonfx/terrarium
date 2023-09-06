// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package cli

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestFarmModuleList_Validate(t *testing.T) {
	tests := []struct {
		name    string
		list    FarmModuleList
		wantErr bool
	}{
		{
			name: "empty name",
			list: FarmModuleList{
				Farm: []FarmModuleRef{
					{
						Source:  "http://rowdy-watcher.info",
						Version: "4.5.6",
						Name:    "",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "empty source",
			list: FarmModuleList{
				Farm: []FarmModuleRef{
					{
						Source:  "",
						Version: "1.2.3",
						Name:    "solution",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "duplicate export name",
			list: FarmModuleList{
				Farm: []FarmModuleRef{
					{
						Source:  "http://rowdy-watcher.info",
						Version: "4.5.6",
						Name:    "solution",
					},
					{
						Source:  "https://knobby-courtroom.net",
						Version: "1.2.3",
						Name:    "solution",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "duplicate reference",
			list: FarmModuleList{
				Farm: []FarmModuleRef{
					{
						Source:  "http://bite-sized-scrap.org",
						Version: "9.4.6",
						Name:    "Wooden",
					},
					{
						Source:  "http://bite-sized-scrap.org",
						Version: "9.4.6",
						Name:    "navigate",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "valid module list",
			list: FarmModuleList{
				Farm: []FarmModuleRef{
					{
						Source:  "http://defiant-forum.com",
						Version: "0.0.0",
						Name:    "Concrete",
					},
					{
						Source:  "https://heavy-caviar.name",
						Version: "21.5.4",
						Name:    "synthesizing",
					},
					{
						Source:  "http://cultured-subscription.com",
						Version: "8.8.8",
						Name:    "generating",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.list.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("FarmModuleList.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFarmModuleRef_WriteFile(t *testing.T) {
	tests := []struct {
		name         string
		m            FarmModuleRef
		wantDirPath  string
		wantFilePath string
		wantErr      bool
	}{
		{
			name: "invalid name",
			m: FarmModuleRef{
				Name:    "Home/",
				Source:  "Carolina",
				Version: "transmit",
				Export:  true,
			},
			wantErr: true,
		},
		{
			name: "valid reference",
			m: FarmModuleRef{
				Name:    "Home",
				Source:  "Carolina",
				Version: "transmit",
				Export:  true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDirPath, gotFilePath, err := tt.m.CreateTerraformFile()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, gotFilePath)
				assert.True(t, strings.HasPrefix(gotFilePath, gotDirPath))
			}
		})
	}
}

func TestLoadFarmModules(t *testing.T) {
	tests := []struct {
		name    string
		list    FarmModuleList
		wantErr bool
	}{
		{
			name: "empty list",
			list: FarmModuleList{
				Farm: []FarmModuleRef{},
			},
			wantErr: true,
		},
		{
			name: "invalid list",
			list: FarmModuleList{
				Farm: []FarmModuleRef{
					{
						Name:    "",
						Source:  "Automotive",
						Version: "cross-media",
						Export:  true,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "non-empty list",
			list: FarmModuleList{
				Farm: []FarmModuleRef{
					{
						Name:    "Wooden",
						Source:  "Automotive",
						Version: "cross-media",
						Export:  true,
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := LoadFarmModules("8065c60d-f7f8-4a25-b000-1a7578ac665e")
			assert.Error(t, err)

			fp, err := os.CreateTemp("", "*")
			if err != nil {
				t.Fatal(err)
			}
			b, err := yaml.Marshal(tt.list)
			if err != nil {
				t.Fatal(err)
			}
			fp.Write(b)

			got, err := LoadFarmModules(fp.Name())
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.list, got)
			}
		})
	}
}
