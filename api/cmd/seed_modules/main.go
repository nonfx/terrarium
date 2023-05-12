package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/cldcvr/terrarium/api/db"
	"github.com/cldcvr/terrarium/api/pkg/terraform-config-inspect/tfconfig"
	"gorm.io/gorm"
)

var resourceTypeByName map[string]*db.TFResourceType

type tfValue interface {
	GetName() string
	GetDescription() string
	IsRequired() bool
	IsComputed() bool
}

func createAttributeRecord(g *gorm.DB, moduleDB *db.TFModule, v tfValue, varAttributePath string, res tfconfig.ResourceAttributeReference) (*db.TFModuleAttribute, error) {
	resDB, ok := resourceTypeByName[res.ResourceType]
	if !ok {
		if res.ResourceType == "module" {
			return nil, nil // module reference was not resolved - the parser does not support remote module references
		}
		resDB = &db.TFResourceType{}
		if err := g.First(resDB, db.TFResourceType{
			// ProviderID:   provDB.ID,  // there may be more than a single provider (e.g. random_password)
			ResourceType: res.ResourceType,
		}).Error; err != nil {
			return nil, nil // skip unkown resources (e.g. need to populate more resource types)
		}
		resourceTypeByName[res.ResourceType] = resDB
	}

	resourceAttrDB := &db.TFResourceAttribute{}
	if err := g.First(&resourceAttrDB, db.TFResourceAttribute{
		ResourceTypeID: resDB.ID,
		ProviderID:     resDB.ProviderID,
		AttributePath:  strings.Join(res.AttributePath, "."),
	}).Error; err != nil {
		return nil, fmt.Errorf("Unknown resource-attribute record: %v", err)
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

	if err := moduleAttrDB.Create(g); err != nil {
		return nil, fmt.Errorf("Error creating module-attribute record: %v", err)
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
	for _, moduleSrcDir := range []string{
		"./cache_data/modules/terraform-aws-vpc",
		"./cache_data/modules/terraform-aws-rds",
		"./cache_data/modules/terraform-aws-security-group",
		"./cache_data/modules/terraform-aws-eks",
		"./cache_data/modules/terraform-aws-s3-bucket",
	} {
		config, _ := tfconfig.LoadModule(moduleSrcDir)

		moduleDB := &db.TFModule{
			ModuleName: config.Path,
			Source:     config.Path,
		}
		if err := moduleDB.Create(g); err != nil {
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
	}

}
