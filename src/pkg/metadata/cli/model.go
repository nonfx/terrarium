// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package cli

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/rotisserie/eris"
	"gopkg.in/yaml.v3"
)

type FarmModuleRef struct {
	Name    string `yaml:"name"`
	Source  string `yaml:"source"`
	Version string `yaml:"version,omitempty"`
	Export  bool   `yaml:"export,omitempty"`
}

// CreateTerraformFile writes a terraform module file.
// The file wil be stored in a sub-directory under a provided root directoy.
// If the root path is empty the file will be stored in a temporary directory instead.
// When using a temporary directory each module will need to be freshly initialized by 'terraform init' every time.
func (r FarmModuleRef) CreateTerraformFile(root string) (dirPath string, filePath string, err error) {
	if root != "" {
		// use a given root directory
		if fp, err := os.Stat(root); err != nil {
			return "", "", eris.Wrapf(err, "provided working directory '%s' could not be accessed", root)
		} else if !fp.IsDir() {
			return "", "", eris.Errorf("provided working directory path '%s' is not a directory", root)
		}

		dirPath = path.Join(root, r.Name)
		if fp, err := os.Stat(dirPath); os.IsNotExist(err) {
			// create a new module output dir inside the root if it does not exist
			if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
				return "", "", eris.Wrapf(err, "could not create module output directory '%s'", dirPath)
			}
		} else if err != nil {
			return "", "", eris.Wrapf(err, "module output directory '%s' could not be accessed", dirPath)
		} else if !fp.IsDir() {
			return "", "", eris.Errorf("module output directory path '%s' is not a directory", dirPath)
		}
	} else {
		// use TEMP dir instead
		dirPath, err = os.MkdirTemp("", fmt.Sprintf("tr_%s_*", r.Name))
		if err != nil {
			return "", "", eris.Wrap(err, "could not create output directory")
		}
	}

	// always overwrite the main.tf file to make sure the executed TF code is consistent with the module list entry
	fp, err := os.Create(path.Join(dirPath, "main.tf"))
	if err != nil {
		return "", "", eris.Wrap(err, "could not open output file")
	}
	defer fp.Close()
	if _, err := fp.WriteString(r.ToTerraform()); err != nil {
		return "", "", eris.Wrapf(err, "could not write to output file '%s'", fp.Name())
	}

	filePath = fp.Name()
	return
}

func (r FarmModuleRef) ToTerraform() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("module \"%s\" {\n", r.Name))
	b.WriteString(fmt.Sprintf("	source = \"%s\"\n", r.Source))
	if r.Version != "" {
		b.WriteString(fmt.Sprintf("	version = \"%s\"\n", r.Version))
	}
	b.WriteString("}\n")
	return b.String()
}

type FarmModuleList struct {
	Farm []FarmModuleRef `yaml:"farm"`
}

func LoadFarmModules(listFilePath string) (FarmModuleList, error) {
	moduleList, err := loadFarmModules(listFilePath)
	if err != nil {
		return moduleList, eris.Wrapf(err, "failed to load farm module list file '%s'", listFilePath)
	} else if len(moduleList.Farm) < 1 {
		return moduleList, eris.Errorf("farm module list file '%s' is empty", listFilePath)
	}
	if err := moduleList.Validate(); err != nil {
		return moduleList, eris.Wrapf(err, "farm module list file '%s' has invalid items", listFilePath)
	}
	return moduleList, nil
}

func loadFarmModules(listFilePath string) (moduleList FarmModuleList, err error) {
	listFile, err := os.ReadFile(listFilePath)
	if err != nil {
		return
	}
	if err = yaml.Unmarshal(listFile, &moduleList); err != nil {
		return
	}
	return
}

func (l FarmModuleList) Validate() error {
	itemCount := len(l.Farm)
	uniqueModuleReferences := make(map[string]*FarmModuleRef, itemCount)
	uniqueExportNames := make(map[string]int, itemCount)
	for i, item := range l.Farm {
		if item.Name == "" {
			return eris.Errorf("module at index %d must have a unique name", i)
		}
		if item.Source == "" {
			return eris.Errorf("module '%s' must have a source", item.Name)
		}
		if _, exists := uniqueExportNames[item.Name]; exists {
			return eris.Errorf("module '%s' at index %d has a duplicate name", item.Name, i)
		}
		uniqueExportNames[item.Name] = i
		refKey := fmt.Sprintf("%s@%s", item.Source, item.Version)
		if found, exists := uniqueModuleReferences[refKey]; exists {
			return eris.Errorf("module '%s' is duplicate of module '%s'", item.Name, found.Name)
		}
		uniqueModuleReferences[refKey] = &item
	}
	return nil
}
