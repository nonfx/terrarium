package modules

import (
	"fmt"
	"log"
	"strings"

	"github.com/cldcvr/terraform-config-inspect/tfconfig"
	"github.com/cldcvr/terrarium/src/cli/internal/config"
	"github.com/cldcvr/terrarium/src/pkg/db"
	"github.com/spf13/cobra"
)

const moduleSchemaFilePath = "terraform/.terraform/modules/modules.json"

var resourceTypeByName map[string]*db.TFResourceType

type tfValue interface {
	GetName() string
	GetDescription() string
	IsRequired() bool
	IsComputed() bool
}

var modulesCmd = &cobra.Command{
	Use:   "modules",
	Short: "Scrapes Terraform modules and attributes from the farm directory",
	Long:  "The 'modules' command scrapes all Terraform modules and their attributes from the specified farm directory.",
	Run: func(cmd *cobra.Command, args []string) {
		main()
	},
}

func GetCmd() *cobra.Command {
	return modulesCmd
}

func createAttributeRecord(g db.DB, moduleDB *db.TFModule, v tfValue, varAttributePath string, res tfconfig.AttributeReference) (*db.TFModuleAttribute, error) {
	if res.Type() == "module" {
		return nil, nil // module reference was not resolved - the parser does not support remote module references
	} else if res.Type() == "var" {
		return nil, nil // output returns another variable (i.e. passthrough) and does not directly map to a resource
	} else if res.Path() == "" {
		return nil, nil // returning the entire resource
	}

	resDB, ok := resourceTypeByName[res.Type()]
	if !ok {
		resDB = &db.TFResourceType{}
		if err := g.GetTFResourceType(resDB, &db.TFResourceType{
			// ProviderID:   provDB.ID,  // there may be more than a single provider (e.g. random_password)
			ResourceType: res.Type(),
		}); err != nil {
			return nil, nil // skip unkown resources (e.g. need to populate more resource types)
		}
		resourceTypeByName[res.Type()] = resDB
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
	g, err := config.DBConnect()
	if err != nil {
		panic(err)
	}

	resourceTypeByName = make(map[string]*db.TFResourceType)

	// load modules
	log.Println("Loading modules...")

	configs, _, err := tfconfig.LoadModulesFromResolvedSchema(moduleSchemaFilePath, tfconfig.FilterModulesOmitLocal, tfconfig.FilterModulesOmitHidden)
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
		record.Version = db.Version(config.Metadata.Version)
	}
	return &record
}
