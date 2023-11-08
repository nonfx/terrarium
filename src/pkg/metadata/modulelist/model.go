// Copyright (c) Ollion
// SPDX-License-Identifier: Apache-2.0

package modulelist

import (
	"fmt"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/cldcvr/terrarium/src/pkg/utils"
	"github.com/rotisserie/eris"
	"gopkg.in/yaml.v3"
)

type FarmModuleList struct {
	Farm []FarmModuleRef `yaml:"farm"`
}

type FarmModuleRef struct {
	Name    string `yaml:"name"`
	Source  string `yaml:"source"`
	Version string `yaml:"version,omitempty"`
	Export  bool   `yaml:"export,omitempty"`
	Group   string `yaml:"group,omitempty"`
}

type FarmModuleGroupsMap map[string]FarmModuleGroup

type FarmModuleGroup struct {
	Name    string
	Export  bool
	Modules []FarmModule
}

type FarmModule struct {
	Name    string `yaml:"name"`
	Source  string `yaml:"source"`
	Version string `yaml:"version,omitempty"`
}

var moduleTmpl = template.Must(template.New("tfmodules").Parse(heredoc.Doc(`
	{{range . -}}
	module "{{.Name}}" {
		source = "{{.Source}}"
		{{if .Version}}version = "{{.Version}}"{{end}}
	}

	{{end}}
`)))

