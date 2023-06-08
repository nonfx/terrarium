package tfconfig

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	tfversion "github.com/hashicorp/go-version"
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

func FilterModulesKeepTopLevel(s ResolvedModulesSchema) ResolvedModulesSchema {
	filteredModules := make([]ModuleReference, 0, len(s.Modules))
	for _, item := range s.Modules {
		if !strings.Contains(item.Key, ".") && item.Key != "" {
			filteredModules = append(filteredModules, item)
		}
	}

	return ResolvedModulesSchema{
		Modules: filteredModules,
		baseDir: s.baseDir,
	}
}

func FilterModulesOmitLocal(s ResolvedModulesSchema) ResolvedModulesSchema {
	filteredModules := make([]ModuleReference, 0, len(s.Modules))
	for _, item := range s.Modules {
		if !strings.HasPrefix(item.Source, ".") && item.Key != "" {
			filteredModules = append(filteredModules, item)
		}
	}

	return ResolvedModulesSchema{
		Modules: filteredModules,
		baseDir: s.baseDir,
	}
}

func FilterModulesOmitHidden(s ResolvedModulesSchema) ResolvedModulesSchema {
	filteredModules := make([]ModuleReference, 0, len(s.Modules))
	for _, item := range s.Modules {
		if !strings.HasPrefix(item.Key, "tr-hide-") {
			filteredModules = append(filteredModules, item)
		}
	}

	return ResolvedModulesSchema{
		Modules: filteredModules,
		baseDir: s.baseDir,
	}
}

func (s ResolvedModulesSchema) GetPaths() (paths []string) {
	paths = make([]string, len(s.Modules))

	for i, m := range s.Modules {
		paths[i] = s.buildAbsPath(m.Dir)
	}

	return
}

func (s ResolvedModulesSchema) buildAbsPath(relPath string) string {
	return filepath.Clean(filepath.Join(s.baseDir, relPath))
}

func (s ResolvedModulesSchema) Get(source string, version string) (path string) {
	skipVersionCheck := version == ""
	if requiredVersion, err := tfversion.NewConstraint(version); skipVersionCheck || err == nil {
		for _, item := range s.Modules {
			if moduleSource := strings.TrimPrefix(item.Source, "registry.terraform.io/"); moduleSource == source {
				if moduleVersion, err := tfversion.NewVersion(item.Version); skipVersionCheck || err == nil && requiredVersion.Check(moduleVersion) {
					return s.buildAbsPath(item.Dir)
				}
			}
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

func (m ModuleReference) GetNormalizedKey() string {
	dotSplits := strings.Split(m.Key, ".")
	lastIndex := len(dotSplits) - 1
	if lastIndex < 0 {
		return m.Key
	}

	return dotSplits[lastIndex]
}

type ResolvedModuleSchemaFilter func(ResolvedModulesSchema) ResolvedModulesSchema

func LoadModulesFromResolvedSchema(schemaFilePath string, filters ...ResolvedModuleSchemaFilter) ([]*Module, []Diagnostics, error) {
	resolvedModules := ResolvedModulesSchema{}
	if err := resolvedModules.UnmarshalFromFile(schemaFilePath); err != nil {
		return nil, nil, err
	}

	filteredModules := resolvedModules
	for _, filter := range filters {
		filteredModules = filter(filteredModules)
	}

	filteredModulePaths := filteredModules.GetPaths()

	modules := make([]*Module, 0, len(filteredModulePaths))
	diags := make([]Diagnostics, 0, len(filteredModulePaths))

	for _, moduleFilePath := range filteredModulePaths {
		mod, diag := LoadModule(moduleFilePath, &resolvedModules)
		modules = append(modules, mod)
		diags = append(diags, diag)
	}

	return modules, diags, nil
}
