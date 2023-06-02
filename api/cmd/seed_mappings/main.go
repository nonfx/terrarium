package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/cldcvr/terrarium/api/db"
	"github.com/cldcvr/terrarium/api/pkg/terraform-config-inspect/tfconfig"
	"golang.org/x/exp/slices"
)

const moduleSchemaFilePath = "terraform/.terraform/modules/modules.json"

var resourceTypeByName map[string]*db.TFResourceType

func createMappingsForResourceInputs(g db.DB, dstRes *tfconfig.Resource) ([]*db.TFResourceAttributesMapping, error) {
	mappings := make([]*db.TFResourceAttributesMapping, 0)
	// provDB := &db.TFProvider{}
	// if err := g.GetTFProvider(provDB, &db.TFProvider{
	// 	Name: res.Provider.Name,
	// }); err != nil {
	// 	return nil, nil // skip unkown resources (e.g. need to populate more resource types)
	// }
	for dstResInputName, srcRes := range dstRes.Inputs {
		if !slices.Contains([]string{"", "module", "var", "local", "each"}, srcRes.ResourceType) {
			srcResDB, ok := resourceTypeByName[srcRes.ResourceType]
			if !ok {
				srcResDB = &db.TFResourceType{}
				if err := g.GetTFResourceType(srcResDB, &db.TFResourceType{
					// ProviderID:   provDB.ID,
					ResourceType: srcRes.ResourceType,
				}); err != nil {
					return nil, nil // skip unkown resources (e.g. need to populate more resource types)
				}
				resourceTypeByName[srcRes.ResourceType] = srcResDB
			}

			srcAttrDB := &db.TFResourceAttribute{}
			if err := g.GetTFResourceAttribute(srcAttrDB, &db.TFResourceAttribute{
				ResourceTypeID: srcResDB.ID,
				ProviderID:     srcResDB.ProviderID,
				AttributePath:  strings.Join(srcRes.AttributePath, "."),
			}); err != nil {
				return nil, fmt.Errorf("unknown resource-attribute record: %v", err)
			}

			dstResDB, ok := resourceTypeByName[dstRes.Type]
			if !ok {
				dstResDB = &db.TFResourceType{}
				if err := g.GetTFResourceType(dstResDB, &db.TFResourceType{
					// ProviderID:   provDB.ID,
					ResourceType: dstRes.Type,
				}); err != nil {
					return nil, nil // skip unkown resources (e.g. need to populate more resource types)
				}
				resourceTypeByName[dstRes.Type] = dstResDB
			}

			dstAttrDB := &db.TFResourceAttribute{}
			if err := g.GetTFResourceAttribute(dstAttrDB, &db.TFResourceAttribute{
				ResourceTypeID: dstResDB.ID,
				ProviderID:     dstResDB.ProviderID,
				AttributePath:  dstResInputName,
			}); err != nil {
				return nil, fmt.Errorf("unknown resource-attribute record: %v", err)
			}

			mappingDB := &db.TFResourceAttributesMapping{
				InputAttributeID:  dstAttrDB.ID,
				OutputAttributeID: srcAttrDB.ID,
			}
			if _, err := g.CreateTFResourceAttributesMapping(mappingDB); err != nil {
				return nil, fmt.Errorf("error creating attribut-mapping record: %v", err)
			}
			mappings = append(mappings, mappingDB)
		}
	}
	return mappings, nil
}

func createMappingsForModule(g db.DB, config *tfconfig.Module) (mappings []*db.TFResourceAttributesMapping, resourceCount int, err error) {
	log.Printf("Processing module '%s'...\n", config.Path)
	for dstResName, dstRes := range config.ManagedResources {
		log.Printf("Processing resource declaration '%s'...\n", dstResName)
		mappings, err = createMappingsForResourceInputs(g, dstRes)
		if err != nil {
			panic(err)
		}
		mappingsCreatedCount := len(mappings)
		log.Printf("Created %d mappings for resource declaration '%s'\n", mappingsCreatedCount, dstResName)
	}
	resourceCount += len(config.ManagedResources)

	// process sub-modules
	for _, dstMod := range config.ModuleCalls {
		if dstMod.Module != nil {
			subMappings, subResCount, err := createMappingsForModule(g, dstMod.Module)
			if err != nil {
				return subMappings, subResCount, err
			}
			mappings = append(mappings, subMappings...)
			resourceCount += subResCount
		}
	}
	return
}

func main() {
	g, err := db.Connect()
	if err != nil {
		panic(err)
	}

	resourceTypeByName = make(map[string]*db.TFResourceType)
	log.Println("Loading modules...")
	configs, _, err := tfconfig.LoadModulesFromResolvedSchema(moduleSchemaFilePath)
	if err != nil {
		panic(err)
	}
	moduleCount := len(configs)
	log.Printf("Loaded %d modules\n", moduleCount)

	totalResourceDeclarationsCount := 0
	totalMappingsCreatedCount := 0
	for _, config := range configs {
		mappings, resourceCount, err := createMappingsForModule(g, config)
		if err != nil {
			panic(err)
		}
		totalResourceDeclarationsCount += resourceCount
		totalMappingsCreatedCount += len(mappings)
	}
	log.Printf("Processed %d resource declarations in %d modules and created %d mappings...\n", totalResourceDeclarationsCount, moduleCount, totalMappingsCreatedCount)
}