// CreateTerraformFile writes a terraform module file.
// The file wil be stored in a sub-directory under a provided root directory.
// If the root path is empty the file will be stored in a temporary directory instead.
// When using a temporary directory each module-group will need to be freshly initialized by 'terraform init' every time.
func (g FarmModuleGroup) CreateTerraformFile(rootDir string) (dirPath string, err error) {
	dirPath, err = g.PrepareDir(rootDir)
	if err != nil {
		return
	}

	// always overwrite the main.tf file to make sure the executed TF code is consistent with the module list entry
	fp, err := os.Create(path.Join(dirPath, "main.tf"))
	if err != nil {
		return "", eris.Wrap(err, "could not open output file")
	}
	defer fp.Close()

	str, err := g.ToTerraform()
	if err != nil {
		return
	}

	if _, err := fp.WriteString(str); err != nil {
		return "", eris.Wrapf(err, "could not write to output file '%s'", fp.Name())
	}

	return
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

// Validate ensures that the FarmModuleList adheres to the following criteria:
// - Each module's name and source must be non-empty.
// - Each module's name must be unique.
// - Each combination of a module's source and version must be unique.
// - A group name must not conflict with a module name.
// - All modules within the same group must have the same 'export' flag value.
func (l FarmModuleList) Validate() error {
	itemCount := len(l.Farm)
	uniqueModuleReferences := make(map[string]*FarmModuleRef, itemCount)
	uniqueExportNames := make(map[string]int, itemCount)
	uniqueGroupVal := map[string]bool{}

	for i, item := range l.Farm {
		if err := item.validateBasicFields(uniqueExportNames, uniqueGroupVal); err != nil {
			return err
		}
		uniqueExportNames[item.Name] = i

		refKey := fmt.Sprintf("%s@%s", item.Source, item.Version)
		if err := item.validateUniqueModuleReferences(uniqueModuleReferences, refKey); err != nil {
			return err
		}
		uniqueModuleReferences[refKey] = &item

		if err := item.validateGroup(uniqueExportNames, uniqueGroupVal); err != nil {
			return err
		}
		uniqueGroupVal[item.Group] = item.Export
	}
	return nil
}

// validateBasicFields checks the basic fields of a FarmModuleRef:
// - The name must be non-empty.
// - The source must be non-empty.
// - The name must be unique.
func (item FarmModuleRef) validateBasicFields(uniqueExportNames map[string]int, uniqueGroupVal map[string]bool) error {
	if item.Name == "" {
		return eris.Errorf("module must have a non-empty name")
	}
	if item.Source == "" {
		return eris.Errorf("module '%s' must have a source", item.Name)
	}
	if _, exists := uniqueExportNames[item.Name]; exists {
		return eris.Errorf("module '%s' has a duplicate name", item.Name)
	}
	if _, exists := uniqueGroupVal[item.Name]; exists {
		return eris.Errorf("module '%s' has name conflict with a group name", item.Name)
	}
	return nil
}

// validateUniqueModuleReferences ensures that each combination of source and version is unique.
func (item FarmModuleRef) validateUniqueModuleReferences(uniqueModuleReferences map[string]*FarmModuleRef, refKey string) error {
	if found, exists := uniqueModuleReferences[refKey]; exists {
		return eris.Errorf("module '%s' is duplicate of module '%s'", item.Name, found.Name)
	}
	return nil
}

// validateGroup checks the following:
// - The group name must not conflict with a module name.
// - All modules within the same group must have the same 'export' flag value.
func (item FarmModuleRef) validateGroup(uniqueExportNames map[string]int, uniqueGroupVal map[string]bool) error {
	if item.Group == "" {
		return nil
	}
	if _, exists := uniqueExportNames[item.Group]; exists {
		return eris.Errorf("group '%s' has name conflict with a module name", item.Group)
	}
	if wantExport, exists := uniqueGroupVal[item.Group]; exists && wantExport != item.Export {
		return eris.Errorf("group '%s' has conflicting export flag! all modules in a group must have same export value", item.Group)
	}
	return nil
}

func (g FarmModuleGroup) PrepareDir(rootDir string) (dirPath string, err error) {
	if rootDir != "" {
		dirPath, err = utils.SetupDir(path.Join(rootDir, g.Name))
		if err != nil {
			return "", eris.Wrapf(err, "failed to use the root directory: %s, for module: %s", rootDir, g.Name)
		}
	} else {
		dirPath, err = os.MkdirTemp("", fmt.Sprintf("tr_%s_*", g.Name))
		if err != nil {
			return "", eris.Wrapf(err, "could not allocate temporary directory for module: %s", g.Name)
		}
	}

	return
}

// GetGroupName returns the group name for a FarmModuleRef.
// If the Group field is not empty, it returns that.
// Otherwise, it returns the Name field.
func (r FarmModuleRef) GetGroupName() string {
	if r.Group != "" {
		return r.Group
	}

	return r.Name
}

// Groups organizes the FarmModuleRefs in the FarmModuleList into a FarmModuleGroupsMap.
// Each group in the map is identified by its name and contains an array of FarmModules.
// The Export flag for each group is also set based on the FarmModuleRef.
func (l FarmModuleList) Groups() FarmModuleGroupsMap {
	groups := FarmModuleGroupsMap{}
	for _, r := range l.Farm {
		groupName := r.GetGroupName()
		g := groups[groupName]
		g.Name, g.Export = groupName, r.Export
		g.Modules = append(g.Modules, FarmModule{
			Name:    r.Name,
			Source:  r.Source,
			Version: r.Version,
		})
		groups[groupName] = g
	}

	return groups
}

func (refArr FarmModuleGroup) ToTerraform() (string, error) {
	str := &strings.Builder{}
	err := moduleTmpl.Execute(str, refArr.Modules)
	if err != nil {
		return "", eris.Wrapf(err, "failed to execute tf-module template on group: %s", refArr.Name)
	}

	return str.String(), nil
}

// FilterExport filters the FarmModuleGroupsMap based on the Export flag.
// It returns a new FarmModuleGroupsMap containing only the groups that match the given Export flag.
func (groups FarmModuleGroupsMap) FilterExport(wantExport bool) (filteredGroups FarmModuleGroupsMap) {
	filteredGroups = FarmModuleGroupsMap{}
	for k, g := range groups {
		if g.Export == wantExport {
			filteredGroups[k] = g
		}
	}

	return
}
