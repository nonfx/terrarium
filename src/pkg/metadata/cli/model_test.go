// Copyright (c) Ollion
// SPDX-License-Identifier: Apache-2.0

package cli

import (
	"os"
	"path"
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
	existingFile, err := os.CreateTemp("", "*")
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name         string
		m            FarmModuleRef
		root         string
		wantDirPath  string
		wantFilePath string
		wantErr      bool
	}{
		{
			name: "invalid name",
			root: "",
			m: FarmModuleRef{
				Name:    "Home/",
				Source:  "Carolina",
				Version: "transmit",
				Export:  true,
			},
			wantErr: true,
		},
		{
			name: "valid reference without root",
			root: "",
			m: FarmModuleRef{
				Name:    "Home",
				Source:  "Carolina",
				Version: "transmit",
				Export:  true,
			},
			wantErr: false,
		},
		{
			name: "non-existing root dir",
			root: "/etc/defaults/f0bc9575-5384-473c-a4fd-98a5ec9f3a86",
			m: FarmModuleRef{
				Name:    "qui",
				Source:  "uniform",
				Version: "bandwidth",
				Export:  true,
			},
			wantErr: true,
		},
		{
			name: "invalid root dir",
			root: existingFile.Name(),
			m: FarmModuleRef{
				Name:    "mobile",
				Source:  "Solomon",
				Version: "wireless",
				Export:  true,
			},
			wantErr: true,
		},
		{
			name: "invalid work module dir",
			root: path.Dir(existingFile.Name()),
			m: FarmModuleRef{
				Name:    path.Base(existingFile.Name()),
				Source:  "Solomon",
				Version: "wireless",
				Export:  true,
			},
			wantErr: true,
		},
		{
			name: "invalid work module dir name",
			root: os.TempDir(),
			m: FarmModuleRef{
				Name:    "Qui culpa ad. Saepe voluptatum et earum rem at officiis. Nihil voluptas earum vel accusantium qui. Rerum voluptatem corrupti necessitatibus dolores non. Ab quo alias et saepe quia quia similique sunt. Esse ut nihil aperiam. Qui nobis similique voluptates repellat. Enim laudantium qui quae eos. Sed voluptatem quia unde nemo nisi. Officia tempora blanditiis est quas soluta. Aliquam unde necessitatibus fugiat et culpa quidem ut illum adipisci.",
				Source:  "Solomon",
				Version: "wireless",
				Export:  true,
			},
			wantErr: true,
		},
		{
			name: "valid reference with root",
			root: os.TempDir(),
			m: FarmModuleRef{
				Name:    "Plastic",
				Source:  "Bridge",
				Version: "tan",
				Export:  true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDirPath, gotFilePath, err := tt.m.CreateTerraformFile(tt.root)
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
