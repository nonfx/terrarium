package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/cldcvr/terrarium/api/db"
	"github.com/cldcvr/terrarium/api/pkg/terraform-config-inspect/tfconfig"
)

const moduleSchemaFilePath = "terraform/.terraform/modules/modules.json"

var resourceTypeByName map[string]*db.TFResourceType

type tfValue interface {
	GetName() string
	GetDescription() string
	IsRequired() bool
	IsComputed() bool
}

func createAttributeRecord(g db.DB, moduleDB *db.TFModule, v tfValue, varAttributePath string, res tfconfig.ResourceAttributeReference) (*db.TFModuleAttribute, error) {
	resDB, ok := resourceTypeByName[res.ResourceType]
	if !ok {
		if res.ResourceType == "module" {
			return nil, nil // module reference was not resolved - the parser does not support remote module references
		} else if res.ResourceType == "var" {
			return nil, nil // output returns another variable (i.e. passthrough) and does not directly map to a resource
		}
		resDB = &db.TFResourceType{}
		if err := g.GetTFResourceType(resDB, &db.TFResourceType{
			// ProviderID:   provDB.ID,  // there may be more than a single provider (e.g. random_password)
			ResourceType: res.ResourceType,
		}); err != nil {
			return nil, nil // skip unkown resources (e.g. need to populate more resource types)
		}
		resourceTypeByName[res.ResourceType] = resDB
	}

	resourceAttrDB := &db.TFResourceAttribute{}
	if err := g.GetTFResourceAttribute(resourceAttrDB, &db.TFResourceAttribute{
		ResourceTypeID: resDB.ID,
		ProviderID:     resDB.ProviderID,
		AttributePath:  strings.Join(res.AttributePath, "."),
	}); err != nil {
		// VAN-4158: the exact match on path may fail, we need to match by prefix instead
		// we store all sub-paths such as 'rule.noncurrent_version_expiration.newer_noncurrent_versions'
		// but the output or input variable may be refering to 'rule.noncurrent_version_expiration' portion only
		// we should treat each "path-level" as an attribute of its own - other modules may use it as input
		// return nil, fmt.Errorf("unknown resource-attribute record: %v", err)
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
		return nil, fmt.Errorf("error creating module-attribute record: %v", err)
	}

	return moduleAttrDB, nil
}

func main() {
	g, err := db.Connect()
	if err != nil {
		panic(err)
	}

	resourceTypeByName = make(map[string]*db.TFResourceType)

	// load modules
	log.Println("Loading modules...")
	configs, _, err := tfconfig.LoadModulesFromResolvedSchema(moduleSchemaFilePath)
	if err != nil {
		panic(err)
	}
	log.Printf("Loaded %d modules\n", len(configs))

	for _, config := range configs {
		log.Printf("Processing module '%s'...\n", config.Path)

		moduleDB := toTFModule(config)
		if _, err := g.CreateTFModule(moduleDB); err != nil {
			log.Println("Error creating module record:", err)
			return
		}

		for varName, v := range config.Variables {
			if varAttrReferences, ok := config.Inputs[varName]; ok { // found a resolution for this variable to resource attribute
				for varAttributePath, resourceReferences := range varAttrReferences {
					for _, res := range resourceReferences {
						if attr, err := createAttributeRecord(g, moduleDB, v, varAttributePath, res); err != nil {
							log.Println("Error creating module input-attribute:", err)
							return
						} else if attr == nil {
							continue
						}
					}
				}
			}
		}

		for _, o := range config.Outputs {
			if attr, err := createAttributeRecord(g, moduleDB, o, "", o.Value); err != nil {
				log.Println("Error creating module output-attribute:", err)
				return
			} else if attr == nil {
				continue
			}
		}
		log.Printf("Module '%s' done processing\n", config.Path)
	}

}

func toTFModule(config *tfconfig.Module) *db.TFModule {
	record := db.TFModule{
		ModuleName: config.Path,
		Source:     config.Path,
	}
	if config.Metadata != nil {
		record.ModuleName = config.Metadata.Name
		record.Source = config.Metadata.Source
		record.Version = config.Metadata.Version
	}
	return &record
}
