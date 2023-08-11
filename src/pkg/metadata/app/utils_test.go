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
					Use: depPostgres,
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

	apps.SetDefaults()
	for _, app := range apps {
		assert.Equal(t, strings.ToUpper(app.ID), app.EnvPrefix)
		assert.Equal(t, app.ID, app.Compute.ID)
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

	deps := apps.GetDependenciesByType(depPostgres)
	assert.Len(t, deps, 2)
	assert.Equal(t, "testdep2", deps[0].ID)
	assert.Equal(t, "testdep4", deps[1].ID)
}

func TestGetUniqueDependencyTypes(t *testing.T) {
	apps := getAppsTest()

	types := apps.GetUniqueDependencyTypes()
	assert.Len(t, types, 3)
	assert.Contains(t, types, depWeb)
	assert.Contains(t, types, depRedis)
	assert.Contains(t, types, depPostgres)
}
