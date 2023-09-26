// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package dependencies

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cldcvr/terrarium/src/pkg/db"
	"github.com/cldcvr/terrarium/src/pkg/jsonschema"
	"github.com/cldcvr/terrarium/src/pkg/metadata/dependency"
	"github.com/cldcvr/terrarium/src/pkg/metadata/taxonomy"
	"github.com/google/uuid"
	"github.com/rotisserie/eris"
	"gopkg.in/yaml.v3"
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
		return eris.Wrapf(err, "error parsing YAML file %s", path)
	}

	for _, dep := range yamlData["dependency-interfaces"] {
		err = processDependency(g, dep)
		if err != nil {
			return eris.Wrapf(err, "error while updating db")
		}
	}
	return nil
}

func processDependency(g db.DB, dep dependency.Interface) error {
	var taxonomyID uuid.UUID
	if dep.Taxonomy != "" {
		var dbTax db.Taxonomy
		// Split the taxonomy string into levels
		levels := taxonomy.NewTaxonomy(dep.Taxonomy).Split()
		// Please refer to TER-209 for more details to update the following snippet of code to match the
		// taxonomy levels in the dependency interface yaml to the database
		for i, level := range levels {
			tax, err := g.GetTaxonomyByFieldName(fmt.Sprintf("level%d", i+1), level)
			if err != nil {
				return eris.Wrap(err, "error getting taxonomy data")
			}
			dbTax = tax
		}
		taxonomyID = dbTax.ID
	}

	// Create a db.Dependency instance
	dbDep := &db.Dependency{
		TaxonomyID:  taxonomyID,
		InterfaceID: dep.ID,
		Title:       dep.Title,
		Description: dep.Description,
	}

	_, err := g.CreateDependencyInterface(dbDep)
	if err != nil {
		return eris.Wrap(err, "error updating the database")
	}

	attrs := []struct {
		Node     *jsonschema.Node
		Computed bool
	}{
		{dep.Inputs, false}, // For inputs, Computed is false
		{dep.Outputs, true}, // For outputs, Computed is true
	}

	for _, attr := range attrs {
		dbAttr := &db.DependencyAttribute{
			DependencyID: dbDep.ID,
			Name:         dep.Title,
			Schema:       attr.Node,
			Computed:     attr.Computed,
		}

		_, err = g.CreateDependencyAttribute(dbAttr)
		if err != nil {
			return eris.Wrap(err, "error creating dependency attribute")
		}
	}
	return nil
}
