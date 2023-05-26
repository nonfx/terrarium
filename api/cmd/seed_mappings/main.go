package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/cldcvr/terrarium/api/db"
	"github.com/cldcvr/terrarium/api/pkg/terraform-config-inspect/tfconfig"
	"golang.org/x/exp/slices"
)

const linkingModulePath = "terraform/"

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

func main() {
	g, err := db.Connect()
	if err != nil {
		panic(err)
	}

	resourceTypeByName = make(map[string]*db.TFResourceType)
	log.Println("Loading root module...")
	config, _ := tfconfig.LoadModule(linkingModulePath, &tfconfig.ResolvedModulesSchema{}) // no module schema: do not descend into sub-modules
	log.Println("Loaded 1 module")

	log.Println("Processing resource declarations...")
	for dstResName, dstRes := range config.ManagedResources {
		log.Printf("Processing resource declaration '%s'...\n", dstResName)
		mappings, err := createMappingsForResourceInputs(g, dstRes)
		if err != nil {
			panic(err)
		}
		log.Printf("Created %d mappings for resource declaration '%s'\n", len(mappings), dstResName)
	}
	log.Printf("Processed %d resource declarations...\n", len(config.ManagedResources))

}
