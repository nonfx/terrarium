// Copyright (c) Ollion
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"errors"

	"github.com/cldcvr/terrarium/src/pkg/metadata/app"
	"github.com/cldcvr/terrarium/src/pkg/metadata/platform"
	pkgutils "github.com/cldcvr/terrarium/src/pkg/utils"
	"github.com/rotisserie/eris"
)

var (
	ErrComponentNotImplemented = eris.New("component is not implemented in the platform")
)

// MatchAppAndPlatform validate app dependency inputs and set default inputs
func MatchAppAndPlatform(pm *platform.PlatformMetadata, apps app.Apps, ignoreUnimplemented bool) (err error) {
	for _, app := range pkgutils.ToRefArr(apps) {
		deps := pkgutils.ToRefArr(app.Dependencies)
		if app.Compute.ID != "" {
			deps = append(deps, &app.Compute)
		}

		for _, dep := range deps {
			if !dep.NoProvision {
				err = validateDependency(pm, dep)
				if ignoreUnimplemented && errors.Is(err, ErrComponentNotImplemented) {
					err = nil
				} else if err != nil {
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
		return eris.Wrapf(ErrComponentNotImplemented, "'%s.%s'", appDep.ID, appDep.Use)
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
