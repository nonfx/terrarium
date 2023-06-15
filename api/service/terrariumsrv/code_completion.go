package terrariumsrv

import (
	"bytes"
	"context"
	"embed"
	"html/template"
	"strings"

	"github.com/cldcvr/terrarium/api/db"
	"github.com/cldcvr/terrarium/api/pkg/pb/terrariumpb"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/google/uuid"
	"golang.org/x/exp/slices"
)

const trCommentPrefix = "## TERRARIUM MODULE ##"

//go:embed hcl_go_tmpl
var tmplFileReader embed.FS

var goTfTemplates *template.Template

func init() {
	goTfTemplates = template.Must(template.ParseFS(tmplFileReader, "hcl_go_tmpl/module.go.tpl"))
}

func (s Service) CodeCompletion(ctx context.Context, req *terrariumpb.CompletionRequest) (*terrariumpb.CompletionResponse, error) {
	existingModulesMap := parseModulesFromCode(req.CodeContext)

	allModules, err := s.fetchModulesRecursive(ctx, req.Modules)
	if err != nil {
		return nil, err
	}

	suggestion, err := modulesToHclTemplate(allModules, req.Modules, existingModulesMap)
	if err != nil {
		return nil, err
	}

	return &terrariumpb.CompletionResponse{
		Suggestions: []string{suggestion},
	}, nil
}

// parseModulesFromCode find all the existing modules being used from the code provided as context to completion.
func parseModulesFromCode(code string) (moduleSrcMap map[string]string) {
	lines := strings.Split(code, "\n")
	moduleSrcMap = map[string]string{}
	for _, cl := range lines {
		if !strings.HasPrefix(strings.TrimSpace(cl), trCommentPrefix) {
			continue
		}

		tokens := strings.Split(cl, "##")
		src := strings.TrimSpace(tokens[2])
		moName := strings.TrimSpace(tokens[3])

		moduleSrcMap[src] = moName
	}

	return
}

// fetchModulesRecursive fetches terraform modules recursively by it's dependencies.
// For example: if Module A depends on Module B, and Module B depends on Module C, then this returns modules A, B, & C.
// Along with each of their input attributes to output attribute mapping pre-populated in the `TFModule.Attributes` array.
func (s Service) fetchModulesRecursive(ctx context.Context, modules []string) ([]*db.TFModule, error) {
	modulesToFetch := mapset.NewSet(modules...)
	modulesFetched := mapset.NewSet([]string{}...)
	allModules := []*db.TFModule{} // maintain sequence of modules. [<requested modules>, <its dependencies>, <its dependency's dependencies>, ...]

	for modulesToFetch.Cardinality() > 0 {
		modules, err := s.db.FindOutputMappingsByModuleID(parseUUIDArr(modulesToFetch.ToSlice())...)
		if err != nil {
			return nil, err
		}

		modulesToFetch.Clear()
		for i := range modules {
			m := &modules[i]

			connectedModuleIds := normalizeModuleDependency(m) // normalize the module

			allModules = append(allModules, m)           // add to the list
			modulesFetched.Add(m.ID.String())            // mark the module being fetched
			modulesToFetch.Append(connectedModuleIds...) // add all it's dependencies to fetch list
		}

		modulesToFetch = modulesToFetch.Difference(modulesFetched)
	}

	return allModules, nil
}

// normalizeModuleDependency simplifies output module attributes to one
// if there are multiple modules that resolve an input attribute.
func normalizeModuleDependency(m *db.TFModule) []string {
	connectedModuleIdSet := mapset.NewSet([]string{}...)

	for i := range m.Attributes {
		attr := &m.Attributes[i]
		// Skip iteration if no ResourceAttribute
		if attr.ResourceAttribute == nil {
			continue
		}

		normalizeOutputMappings(attr, connectedModuleIdSet)
	}

	return connectedModuleIdSet.ToSlice()
}

