// Copyright (c) Ollion
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"github.com/cldcvr/terrarium/src/pkg/metadata/app"
	"github.com/cldcvr/terrarium/src/pkg/metadata/platform"
	pkgutils "github.com/cldcvr/terrarium/src/pkg/utils"
	"github.com/rotisserie/eris"
)

// MatchAppAndPlatform validate app dependency inputs and set default inputs
func MatchAppAndPlatform(pm *platform.PlatformMetadata, apps app.Apps) (err error) {
	for _, app := range pkgutils.ToRefArr(apps) {
		deps := pkgutils.ToRefArr(app.Dependencies)
		if app.Compute.ID != "" {
			deps = append(deps, &app.Compute)
		}

		for _, dep := range deps {
			if !dep.NoProvision {
				err = validateDependency(pm, dep)
				if err != nil {
					return
				}
			}
		}
	}

	return nil
}

// validateDependency validate and set default values to the dependency inputs.
func validateDependency(pm *platform.PlatformMetadata, appDep *app.Dependency) error {
	comp := pm.Components.GetByID(appDep.Use)
	if comp == nil {
		return eris.Errorf("component '%s.%s' is not implemented in the platform", appDep.ID, appDep.Use)
	}

	if appDep.Inputs == nil {
		appDep.Inputs = map[string]interface{}{}
	}

	err := comp.Inputs.Validate(appDep.Inputs)
	if err != nil {
		return eris.Wrapf(err, "component '%s.%s' does not contain a valid set of inputs", appDep.ID, appDep.Use)
	}

	comp.Inputs.ApplyDefaultsToMSI(appDep.Inputs)
	return nil
}
