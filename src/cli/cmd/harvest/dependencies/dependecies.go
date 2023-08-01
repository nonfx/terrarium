package dependencies

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cldcvr/terrarium/src/cli/internal/config"
	"github.com/cldcvr/terrarium/src/pkg/db"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var depIfaceDirectoryFlag string

var dependencyCmd = &cobra.Command{
	Use:   "dependencies",
	Short: "Harvests dependencies from the given directory",
	Long:  "Harvests dependencies from all the yaml/yml files in the directory provided and adds it to the database.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return main()
	},
}

func GetCmd() *cobra.Command {
	addFlags()
	return dependencyCmd
}

func addFlags() {
	dependencyCmd.Flags().StringVarP(&depIfaceDirectoryFlag, "dir", "d", "", "path to dependency directory")
}

func main() error {
	err := processYAMLFiles(depIfaceDirectoryFlag)
	if err != nil {
		return err
	}

	return nil
}

func processYAMLFiles(directory string) error {
	return filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if the file is a YAML file (ends with .yaml or .yml)
		if !info.IsDir() && (strings.HasSuffix(info.Name(), ".yaml") || strings.HasSuffix(info.Name(), ".yml")) {
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			// Process the YAML data and insert into the database
			err = processYAMLData(path, data)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func processYAMLData(path string, data []byte) error {
	g, err := config.DBConnect()
	if err != nil {
		return err
	}
	var yamlData map[string][]db.Dependency

	err = yaml.Unmarshal(data, &yamlData)
	if err != nil {
		return fmt.Errorf("error parsing YAML file %s: %w", path, err)
	}

	for _, dep := range yamlData["dependency-interface"] {
		dep.ID = uuid.New()

		// Convert inputs and outputs to JSON before insertion
		inputsJSON, err := toJSONString(dep.Inputs)
		if err != nil {
			return fmt.Errorf("error converting inputs to JSON: %w", err)
		}
		dep.InputsJSON = inputsJSON

		outputsJSON, err := toJSONString(dep.Outputs)
		if err != nil {
			return fmt.Errorf("error converting outputs to JSON: %w", err)
		}
		dep.OutputsJSON = outputsJSON

		// Call CreateDependencyInterface with the updated Dependency object
		_, err = g.CreateDependencyInterface(&dep)
		if err != nil {
			return fmt.Errorf("error updating the database: %w", err)
		}
		fmt.Fprintf(os.Stdout, "Data inserted successfully!\n")
	}

	return nil
}

func toJSONString(data interface{}) (string, error) {
	switch v := data.(type) {
	case map[string]interface{}:
		jsonData := make(map[string]interface{})
		for key, val := range v {
			valJSON, err := toJSONString(val)
			if err != nil {
				return "", err
			}
			jsonData[key] = valJSON
		}
		jsonBytes, err := json.Marshal(jsonData)
		if err != nil {
			return "", err
		}
		return string(jsonBytes), nil
	case []interface{}:
		jsonData := make([]interface{}, len(v))
		for i, val := range v {
			valJSON, err := toJSONString(val)
			if err != nil {
				return "", err
			}
			jsonData[i] = valJSON
		}
		jsonBytes, err := json.Marshal(jsonData)
		if err != nil {
			return "", err
		}
		return string(jsonBytes), nil
	default:
		return fmt.Sprintf("%v", v), nil
	}
}