// normalizeOutputMappings makes selection of one mapping of module attribute. From multiple available.
// There are one to many relations on two levels - input resource attribute to output resource attribute mappings,
// and resource attribute to module attribute mappings, which are both normalized to one.
func normalizeOutputMappings(moInAttr *db.TFModuleAttribute, connectedModuleIdSet mapset.Set[string]) {
	var chosenResOutMapping *db.TFResourceAttributesMapping
	var chosenMoOutAttr *db.TFModuleAttribute
	for i := range moInAttr.ResourceAttribute.OutputMappings {
		resOutMapping := &moInAttr.ResourceAttribute.OutputMappings[i]

		for j := range resOutMapping.OutputAttribute.RelatedModuleAttrs {
			moOutAttr := &resOutMapping.OutputAttribute.RelatedModuleAttrs[j]

			if chooseModuleOutputAttr(chosenMoOutAttr, moOutAttr, moInAttr) {
				chosenResOutMapping = resOutMapping
				chosenMoOutAttr = moOutAttr
			}

		}

	}

	if chosenMoOutAttr == nil {
		moInAttr.ResourceAttribute.OutputMappings = []db.TFResourceAttributesMapping{}
		return
	}

	chosenResOutMapping.OutputAttribute.RelatedModuleAttrs = []db.TFModuleAttribute{*chosenMoOutAttr}
	moInAttr.ResourceAttribute.OutputMappings = []db.TFResourceAttributesMapping{*chosenResOutMapping}
	connectedModuleIdSet.Add(chosenMoOutAttr.ModuleID.String())
}

// chooseModuleOutputAttr returns true if the "new output attribute" must be chosen against the previous
// based on module version number
func chooseModuleOutputAttr(previousOutAttr, newOutAttr, inAttr *db.TFModuleAttribute) bool {
	if inAttr.ModuleID == newOutAttr.ModuleID {
		return false
	}

	if previousOutAttr == nil {
		return true
	}

	return previousOutAttr.Module.Version.Compare(newOutAttr.Module.Version) > 0
}

// parseUUIDArr parse array of string UUIDs into `uuid.UUID` type.
// Panics on failure, array elements in response corresponds to the same element indexes from inputs
func parseUUIDArr(uuidStrArr []string) []uuid.UUID {
	res := make([]uuid.UUID, len(uuidStrArr))
	for i, uuidStr := range uuidStrArr {
		res[i] = uuid.MustParse(uuidStr)
	}

	return res
}

// modulesToHclTemplate generate HCL code for multiple module objects
func modulesToHclTemplate(allModules []*db.TFModule, reqModules []string, existingModulesMap map[string]string) (string, error) {
	suggestion := ""

	for _, moduleDef := range allModules {
		if existingModulesMap[moduleDef.Source] != "" && !slices.Contains(reqModules, moduleDef.ID.String()) {
			// don't include module(s) that already exists, unless it's in direct request
			continue
		}

		tfCode, err := moduleToHclTemplate(moduleDef, existingModulesMap)
		if err != nil {
			return "", err
		}

		suggestion = tfCode + suggestion
	}

	return suggestion, nil
}

// moduleToHclTemplate generate HCL code for one terraform module
func moduleToHclTemplate(moduleDef *db.TFModule, existingModulesMap map[string]string) (string, error) {
	// update dependency name from existing module map
	for i, attr := range moduleDef.Attributes {
		if attr.ResourceAttribute != nil && len(attr.ResourceAttribute.OutputMappings) > 0 && len(attr.ResourceAttribute.OutputMappings[0].OutputAttribute.RelatedModuleAttrs) > 0 {
			module := &moduleDef.Attributes[i].ResourceAttribute.OutputMappings[0].OutputAttribute.RelatedModuleAttrs[0].Module
			src := module.Source
			if existingModulesMap[src] != "" {
				module.ModuleName = existingModulesMap[src]
			}
		}
	}

	tfCode := bytes.NewBuffer([]byte{})
	err := goTfTemplates.ExecuteTemplate(tfCode, "module_call", moduleDef)
	if err != nil {
		return "", err
	}

	return tfCode.String(), nil
}
