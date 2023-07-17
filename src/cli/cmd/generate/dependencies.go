package generate

import (
	"os"
	"path/filepath"

	"github.com/rotisserie/eris"
	"gopkg.in/yaml.v3"
)

type AppDependency struct {
	ID      string
	App     string
	Type    string
	Inputs  map[string]interface{}
	Outputs map[string]interface{}
}

type AppDependencies []AppDependency

type TerrariumFile struct {
	App TerrariumFile_App
}

type TerrariumFile_App struct {
	Name         string
	Service      AppDependency
	Dependencies map[string]AppDependency
}

// Parse is responsible for parsing the terrarium dependencies at the provided path.
func (appDeps *AppDependencies) Parse(appPath string) error {
	// If appPath is a directory, update appPath to point to 'terrarium.yaml' within the directory
	fi, err := os.Stat(appPath)
	if err != nil {
		return eris.Wrap(err, "failed to get app path info")
	}

	if fi.IsDir() {
		appPath = filepath.Join(appPath, "terrarium.yaml")
	}

	// Read the yaml file
	data, err := os.ReadFile(appPath)
	if err != nil {
		return eris.Wrap(err, "failed to read the app terrarium metadata file")
	}

	// Unmarshal the data into TerrariumFile
	var terrariumFile TerrariumFile
	if err := yaml.Unmarshal(data, &terrariumFile); err != nil {
		return eris.Wrap(err, "failed to unmarshal yaml data")
	}

	// Set AppDependency for Service
	terrariumFile.App.Service.ID = terrariumFile.App.Name
	terrariumFile.App.Service.App = terrariumFile.App.Name
	*appDeps = append(*appDeps, terrariumFile.App.Service)

	// Set AppDependency for Dependencies
	for id, dep := range terrariumFile.App.Dependencies {
		dep.ID = id
		dep.App = terrariumFile.App.Name
		*appDeps = append(*appDeps, dep)
	}

	return nil
}

// GetById retrieves the AppDependency associated with a specific id
func (appDeps AppDependencies) GetById(id string) *AppDependency {
	for _, appDep := range appDeps {
		if appDep.ID == id {
			return &appDep
		}
	}
	return nil
}

// FilterByApp retrieves the AppDependencies associated with a specific application interface
func (appDeps AppDependencies) FilterByApp(appInterface string) AppDependencies {
	var filteredDeps AppDependencies
	for _, appDep := range appDeps {
		if appDep.App == appInterface {
			filteredDeps = append(filteredDeps, appDep)
		}
	}
	return filteredDeps
}

// FilterByType retrieves the AppDependencies associated with a specific type
func (appDeps AppDependencies) FilterByType(appTypes string) AppDependencies {
	var filteredDeps AppDependencies
	for _, appDep := range appDeps {
		if appDep.Type == appTypes {
			filteredDeps = append(filteredDeps, appDep)
		}
	}
	return filteredDeps
}

// GetUniqueTypes retrieves all unique types from the AppDependencies
func (appDeps AppDependencies) GetUniqueTypes() []string {
	typeMap := make(map[string]bool)
	for _, appDep := range appDeps {
		typeMap[appDep.Type] = true
	}

	// Convert the keys of the typeMap to a slice
	var uniqueTypes []string
	for key := range typeMap {
		uniqueTypes = append(uniqueTypes, key)
	}

	return uniqueTypes
}
