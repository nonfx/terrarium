// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"fmt"
	"strings"

	"github.com/cldcvr/terrarium/src/pkg/metadata/app"
	"github.com/cldcvr/terrarium/src/pkg/metadata/platform"
	"github.com/hoisie/mustache"
)

const envValTemp = `{{ tr_component_%s_%s.value.%s }}` // 1. component-id, 2. output-name & 3. app-dependency-id

// EnvVars array of env variables
type EnvVars []EnvVar

// EnvVar env variable Object
type EnvVar struct {
	Key, Value string
}

// GetAppEnvTemplate based on the app-dependencies, and component metadata,
// generate a template for env variables such that the variables can be rendered later using
// the terraform state output object.
func GetAppEnvTemplate(pm *platform.PlatformMetadata, app app.App) EnvVars {
	envVars := EnvVars{}
	for _, appDep := range app.GetDependencies() {
		comp := pm.Components.GetByID(appDep.Use)
		if comp == nil || comp.Outputs == nil {
			continue
		}

		depDefaults := getDepDefaults(comp, appDep)

		finalOutputs := getRenderedOutputs(&appDep, depDefaults)

		prefix := getEnvVarPrefix(app.EnvPrefix, appDep.EnvPrefix)

		for k, v := range finalOutputs {
			varName := prefix + k
			envVars = append(envVars, EnvVar{varName, v})
		}
	}

	return envVars
}

// getDepDefaults returns the default environment variables for a given component.
func getDepDefaults(comp *platform.Component, appDep app.Dependency) map[string]string {
	depDefaults := map[string]string{}
	for outputName := range comp.Outputs.Properties {
		depDefaults[outputName] = fmt.Sprintf(envValTemp, comp.ID, outputName, appDep.ID)
	}
	return depDefaults
}

// getRenderedOutputs updates the Outputs of appDep based on depDefaults.
func getRenderedOutputs(appDep *app.Dependency, depDefaults map[string]string) (finalOutputs map[string]string) {
	if len(appDep.Outputs) > 0 {
		finalOutputs = make(map[string]string, len(appDep.Outputs))
		for k, vTemp := range appDep.Outputs {
			finalOutputs[k] = mustache.Render(vTemp, depDefaults)
		}
	} else {
		finalOutputs = make(map[string]string, len(depDefaults))
		for k, v := range depDefaults {
			finalOutputs[strings.ToUpper(k)] = v
		}
	}

	return
}

// getEnvVarPrefix returns the environment variable prefix based on app and appDep.
func getEnvVarPrefix(appPrefix, depPrefix string) string {
	prefix := ""
	if appPrefix != "" {
		prefix += appPrefix + "_"
	}
	if depPrefix != "" {
		prefix += depPrefix + "_"
	}
	return prefix
}

func (e EnvVar) Render(quoteVal bool) string {
	tmpl := "%s=%s"
	if quoteVal {
		tmpl = "%s=%q"
	}
	return fmt.Sprintf(tmpl, e.Key, e.Value)
}

func (vars EnvVars) render(quoteVal bool) string {
	allVars := ""
	for _, v := range vars {
		allVars += v.Render(quoteVal) + "\n"
	}
	return allVars
}

func (vars EnvVars) Render() string {
	return vars.render(false)
}

func (vars EnvVars) RenderWithQuotes() string {
	return vars.render(true)
}

//  Implement sort.Interface to make EnvVars sortable by var name

func (vars EnvVars) Len() int {
	return len(vars)
}

func (vars EnvVars) Less(i, j int) bool {
	return vars[i].Key < vars[j].Key
}

func (vars EnvVars) Swap(i, j int) {
	vars[i], vars[j] = vars[j], vars[i]
}
