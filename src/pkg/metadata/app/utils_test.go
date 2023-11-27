// Copyright (c) Ollion
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	depWeb      = "compute_web"
	depPostgres = "postgres"
	depRedis    = "redis"
)

// initializing test data
func getAppsTest() Apps {
	return Apps{
		{
			ID: "testapp1",
			Compute: Dependency{
				ID:  "testapp1",
				Use: depWeb,
			},
			Dependencies: Dependencies{
				{
					ID:  "testdep2",
					Use: depPostgres + "@11",
				},
				{
					ID:          "testdep3",
					Use:         depRedis,
					NoProvision: true,
				},
			},
		},
		{
			ID: "testapp2",
			Compute: Dependency{
				ID:  "testapp2",
				Use: depWeb,
			},
			Dependencies: Dependencies{
				{
					ID:  "testdep3",
					Use: depRedis,
				},
				{
					ID:  "testdep4",
					Use: depPostgres,
				},
				{
					ID:  "testdep5",
					Use: depRedis,
				},
			},
		},
	}
}

func TestSetDefaults(t *testing.T) {
	apps := getAppsTest()

	apps[0].Compute.ID = ""
	apps[0].Dependencies[0].ID = ""

	apps.SetDefaults()
	for _, app := range apps {
		assert.Equal(t, strings.ToUpper(app.ID), app.EnvPrefix)
		assert.Equal(t, app.ID, app.Compute.ID)
		for _, dep := range app.Dependencies {
			assert.Equal(t, strings.ToUpper(dep.ID), dep.EnvPrefix)
		}
	}

	assert.Equal(t, apps[0].Dependencies[0].ID, depPostgres)
	assert.Equal(t, apps[0].Dependencies[0].Use, depPostgres)
	assert.Equal(t, apps[0].Dependencies[0].Inputs["version"], "11")
}

