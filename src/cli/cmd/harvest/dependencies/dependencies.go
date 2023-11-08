// Copyright (c) Ollion
// SPDX-License-Identifier: Apache-2.0

package dependencies

import (
	"os"
	"path/filepath"

	"github.com/cldcvr/terrarium/src/pkg/db"
	"github.com/cldcvr/terrarium/src/pkg/jsonschema"
	"github.com/cldcvr/terrarium/src/pkg/metadata/dependency"
	"github.com/cldcvr/terrarium/src/pkg/utils"
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
		if info.IsDir() || !utils.IsYaml(info.Name()) {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return eris.Wrapf(err, "error reading file '%s'", path)
		}

		// Process the YAML data and insert into the database
		err = processYAMLData(g, data)
		if err != nil {
			return eris.Wrapf(err, "error processing file '%s'", path)
		}

		return nil
	})
}

// processYAMLData processes the YAML data and inserts it into the database.
func processYAMLData(g db.DB, data []byte) error {
	var yamlData dependency.File

	err := yaml.Unmarshal(data, &yamlData)
	if err != nil {
		return eris.Wrapf(err, "error parsing YAML content")
	}

	for _, dep := range yamlData.DependencyInterfaces {
		err = processDependency(g, dep)
		if err != nil {
			return eris.Wrapf(err, "error while processing interface '%s'", dep.ID)
		}
	}
	return nil
}

func processDependency(g db.DB, dep dependency.Interface) error {
	taxonomyID, err := getTaxonomy(g, dep)
	if err != nil {
		return err
	}

	// Create a db.Dependency instance
	dbDep := &db.Dependency{
		InterfaceID: dep.ID,
		Title:       dep.Title,
		Description: dep.Description,
	}

	if taxonomyID != uuid.Nil {
		dbDep.TaxonomyID = &taxonomyID
	}

	_, err = g.CreateDependencyInterface(dbDep)
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
		if attr.Node == nil || attr.Node.Properties == nil {
			continue
		}

		for k, n := range attr.Node.Properties {
			dbAttr := &db.DependencyAttribute{
				DependencyID: dbDep.ID,
				Name:         k,
				Schema:       n,
				Computed:     attr.Computed,
			}

			_, err = g.CreateDependencyAttribute(dbAttr)
			if err != nil {
				return eris.Wrap(err, "error creating dependency attribute")
			}
		}
	}
	return nil
}

func getTaxonomy(g db.DB, dep dependency.Interface) (id uuid.UUID, err error) {
	levels := dep.Taxonomy.Split()
	if len(levels) == 0 {
		// no error, empty id.
		return
	}

	dbTax := db.TaxonomyFromLevels(levels...)
	id, err = g.CreateTaxonomy(dbTax)
	if err != nil {
		err = eris.Wrapf(err, "failed to get or create taxon '%s'", dep.Taxonomy)
		return
	}

	return
}
