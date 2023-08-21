package dependencies

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cldcvr/terrarium/src/pkg/db"
	"github.com/cldcvr/terrarium/src/pkg/metadata/dependency"
	"github.com/cldcvr/terrarium/src/pkg/metadata/taxonomy"
	"github.com/google/uuid"
	"gopkg.in/yaml.v2"
)

// processYAMLFiles recursively processes YAML files in the specified directory.
func processYAMLFiles(g db.DB, directory string) error {
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
		err = processYAMLData(g, path, data)
		if err != nil {
			return err
		}

		return nil
	})
}

// processYAMLData processes the YAML data and inserts it into the database.
func processYAMLData(g db.DB, path string, data []byte) error {
	var yamlData map[string][]dependency.Interface

	err := yaml.Unmarshal(data, &yamlData)
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
			InterfaceID: dep.ID,
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

	var taxonomy db.Taxonomy
	for i, level := range levels {
		tax, err := g.GetTaxonomyByFieldName(fmt.Sprintf("level%d", i+1), level)
		if err != nil {
			return uuid.UUID{}, err
		}
		taxonomy = tax
	}

	return taxonomy.ID, nil
}