func TestAppsValidate(t *testing.T) {
	tests := []struct {
		name        string
		apps        Apps
		expectError string
	}{
		{
			name: "Valid Apps",
			apps: func() Apps {
				return getAppsTest()
			}(),
		},
		{
			name: "Empty App ID",
			apps: func() Apps {
				apps := getAppsTest()
				apps[0].ID = ""
				return apps
			}(),
			expectError: "app id must not be empty",
		},
		{
			name: "Duplicate App ID",
			apps: func() Apps {
				apps := getAppsTest()
				apps[0].ID = apps[1].ID
				return apps
			}(),
			expectError: "duplicate app ID: testapp2",
		},
		{
			name: "Empty Dependency ID",
			apps: func() Apps {
				apps := getAppsTest()
				apps[0].Dependencies[0].ID = ""
				return apps
			}(),
			expectError: "dependency `id` field must not be empty",
		},
		{
			name: "Duplicate Dependency ID",
			apps: func() Apps {
				apps := getAppsTest()
				apps[0].Dependencies[1].NoProvision = false
				return apps
			}(),
			expectError: "duplicate dependency ID: testdep3",
		},
		{
			name: "Un-provisioned Shared Dependency",
			apps: func() Apps {
				apps := getAppsTest()
				apps = apps[:1]
				return apps
			}(),
			expectError: "shared dependency not provisioned: testdep3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.apps.Validate()

			if tt.expectError != "" {
				assert.ErrorContains(t, err, tt.expectError)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestAppValidate(t *testing.T) {
	tests := []struct {
		name        string
		app         App
		expectError string
	}{
		{
			name: "Valid App",
			app:  getAppsTest()[1],
		},
		{
			name: "Valid App with NoProvision",
			app:  getAppsTest()[0],
		},
		{
			name: "Empty App ID",
			app: func() App {
				app := getAppsTest()[0]
				app.ID = ""
				return app
			}(),
			expectError: "app id must not be empty",
		},
		{
			name: "Empty Dependency ID",
			app: func() App {
				app := getAppsTest()[0]
				app.Dependencies[0].ID = ""
				return app
			}(),
			expectError: "dependency `id` field must not be empty for",
		},
		{
			name: "Duplicate Dependency ID",
			app: func() App {
				app := getAppsTest()[1]
				app.Dependencies[0].ID = app.Dependencies[1].ID
				return app
			}(),
			expectError: "duplicate dependency ID: testdep4",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.app.Validate()

			if tt.expectError != "" {
				assert.ErrorContains(t, err, tt.expectError)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestDependencyValidate(t *testing.T) {
	tests := []struct {
		name    string
		dep     Dependency
		wantErr string
	}{
		{
			name: "Valid Dependency",
			dep: Dependency{
				ID:  "dep1",
				Use: "dep",
			},
		},
		{
			name: "Empty Dependency ID",
			dep: Dependency{
				ID:  "",
				Use: "dep",
			},
			wantErr: "dependency `id` field must not be empty",
		},
		{
			name: "Empty Dependency Use field",
			dep: Dependency{
				ID:  "dep1",
				Use: "",
			},
			wantErr: "dependency `use` field must not be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.dep.Validate()

			if tt.wantErr != "" {
				assert.ErrorContains(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetAppByID(t *testing.T) {
	apps := getAppsTest()

	app := apps.GetAppByID("testapp2")
	assert.NotNil(t, app)
	assert.Equal(t, "testapp2", app.ID)

	assert.Nil(t, apps.GetAppByID("invalid"))
}

func TestGetDependencies(t *testing.T) {
	apps := getAppsTest()

	app := apps.GetAppByID("testapp1")
	deps := app.GetDependencies()
	assert.Len(t, deps, 3)
	assert.Equal(t, "testdep2", deps[0].ID)
	assert.Equal(t, "testdep3", deps[1].ID)
	assert.Equal(t, "testapp1", deps[2].ID)
}

func TestGetDependenciesToProvision(t *testing.T) {
	apps := getAppsTest()

	app := apps.GetAppByID("testapp1")
	deps := app.GetDependencies().GetDependenciesToProvision()
	assert.Len(t, deps, 2)
	assert.Equal(t, "testdep2", deps[0].ID)
	assert.Equal(t, "testapp1", deps[1].ID)
}

func TestGetDependenciesByAppID(t *testing.T) {
	apps := getAppsTest()

	deps := apps.GetDependenciesByAppID("testapp1")
	assert.Len(t, deps, 3)
	assert.Equal(t, "testdep2", deps[0].ID)
	assert.Equal(t, "testdep3", deps[1].ID)
	assert.Equal(t, "testapp1", deps[2].ID)

	assert.Empty(t, apps.GetDependenciesByAppID("invalid"))
}

func TestGetDependenciesByType(t *testing.T) {
	apps := getAppsTest()
	apps.SetDefaults()

	deps := apps.GetDependenciesByType(depPostgres)
	assert.Len(t, deps, 2)
	assert.Equal(t, "testdep2", deps[0].ID)
	assert.Equal(t, "testdep4", deps[1].ID)
}

func TestGetUniqueDependencyTypes(t *testing.T) {
	apps := getAppsTest()
	apps.SetDefaults()

	types := apps.GetUniqueDependencyTypes()
	assert.Len(t, types, 3)
	assert.Contains(t, types, depWeb)
	assert.Contains(t, types, depRedis)
	assert.Contains(t, types, depPostgres)
}

func TestNewApp(t *testing.T) {
	tests := []struct {
		name    string
		content []byte
		wantErr bool
	}{
		{
			name:    "valid yaml content",
			content: []byte("\nid: test_app\nname: Test App\n"),
			wantErr: false,
		},
		{
			name:    "invalid yaml content",
			content: []byte("\nid: test_app\nname: Test App\ninvalidYAML\n"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, err := NewApp(tt.content)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, app)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, app)
			}
		})
	}
}

func TestGetInputs(t *testing.T) {
	apps := getAppsTest()
	apps.SetDefaults()

	output := apps.GetDependenciesByType(depPostgres).GetInputs()

	assert.Equal(t, map[string]interface{}{
		"testdep2": map[string]interface{}{
			"version": "11",
		},
		"testdep4": map[string]interface{}{},
	}, output)
}
