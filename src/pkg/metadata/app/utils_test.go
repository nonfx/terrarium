package app

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// initializing test data
func getAppsTest() Apps {
	return Apps{
		{
			ID: "testapp1",
			Service: Dependency{
				ID:   "testapp1",
				Type: "service.web",
			},
			Dependencies: Dependencies{
				{
					ID:   "testdep2",
					Type: "database",
				},
				{
					ID:          "testdep3",
					Type:        "cache",
					NoProvision: true,
				},
			},
		},
		{
			ID: "testapp2",
			Service: Dependency{
				ID:   "testapp2",
				Type: "service.web",
			},
			Dependencies: Dependencies{
				{
					ID:   "testdep3",
					Type: "cache",
				},
				{
					ID:   "testdep4",
					Type: "database",
				},
				{
					ID:   "testdep5",
					Type: "cache",
				},
			},
		},
	}
}

func TestSetDefaults(t *testing.T) {
	apps := getAppsTest()

	apps[0].Service.ID = ""

	apps.SetDefaults()
	for _, app := range apps {
		assert.Equal(t, strings.ToUpper(app.ID), app.EnvPrefix)
		assert.Equal(t, app.ID, app.Service.ID)
		for _, dep := range app.Dependencies {
			assert.Equal(t, strings.ToUpper(dep.ID), dep.EnvPrefix)
		}
	}
}

func TestValidate(t *testing.T) {
	apps := getAppsTest()

	assert.Nil(t, apps.Validate())

	apps[0].ID = ""
	assert.ErrorContains(t, apps.Validate(), "app id must not be empty")

	apps = getAppsTest()
	apps[0].ID = apps[1].ID
	assert.ErrorContains(t, apps.Validate(), "duplicate app ID: testapp2")

	apps = getAppsTest()
	apps[0].Dependencies[0].ID = ""
	assert.ErrorContains(t, apps.Validate(), "dependency id must not be empty")

	apps = getAppsTest()
	apps[0].Dependencies[1].NoProvision = false
	assert.ErrorContains(t, apps.Validate(), "duplicate dependency ID: testdep3")

	apps = getAppsTest()
	apps = apps[:1]
	assert.ErrorContains(t, apps.Validate(), "shared dependency not provisioned: testdep3")
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

	deps := apps.GetDependenciesByType("database")
	assert.Len(t, deps, 2)
	assert.Equal(t, "testdep2", deps[0].ID)
	assert.Equal(t, "testdep4", deps[1].ID)
}

func TestGetUniqueDependencyTypes(t *testing.T) {
	apps := getAppsTest()

	types := apps.GetUniqueDependencyTypes()
	assert.Len(t, types, 3)
	assert.Contains(t, types, "service.web")
	assert.Contains(t, types, "cache")
	assert.Contains(t, types, "database")
}
