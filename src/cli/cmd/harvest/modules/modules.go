// Copyright (c) Ollion
// SPDX-License-Identifier: Apache-2.0

package modules

import (
	"os"
	"strings"

	"github.com/cldcvr/terraform-config-inspect/tfconfig"
	cliconfig "github.com/cldcvr/terrarium/src/cli/internal/config"
	"github.com/cldcvr/terrarium/src/pkg/db"
	"github.com/rotisserie/eris"
)

type tfValue interface {
	GetName() string
	GetDescription() string
	IsRequired() bool
	IsComputed() bool
}

func createAttributeRecord(g db.DB, moduleDB *db.TFModule, v tfValue, varAttributePath string, res tfconfig.AttributeReference, resourceTypeByNameCache map[string]*db.TFResourceType) (*db.TFModuleAttribute, error) {
	if res.Type() == "module" {
		return nil, nil // module reference was not resolved - the parser does not support remote module references
	} else if res.Type() == "var" {
		return nil, nil // output returns another variable (i.e. passthrough) and does not directly map to a resource
	} else if res.Path() == "" {
		return nil, nil // returning the entire resource
	}

	resDB, ok := resourceTypeByNameCache[res.Type()]
	if !ok {
		resDB = &db.TFResourceType{}
		if err := g.GetTFResourceType(resDB, &db.TFResourceType{
			// ProviderID:   provDB.ID,  // there may be more than a single provider (e.g. random_password)
			ResourceType: res.Type(),
		}); err != nil {
			return nil, nil // skip unkown resources (e.g. need to populate more resource types)
		}
		resourceTypeByNameCache[res.Type()] = resDB
	}

	resourceAttrDB := &db.TFResourceAttribute{}
	if err := g.GetTFResourceAttribute(resourceAttrDB, &db.TFResourceAttribute{
		ResourceTypeID: resDB.ID,
		ProviderID:     resDB.ProviderID,
		AttributePath:  res.Path(),
	}); err != nil {
		// VAN-4158: the exact match on path may fail, we need to match by prefix instead
		// we store all sub-paths such as 'rule.noncurrent_version_expiration.newer_noncurrent_versions'
		// but the output or input variable may be refering to 'rule.noncurrent_version_expiration' portion only
		// we should treat each "path-level" as an attribute of its own - other modules may use it as input
		// return nil, eris.Wrap(err, "unknown resource-attribute record")
		return nil, nil
	}

	moduleAttrPathTokens := []string{v.GetName()}
	if varAttributePath != "" {
		moduleAttrPathTokens = append(moduleAttrPathTokens, varAttributePath)
	}

	moduleAttrDB := &db.TFModuleAttribute{
		ModuleID:                       moduleDB.ID,
		ModuleAttributeName:            strings.Join(moduleAttrPathTokens, "."),
		Description:                    v.GetDescription(),
		RelatedResourceTypeAttributeID: resourceAttrDB.ID,
		Optional:                       !v.IsRequired(),
		Computed:                       v.IsComputed(),
	}

	if _, err := g.CreateTFModuleAttribute(moduleAttrDB); err != nil {
		return nil, eris.Wrap(err, "error creating module-attribute record")
	}

	return moduleAttrDB, nil
}

func toTFModule(config *tfconfig.Module) *db.TFModule {
	record := db.TFModule{
		ModuleName: config.Path,
		Source:     config.Path,
		Namespace:  cliconfig.FarmDefault(),
	}
	if config.Metadata != nil {
		record.ModuleName = config.Metadata.Name
		record.Source = config.Metadata.Source
		// filter local module
		if flagIncludeLocal && strings.HasPrefix(config.Metadata.Source, ".") && config.Metadata.Name != "" {
			if strings.TrimSpace(flagTFDir) == "." || strings.TrimSpace(flagTFDir) == "" {
				cwd, _ := os.Getwd()
				flagTFDir = strings.TrimRight(cwd, "\n")
			}
			record.Namespace = flagTFDir
		}
		record.Version = db.Version(config.Metadata.Version)
	}
	return &record
}
