package dependencies

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/cldcvr/terrarium/src/cli/internal/config"
	"github.com/cldcvr/terrarium/src/pkg/db"
	"github.com/cldcvr/terrarium/src/pkg/metadata/dependency"
	"github.com/cldcvr/terrarium/src/pkg/metadata/taxonomy"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var depIfaceDirectoryFlag string

var dependencyCmd = &cobra.Command{
	Use:   "dependencies",
	Short: "Harvests dependencies from the given directory",
	Long: heredoc.Docf(`
		The 'dependencies' command is used to harvest dependency information from YAML or YML files located
		in a specified directory. It parses these files to extract dependency details and stores them in the database
		for further reference.

		To use this command, provide the path to the directory containing the YAML or YML files using the '--dir' flag.
		The command will recursively process all valid YAML files within the directory, extracting information such as
		taxonomy, title, description, inputs, and outputs. The extracted data is then stored in the database.

		Example usage:
  			terrarium dependencies --dir /path/to/yaml/files

		Please ensure that the provided directory contains valid YAML or YML files with the appropriate structure to avoid any errors.
`),
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

// processYAMLFiles recursively processes YAML files in the specified directory.
func processYAMLFiles(directory string) error {
	return filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if the file is a YAML file (ends with .yaml or .yml)
		if info.IsDir() || (!strings.HasSuffix(info.Name(), ".yaml") && !strings.HasSuffix(info.Name(), ".yml")) {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		// Process the YAML data and insert into the database
		err = processYAMLData(path, data)
		if err != nil {
			return err
		}

		return nil
	})
}

// processYAMLData processes the YAML data and inserts it into the database.
func processYAMLData(path string, data []byte) error {
	g, err := config.DBConnect()
	if err != nil {
		return err
	}

	var yamlData map[string][]dependency.Interface

	err = yaml.Unmarshal(data, &yamlData)
	if err != nil {
		return fmt.Errorf("error parsing YAML file %s: %w", path, err)
	}

	for _, dep := range yamlData["dependency-interfaces"] {
		var taxonomyID uuid.UUID

		if dep.Taxonomy != "" {
			// Split the taxonomy string into levels
			levels := taxonomy.NewTaxonomy(dep.Taxonomy).Split()

			// Retrieve the TaxonomyID based on the levels from the taxonomy table
			taxonomyID, err = getTaxonomyID(g, levels)
			if err != nil {
				return fmt.Errorf("error retrieving TaxonomyID: %w", err)
			}
		}

		// Create a db.Dependency instance
		dbDep := &db.Dependency{
			TaxonomyID:  taxonomyID,
			Title:       dep.Title,
			Description: dep.Description,
			Inputs:      dep.Inputs,
			Outputs:     dep.Outputs,
		}

		// Call CreateDependencyInterface with the updated Dependency object
		_, err = g.CreateDependencyInterface(dbDep)
		if err != nil {
			return fmt.Errorf("error updating the database: %w", err)
		}
	}
	return nil
}

// getTaxonomyID retrieves the TaxonomyID based on the taxonomy levels from the taxonomy table.
func getTaxonomyID(g db.DB, levels []string) (uuid.UUID, error) {
	uniqueFields := make(map[string]interface{})
	for i, level := range levels {
		uniqueFields[fmt.Sprintf("level%d", i+1)] = level
	}

	taxonomy, err := g.GetTaxonomyByFieldName("level1", levels[0])
	if err != nil {
		return uuid.UUID{}, err
	}

	for i := 1; i < len(levels); i++ {
		taxonomy, err = g.GetTaxonomyByFieldName(fmt.Sprintf("level%d", i+1), levels[i])
		if err != nil {
			return uuid.UUID{}, err
		}
	}

	return taxonomy.ID, nil
}
