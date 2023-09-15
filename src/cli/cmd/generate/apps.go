// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package generate

import (
	"os"
	"path"
	"sort"

	"github.com/cldcvr/terrarium/src/cli/internal/constants"
	"github.com/cldcvr/terrarium/src/pkg/metadata/app"
	"github.com/cldcvr/terrarium/src/pkg/metadata/platform"
	"github.com/cldcvr/terrarium/src/pkg/metadata/utils"
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
	} else if !slices.Contains([]string{".yml", ".yaml"}, path.Ext(appYamlPath)) {
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

func writeAppsEnv(pm *platform.PlatformMetadata, apps app.Apps) error {
	for _, appObj := range apps {
		vars := utils.GetAppEnvTemplate(pm, appObj)
		sort.Sort(vars)
		fileName := "app_" + appObj.ID + ".env.mustache"
		err := os.WriteFile(path.Join(flagOutDir, fileName), []byte(vars.RenderWithQuotes()), constants.ReadWritePermissions)
		if err != nil {
			return eris.Wrapf(err, "failed to write app env file")
		}
	}

	return nil
}
