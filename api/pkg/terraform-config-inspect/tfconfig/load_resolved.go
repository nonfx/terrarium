package tfconfig

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

type ModuleReference struct {
	Key     string
	Source  string
	Version string
	Dir     string
}
type ResolvedModulesSchema struct {
	Modules []ModuleReference
	baseDir string
}

func (s *ResolvedModulesSchema) UnmarshalFromFile(schemaFilePath string) error {
	modulesFile, err := os.ReadFile(schemaFilePath)
	if err != nil {
		return err
	}
	s.baseDir = filepath.Clean(filepath.Join(filepath.Dir(filepath.Clean(schemaFilePath)), "../../")) // modules.json is in <root>/.terraform/modules

	if err := json.Unmarshal(modulesFile, s); err != nil {
		return err
	}

	return nil
}

func (s ResolvedModulesSchema) GetTopLevelModulePaths() []string {
	topLevelModulePaths := []string{}
	for _, item := range s.Modules {
		if !strings.Contains(item.Key, ".") && item.Key != "" {
			topLevelModulePaths = append(topLevelModulePaths, s.buildAbsPath(item.Dir))
		}
	}
	return topLevelModulePaths
}

func (s ResolvedModulesSchema) buildAbsPath(relPath string) string {
	return filepath.Clean(filepath.Join(s.baseDir, relPath))
}

func (s ResolvedModulesSchema) Get(source string, version string) string {
	for _, item := range s.Modules {
		if item.Source == source && item.Version == version {
			return s.buildAbsPath(item.Dir)
		}
	}
	return ""
}

func (s ResolvedModulesSchema) Find(path string) *ModuleReference {
	relDir, err := filepath.Rel(s.baseDir, path)
	if err == nil {
		for _, item := range s.Modules {
			if item.Dir == relDir {
				return &item
			}
		}
	}
	return nil
}

func LoadModulesFromResolvedSchema(schemaFilePath string) ([]*Module, []Diagnostics, error) {
	resolvedModules := ResolvedModulesSchema{}
	if err := resolvedModules.UnmarshalFromFile(schemaFilePath); err != nil {
		return nil, nil, err
	}

	topLevelModulePaths := resolvedModules.GetTopLevelModulePaths()
	modules := make([]*Module, 0, len(topLevelModulePaths))
	diags := make([]Diagnostics, 0, len(topLevelModulePaths))
	for _, moduleFilePath := range topLevelModulePaths {
		mod, diag := LoadModule(moduleFilePath, &resolvedModules)
		modules = append(modules, mod)
		diags = append(diags, diag)
	}

	return modules, diags, nil
}
