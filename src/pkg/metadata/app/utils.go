// Copyright (c) Ollion
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"strings"

	"github.com/rotisserie/eris"
)

// Validate ensures that each application and its dependencies have unique IDs within the scope of all applications.
// It also checks that shared dependencies are provisioned.
func (apps Apps) Validate() error {
	seenAppIDs := make(map[string]struct{})   // Tracks observed application IDs for uniqueness.
	seenDepIDs := make(map[string]struct{})   // Tracks observed dependency IDs for uniqueness.
	sharedDepIDs := make(map[string]struct{}) // Tracks IDs of dependencies marked as shared.

	for _, app := range apps {
		// Validate each application and its dependencies.
		if err := app.validate(seenAppIDs, seenDepIDs, sharedDepIDs); err != nil {
			return err
		}
	}

	// Ensure that each shared dependency is provisioned.
	for depID := range sharedDepIDs {
		if _, exists := seenDepIDs[depID]; !exists {
			return eris.New("shared dependency not provisioned: " + depID)
		}
	}

	return nil
}

// Validate ensures that the application and its dependencies have unique IDs within the scope of the given applications.
func (app *App) Validate() error {
	return app.validate(map[string]struct{}{}, map[string]struct{}{}, map[string]struct{}{})
}

// validate checks the uniqueness of the app's ID and its dependencies' IDs.
func (app *App) validate(seenAppIDs, seenDepIDs, sharedDepIDs map[string]struct{}) error {
	if app.ID == "" {
		return eris.New("app id must not be empty")
	}

	// Ensure the app's ID is unique.
	if _, exists := seenAppIDs[app.ID]; exists {
		return eris.New("duplicate app ID: " + app.ID)
	}
	seenAppIDs[app.ID] = struct{}{}

	for _, dep := range app.GetDependencies() {
		// Validate each dependency within the context of the app.
		if err := dep.validate(seenDepIDs, sharedDepIDs); err != nil {
			return err
		}
	}

	return nil
}

// Validate if the dependency ID is set
func (dep *Dependency) Validate() error {
	return dep.validate(map[string]struct{}{}, map[string]struct{}{})
}

// validate checks the uniqueness of a dependency's ID and tracks shared dependencies.
func (dep *Dependency) validate(seenDepIDs, sharedDepIDs map[string]struct{}) error {
	if dep.ID == "" {
		return eris.Errorf("dependency `id` field must not be empty for: %s", dep.Use)
	}

	if dep.Use == "" {
		return eris.Errorf("dependency `use` field must not be empty for: %s", dep.ID)
	}

	// Ensure the dependency's ID is unique unless it's marked as NoProvision.
	if _, exists := seenDepIDs[dep.ID]; !dep.NoProvision && exists {
		return eris.New("duplicate dependency ID: " + dep.ID)
	}

	// Track the dependency ID as either shared or seen based on its provisioning status.
	if dep.NoProvision {
		sharedDepIDs[dep.ID] = struct{}{}
	} else {
		seenDepIDs[dep.ID] = struct{}{}
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
