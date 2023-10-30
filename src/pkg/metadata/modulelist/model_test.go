// Copyright (c) Ollion
// SPDX-License-Identifier: Apache-2.0

package modulelist

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestFarmModuleList_Validate(t *testing.T) {
	tests := []struct {
		name    string
		modules FarmModuleList
		wantErr string
	}{
		{
			name: "valid modules with groups",
			modules: FarmModuleList{
				Farm: []FarmModuleRef{
					{Name: "Module1", Source: "Source1", Version: "1.0", Group: "Group1", Export: true},
					{Name: "Module2", Source: "Source2", Version: "1.0", Group: "Group1", Export: true},
				},
			},
		},
		{
			name: "valid module list",
			modules: FarmModuleList{
				Farm: []FarmModuleRef{
					{Source: "http://defiant-forum.com", Version: "0.0.0", Name: "Concrete"},
					{Source: "https://heavy-caviar.name", Version: "21.5.4", Name: "synthesizing"},
					{Source: "http://cultured-subscription.com", Version: "8.8.8", Name: "generating"},
				},
			},
		},
		{
			name: "duplicate module names",
			modules: FarmModuleList{
				Farm: []FarmModuleRef{
					{Name: "Module1", Source: "Source1", Version: "1.0"},
					{Name: "Module1", Source: "Source2", Version: "1.0"},
				},
			},
			wantErr: "module 'Module1' has a duplicate name",
		},
		{
			name: "duplicate module reference",
			modules: FarmModuleList{
				Farm: []FarmModuleRef{
					{Name: "Module1", Source: "Source1", Version: "1.0"},
					{Name: "Module2", Source: "Source1", Version: "1.0"},
				},
			},
			wantErr: "module 'Module2' is duplicate of module 'Module2'",
		},
		{
			name: "conflicting group name and module name",
			modules: FarmModuleList{
				Farm: []FarmModuleRef{
					{Name: "Module1", Source: "Source1", Version: "1.0", Group: "Group1"},
					{Name: "Module2", Source: "Source2", Version: "1.0", Group: "Module1"},
				},
			},
			wantErr: "group 'Module1' has name conflict with a module name",
		},
		{
			name: "conflicting module name and group name",
			modules: FarmModuleList{
				Farm: []FarmModuleRef{
					{Name: "Module1", Source: "Source1", Version: "1.0", Group: "Group1"},
					{Name: "Group1", Source: "Source2", Version: "1.0", Group: "Group1"},
				},
			},
			wantErr: "module 'Group1' has name conflict with a group name",
		},
		{
			name: "empty module name",
			modules: FarmModuleList{
				Farm: []FarmModuleRef{
					{Name: "", Source: "Source1", Version: "1.0"},
				},
			},
			wantErr: "module must have a non-empty name",
		},
		{
			name: "empty module source",
			modules: FarmModuleList{
				Farm: []FarmModuleRef{
					{Name: "Module1", Source: "", Version: "1.0"},
				},
			},
			wantErr: "module 'Module1' must have a source",
		},
		{
			name: "conflicting export flags",
			modules: FarmModuleList{
				Farm: []FarmModuleRef{
					{Name: "Module1", Source: "Source1", Version: "1.0", Group: "Group1", Export: true},
					{Name: "Module2", Source: "Source2", Version: "1.0", Group: "Group1", Export: false},
				},
			},
			wantErr: "group 'Group1' has conflicting export flag! all modules in a group must have same export value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.modules.Validate()
			if tt.wantErr != "" {
				assert.ErrorContains(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
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
		m            FarmModuleGroup
		root         string
		wantDirPath  string
		wantFilePath string
		wantErr      bool
	}{
		{
			name: "invalid name",
			root: "",
			m: FarmModuleGroup{
				Name: "Home/",
				Modules: []FarmModule{
					{
						Name:    "Home/",
						Source:  "Carolina",
						Version: "transmit",
					},
				},
				Export: true,
			},
			wantErr: true,
		},
		{
			name: "valid reference without root",
			root: "",
			m: FarmModuleGroup{
				Name: "Home",
				Modules: []FarmModule{
					{
						Name:    "Home",
						Source:  "Carolina",
						Version: "transmit",
					},
				},
				Export: true,
			},
			wantErr: false,
		},
		{
			name: "non-existing root dir", // creates new directory
			root: path.Join(t.TempDir(), "f0bc9575-5384-473c-a4fd-98a5ec9f3a86"),
			m: FarmModuleGroup{
				Name: "qui",
				Modules: []FarmModule{
					{
						Name:    "qui",
						Source:  "uniform",
						Version: "bandwidth",
					},
				},
				Export: true,
			},
			wantErr: false,
		},
		{
			name: "invalid root dir",
			root: existingFile.Name(),
			m: FarmModuleGroup{
				Name: "mobile",
				Modules: []FarmModule{
					{
						Name:    "mobile",
						Source:  "Solomon",
						Version: "wireless",
					},
				},
				Export: true,
			},
			wantErr: true,
		},
		{
			name: "invalid work module dir",
			root: path.Dir(existingFile.Name()),
			m: FarmModuleGroup{
				Name: path.Base(existingFile.Name()),
				Modules: []FarmModule{
					{
						Name:    path.Base(existingFile.Name()),
						Source:  "Solomon",
						Version: "wireless",
					},
				},
				Export: true,
			},
			wantErr: true,
		},
		{
			name: "invalid work module dir name",
			root: t.TempDir(),
			m: FarmModuleGroup{
				Name: "Qui culpa ad. Saepe voluptatum et earum rem at officiis. Nihil voluptas earum vel accusantium qui. Rerum voluptatem corrupti necessitatibus dolores non. Ab quo alias et saepe quia quia similique sunt. Esse ut nihil aperiam. Qui nobis similique voluptates repellat. Enim laudantium qui quae eos. Sed voluptatem quia unde nemo nisi. Officia tempora blanditiis est quas soluta. Aliquam unde necessitatibus fugiat et culpa quidem ut illum adipisci.",
				Modules: []FarmModule{
					{
						Name:    "Qui culpa ad. Saepe voluptatum et earum rem at officiis. Nihil voluptas earum vel accusantium qui. Rerum voluptatem corrupti necessitatibus dolores non. Ab quo alias et saepe quia quia similique sunt. Esse ut nihil aperiam. Qui nobis similique voluptates repellat. Enim laudantium qui quae eos. Sed voluptatem quia unde nemo nisi. Officia tempora blanditiis est quas soluta. Aliquam unde necessitatibus fugiat et culpa quidem ut illum adipisci.",
						Source:  "Solomon",
						Version: "wireless",
					},
				},
				Export: true,
			},
			wantErr: true,
		},
		{
			name: "valid reference with root",
			root: t.TempDir(),
			m: FarmModuleGroup{
				Name: "Plastic",
				Modules: []FarmModule{
					{
						Name:    "Plastic",
						Source:  "Bridge",
						Version: "tan",
					},
				},
				Export: true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDirPath, err := tt.m.CreateTerraformFile(tt.root)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, gotDirPath)
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

func TestFarmModuleList_Groups(t *testing.T) {
	tests := []struct {
		name  string
		list  FarmModuleList
		group FarmModuleGroupsMap
	}{
		{
			name: "Single Group",
			list: FarmModuleList{
				Farm: []FarmModuleRef{
					{Name: "Module1", Source: "Source1", Version: "1.0", Group: "Group1", Export: true},
					{Name: "Module2", Source: "Source2", Version: "1.0", Group: "Group1", Export: true},
				},
			},
			group: FarmModuleGroupsMap{
				"Group1": FarmModuleGroup{
					Name:   "Group1",
					Export: true,
					Modules: []FarmModule{
						{Name: "Module1", Source: "Source1", Version: "1.0"},
						{Name: "Module2", Source: "Source2", Version: "1.0"},
					},
				},
			},
		},
		{
			name: "Multiple Groups",
			list: FarmModuleList{
				Farm: []FarmModuleRef{
					{Name: "Module1", Source: "Source1", Version: "1.0", Group: "Group1", Export: true},
					{Name: "Module2", Source: "Source2", Version: "1.0", Group: "Group2", Export: false},
					{Name: "Module3", Source: "Source3", Version: "1.0", Export: false},
				},
			},
			group: FarmModuleGroupsMap{
				"Group1": FarmModuleGroup{
					Name:   "Group1",
					Export: true,
					Modules: []FarmModule{
						{Name: "Module1", Source: "Source1", Version: "1.0"},
					},
				},
				"Group2": FarmModuleGroup{
					Name:   "Group2",
					Export: false,
					Modules: []FarmModule{
						{Name: "Module2", Source: "Source2", Version: "1.0"},
					},
				},
				"Module3": FarmModuleGroup{
					Name:   "Module3",
					Export: false,
					Modules: []FarmModule{
						{Name: "Module3", Source: "Source3", Version: "1.0"},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.list.Groups()
			assert.Equal(t, tt.group, result)
		})
	}
}

func TestFarmModuleRef_GetGroupName(t *testing.T) {
	tests := []struct {
		name     string
		ref      FarmModuleRef
		expected string
	}{
		{
			name:     "Group Name Provided",
			ref:      FarmModuleRef{Name: "Module1", Group: "Group1"},
			expected: "Group1",
		},
		{
			name:     "Group Name Empty",
			ref:      FarmModuleRef{Name: "Module1", Group: ""},
			expected: "Module1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.ref.GetGroupName()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFarmModuleGroupsMap_FilterExport(t *testing.T) {
	tests := []struct {
		name        string
		groups      FarmModuleGroupsMap
		wantExport  bool
		expectedMap FarmModuleGroupsMap
	}{
		{
			name: "Filter Export True",
			groups: FarmModuleGroupsMap{
				"Group1": FarmModuleGroup{Export: true},
				"Group2": FarmModuleGroup{Export: false},
			},
			wantExport: true,
			expectedMap: FarmModuleGroupsMap{
				"Group1": FarmModuleGroup{Export: true},
			},
		},
		{
			name: "Filter Export False",
			groups: FarmModuleGroupsMap{
				"Group1": FarmModuleGroup{Export: true},
				"Group2": FarmModuleGroup{Export: false},
			},
			wantExport: false,
			expectedMap: FarmModuleGroupsMap{
				"Group2": FarmModuleGroup{Export: false},
			},
		},
		{
			name:        "Empty Groups",
			groups:      FarmModuleGroupsMap{},
			wantExport:  true,
			expectedMap: FarmModuleGroupsMap{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.groups.FilterExport(tt.wantExport)
			assert.Equal(t, tt.expectedMap, result)
		})
	}
}
