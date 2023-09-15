// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"strings"

	"github.com/rotisserie/eris"
)

// Validates the uniqueness of application and dependency IDs across all apps.
func (apps Apps) Validate() error {
	seenAppIDs := make(map[string]struct{})
	seenDepIDs := make(map[string]struct{})
	sharedDepIDs := make(map[string]struct{})

	for _, app := range apps {
		if app.ID == "" {
			return eris.New("app id must not be empty")
		}

		// App ID must be unique across all apps
		if _, exists := seenAppIDs[app.ID]; exists {
			return eris.New("duplicate app ID: " + app.ID)
		}
		seenAppIDs[app.ID] = struct{}{}

		for i, dep := range app.GetDependencies() {
			if dep.ID == "" {
				return eris.Errorf("dependency id must not be empty for: %s (%d)", dep.Use, i)
			}

			// Dependency ID to provision must be unique in a project
			if _, exists := seenDepIDs[dep.ID]; !dep.NoProvision && exists {
				return eris.New("duplicate dependency ID: " + dep.ID)
			}

			if dep.NoProvision {
				sharedDepIDs[dep.ID] = struct{}{}
			} else {
				seenDepIDs[dep.ID] = struct{}{}
			}
		}
	}

	// shared dependency must have a provisioned instance with same ID.
	for depID := range sharedDepIDs {
		if _, exists := seenDepIDs[depID]; !exists {
			return eris.New("shared dependency not provisioned: " + depID)
		}
	}

	return nil
}

// Sets the default values for the optional fields in the Apps and its Dependencies.
func (apps *Apps) SetDefaults() {
	for i := range *apps {
		(*apps)[i].SetDefaults()
	}
}

// Sets the default values for the optional fields in the App and its Dependencies.
func (app *App) SetDefaults() {
	if app.EnvPrefix == "" {
		app.EnvPrefix = strings.ToUpper(app.ID)
	}

	if app.Compute.ID == "" && app.Compute.Use != "" {
		app.Compute.ID = app.ID
	}

	app.Compute.SetDefaults()

	for i := range app.Dependencies {
		app.Dependencies[i].SetDefaults()
	}
}

func (dep *Dependency) SetDefaults() {
	if dep.Inputs == nil {
		dep.Inputs = map[string]interface{}{}
	}

	if strings.Contains(dep.Use, "@") {
		split := strings.SplitN(dep.Use, "@", 2)
		dep.Use, dep.Inputs["version"] = split[0], split[1]
	}

	if dep.ID == "" {
		dep.ID = dep.Use
	}

	if dep.EnvPrefix == "" {
		dep.EnvPrefix = strings.ToUpper(dep.ID)
	}
}

// GetAppByID returns the first app found with the given ID, or nil if no such app exists.
func (apps Apps) GetAppByID(id string) *App {
	for _, app := range apps {
		if app.ID == id {
			return &app
		}
	}

	return nil
}

// GetDependencies returns the dependencies for the app including it's deployment dependency.
func (app App) GetDependencies() Dependencies {
	allDeps := app.Dependencies
	if app.Compute.ID != "" {
		allDeps = append(app.Dependencies, app.Compute)
	}

	return allDeps
}

// GetDependencies returns the dependencies for the app that needs to be provisioned including it's deployment dependency.
func (allDeps Dependencies) GetDependenciesToProvision() Dependencies {
	filteredDeps := make(Dependencies, 0, len(allDeps))

	for _, dep := range allDeps {
		if !dep.NoProvision {
			filteredDeps = append(filteredDeps, dep)
		}
	}

	return filteredDeps
}

// GetDependenciesByAppID returns the dependencies for the app with the given ID, or nil if the app does not exist.
func (apps Apps) GetDependenciesByAppID(appID string) Dependencies {
	app := apps.GetAppByID(appID)
	if app != nil {
		return app.GetDependencies()
	}

	return nil
}

// GetDependenciesByType returns all dependencies of a given type across all apps.
func (apps Apps) GetDependenciesByType(depType string) Dependencies {
	var deps Dependencies
	for _, app := range apps {
		for _, dep := range app.GetDependencies() {
			if dep.Use == depType && !dep.NoProvision {
				deps = append(deps, dep)
			}
		}
	}

	return deps
}

// GetUniqueDependencyTypes returns a list of unique dependency types across all apps.
func (apps Apps) GetUniqueDependencyTypes() []string {
	seenTypes := make(map[string]struct{})
	for _, app := range apps {
		for _, dep := range app.GetDependencies() {
			seenTypes[dep.Use] = struct{}{}
		}
	}

	var types []string
	for depType := range seenTypes {
		types = append(types, depType)
	}

	return types
}

// GetInputs returns a map of all inputs keyed by dependency identifier
func (allDeps Dependencies) GetInputs() map[string]interface{} {
	result := make(map[string]interface{}, len(allDeps))
	for _, dep := range allDeps {
		result[dep.ID] = dep.Inputs
	}

	return result
}
