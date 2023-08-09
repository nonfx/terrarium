package mappings

import (
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/cldcvr/terraform-config-inspect/tfconfig"
	"github.com/cldcvr/terrarium/src/pkg/db"
	"golang.org/x/exp/slices"
)

var resourceTypeByName map[string]*db.TFResourceType

func createMappingRecord(g db.DB, parent *tfconfig.Module, dstRes *tfconfig.Resource, dstResInputName string, srcRes tfconfig.AttributeReference) (*db.TFResourceAttributesMapping, error) {
	if !slices.Contains([]string{"", "module", "var", "local", "each"}, srcRes.Type()) && srcRes.Path() != "" {
		srcResDB, ok := resourceTypeByName[srcRes.Type()]
		if !ok {
			srcResDB = &db.TFResourceType{}
			if err := g.GetTFResourceType(srcResDB, &db.TFResourceType{
				// ProviderID:   provDB.ID,
				ResourceType: srcRes.Type(),
			}); err != nil {
				return nil, nil // skip unknown resources (e.g. need to populate more resource types)
			}
			resourceTypeByName[srcRes.Type()] = srcResDB
		}

		srcAttrDB := &db.TFResourceAttribute{}
		if err := g.GetTFResourceAttribute(srcAttrDB, &db.TFResourceAttribute{
			ResourceTypeID: srcResDB.ID,
			ProviderID:     srcResDB.ProviderID,
			AttributePath:  srcRes.Path(),
		}); err != nil {
			srcFile, srcLine := srcRes.Pos()
			log.Warnf("unknown resource-attribute record %s: %v", fmtAttrMeta(srcRes.Type(), srcRes.Name(), srcRes.Path(), srcFile, srcLine), err)
			return nil, nil // skip unknown resource attributes (e.g. reference to field in a dynamic type)
		}

		dstResDB, ok := resourceTypeByName[dstRes.Type]
		if !ok {
			dstResDB = &db.TFResourceType{}
			if err := g.GetTFResourceType(dstResDB, &db.TFResourceType{
				// ProviderID:   provDB.ID,
				ResourceType: dstRes.Type,
			}); err != nil {
				return nil, nil // skip unknown resources (e.g. need to populate more resource types)
			}
			resourceTypeByName[dstRes.Type] = dstResDB
		}

		dstAttrDB := &db.TFResourceAttribute{}
		if err := g.GetTFResourceAttribute(dstAttrDB, &db.TFResourceAttribute{
			ResourceTypeID: dstResDB.ID,
			ProviderID:     dstResDB.ProviderID,
			AttributePath:  dstResInputName,
		}); err != nil {
			return nil, fmt.Errorf("unknown resource-attribute record %s: %v", fmtAttrMeta(dstRes.Type, dstRes.Name, dstResInputName, dstRes.Pos.Filename, dstRes.Pos.Line), err)
		}

		mappingDB := &db.TFResourceAttributesMapping{
			InputAttributeID:  dstAttrDB.ID,
			OutputAttributeID: srcAttrDB.ID,
		}
		if _, err := g.CreateTFResourceAttributesMapping(mappingDB); err != nil {
			return nil, fmt.Errorf("error creating attribut-mapping record: %v", err)
		}
		return mappingDB, nil
	}
	return nil, nil // skip unresolvable resources
}

func fmtAttrMeta(resType string, resName string, resAttr string, resFile string, resLine int) string {
	return fmt.Sprintf("[resource='%s.%s'; attribute='%s'; file='%s'; line=%d]", resType, resName, resAttr, resFile, resLine)
}

func createMappingsForResources(g db.DB, parent *tfconfig.Module, resources map[string]*tfconfig.Resource, created *[]*db.TFResourceAttributesMapping) (resourceCount int, err error) {
	for dstResName, dstRes := range resources {
		log.Infof("Processing resource declaration '%s'...\n", dstResName)
		resourceMallingCount := 0
		for dstResInputName, inputValueReferences := range dstRes.References {
			for _, item := range inputValueReferences {
				mapping, err := createMappingRecord(g, parent, dstRes, dstResInputName, item)
				if err != nil {
					return 0, err
				}
				*created = append(*created, mapping)
			}
			resourceMallingCount += len(inputValueReferences)
		}
		log.Infof("Created %d mappings for resource declaration '%s'\n", resourceMallingCount, dstResName)
	}
	return len(resources), nil
}

func createMappingsForModule(g db.DB, config *tfconfig.Module) (mappings []*db.TFResourceAttributesMapping, resourceCount int, err error) {
	log.Infof("Processing module '%s'...\n", config.Path)
	mappings = make([]*db.TFResourceAttributesMapping, 0)
	count, err := createMappingsForResources(g, config, config.ManagedResources, &mappings)
	if err != nil {
		return nil, 0, err
	}
	resourceCount += count

	count, err = createMappingsForResources(g, config, config.DataResources, &mappings)
	if err != nil {
		return nil, 0, err
	}
	resourceCount += count

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