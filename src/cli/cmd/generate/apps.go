package generate

import (
	"os"
	"path"

	"github.com/cldcvr/terrarium/src/pkg/metadata/app"
	"github.com/cldcvr/terrarium/src/pkg/metadata/platform"
	"github.com/rotisserie/eris"
	"golang.org/x/exp/slices"
)

const (
	defaultYAMLFileName = "terrarium.yaml"
)

func fetchApps(appPaths []string) (app.Apps, error) {
	apps := make(app.Apps, len(appPaths))
	for i, appPath := range appPaths {
		content, err := readAppDependency(appPath)
		if err != nil {
			return nil, err
		}

		appObj, err := app.NewApp(content)
		if err != nil {
			return nil, err
		}

		apps[i] = *appObj
	}

	apps.SetDefaults()
	err := apps.Validate()
	if err != nil {
		return nil, err
	}

	return apps, nil
}

// readAppDependency if the given path is a directory, then reads terrarium.yaml file in the directory.
// else if the path points to a YAML file, then returns content of the yaml file.
// otherwise error out
func readAppDependency(appYamlPath string) ([]byte, error) {
	// Check if the path exists and get its FileInfo
	info, err := os.Stat(appYamlPath)
	if err != nil {
		return nil, eris.Wrapf(err, "invalid file path: %s", appYamlPath)
	}

	// If it's a directory, append terrarium.yaml to the path
	if info.IsDir() {
		appYamlPath = path.Join(appYamlPath, defaultYAMLFileName)
	} else if !slices.Contains([]string{"yml", "yaml"}, path.Ext(appYamlPath)) {
		// If it's a file but not a .yaml, return an error
		return nil, eris.New("provided path is not a directory or a .yaml|.yml file")
	}

	// Read the file content
	content, err := os.ReadFile(appYamlPath)
	if err != nil {
		return nil, eris.Wrapf(err, "error reading file at: %s", appYamlPath)
	}

	return content, nil
}

func matchAppAndPlatform(pm *platform.PlatformMetadata, apps app.Apps) (err error) {
	for i, app := range apps {
		if app.Compute.ID != "" {
			err = validateDependency(pm, &(apps[i].Compute))
			if err != nil {
				return
			}
		}

		for j, dep := range app.Dependencies {
			if !dep.NoProvision {
				err = validateDependency(pm, &(apps[i].Dependencies[j]))
				if err != nil {
					return
				}
			}
		}
	}

	return nil
}

func validateDependency(pm *platform.PlatformMetadata, appDep *app.Dependency) error {
	comp := pm.Components.GetByID(appDep.Use)
	if comp == nil {
		return eris.Errorf("component '%s.%s' is not implemented in the platform", appDep.ID, appDep.Use)
	}

	err := comp.Inputs.Validate(appDep.Inputs)
	if err != nil {
		return eris.Wrapf(err, "component '%s.%s' does not contain a valid set of inputs", appDep.ID, appDep.Use)
	}

	comp.Inputs.ApplyDefaultsToMSI(appDep.Inputs)
	return nil
}
