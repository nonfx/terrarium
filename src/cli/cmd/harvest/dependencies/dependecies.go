package dependencies

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cldcvr/terrarium/src/cli/internal/config"
	"github.com/cldcvr/terrarium/src/pkg/db"
	"github.com/cldcvr/terrarium/src/pkg/jsonschema"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"gorm.io/gorm"
)

type Taxonomy string

type Dependency struct {
	Taxonomy    Taxonomy
	Title       string
	Description string
	Inputs      *jsonschema.Node
	Outputs     *jsonschema.Node
}

func (t Taxonomy) Parse() (taxons []string) {
	return strings.Split(string(t), "/")
}

func NewTaxonomy(taxons ...string) Taxonomy {
	return Taxonomy(strings.Join(taxons, "/"))
}

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
	var yamlData map[string][]Dependency

	err = yaml.Unmarshal(data, &yamlData)
	if err != nil {
		return fmt.Errorf("error parsing YAML file %s: %w", path, err)
	}

	for _, dep := range yamlData["dependency-interfaces"] {
		depID := uuid.New()

		// Convert inputs and outputs to JSON before insertion
		inputsJSON, err := toJSONString(dep.Inputs)
		if err != nil {
			return fmt.Errorf("error converting inputs to JSON: %w", err)
		}
		outputsJSON, err := toJSONString(dep.Outputs)
		if err != nil {
			return fmt.Errorf("error converting outputs to JSON: %w", err)
		}

		// Match and retrieve the TaxonomyID from the Taxonomy table
		taxonomyID, err := getTaxonomyID(string(dep.Taxonomy))
		if err != nil {
			return fmt.Errorf("error retrieving TaxonomyID: %w", err)
		}

		// Create a db.Dependency instance
		dbDep := &db.Dependency{
			Model:       db.Model{ID: depID},
			TaxonomyID:  taxonomyID,
			Title:       dep.Title,
			Description: dep.Description,
		}

		// Convert inputsJSON and outputsJSON to *jsonschema.Node
		var inputsNode jsonschema.Node
		err = json.Unmarshal([]byte(inputsJSON), &inputsNode)
		if err != nil {
			return fmt.Errorf("error converting InputsJSON to Node: %w", err)
		}
		dbDep.Inputs = &inputsNode

		var outputsNode jsonschema.Node
		err = json.Unmarshal([]byte(outputsJSON), &outputsNode)
		if err != nil {
			return fmt.Errorf("error converting OutputsJSON to Node: %w", err)
		}
		dbDep.Outputs = &outputsNode

		// Call CreateDependencyInterface with the updated Dependency object
		_, err = g.CreateDependencyInterface(dbDep)
		if err != nil {
			return fmt.Errorf("error updating the database: %w", err)
		}
		fmt.Fprintf(os.Stdout, "Data inserted successfully!\n")
	}

	return nil
}

func getTaxonomyID(t string) (uuid.UUID, error) {
	// Split the taxonomy string
	// levels := dependency.Taxonomy(t).Parse()

	// Query the Taxonomy table based on the levels
	// and retrieve the matching Taxonomy ID
	var taxonomy db.Taxonomy
	// result := config.DBInstance().Where(
	// 	"level1 = ? AND level2 = ? AND level3 = ? AND level4 = ? AND level5 = ? AND level6 = ? AND level7 = ?",
	// 	levels[0], levels[1], levels[2], levels[3], levels[4], levels[5], levels[6],
	// ).First(&taxonomy)
	// if result.Error != nil {
	// 	return uuid.UUID{}, result.Error
	// }

	return taxonomy.ID, nil
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

func getTaxonomyIDByLevels(g *gorm.DB, levels []string) (uuid.UUID, error) {
	uniqueFields := make(map[string]interface{})
	for i, level := range levels {
		uniqueFields[fmt.Sprintf("level%d", i+1)] = level
	}

	var taxonomy db.Taxonomy
	// err := getByUniqueFields(g, uniqueFields, &taxonomy)
	// if err != nil {
	// 	return uuid.Nil, err
	// }

	return taxonomy.ID, nil
}
