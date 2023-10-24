// Copyright (c) Ollion
// SPDX-License-Identifier: Apache-2.0

package mappings

import (
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/cldcvr/terraform-config-inspect/tfconfig"
	"github.com/cldcvr/terrarium/src/pkg/db"
	"github.com/rotisserie/eris"
	"golang.org/x/exp/slices"
)

func createMappingRecord(g db.DB, parent *tfconfig.Module, dstRes *tfconfig.Resource, dstResInputName string, srcRes tfconfig.AttributeReference, resourceTypeByNameCache map[string]*db.TFResourceType) (*db.TFResourceAttributesMapping, error) {
	if !slices.Contains([]string{"", "module", "var", "local", "each"}, srcRes.Type()) && srcRes.Path() != "" {
		srcResDB, ok := resourceTypeByNameCache[srcRes.Type()]
		if !ok {
			srcResDB = &db.TFResourceType{}
			if err := g.GetTFResourceType(srcResDB, &db.TFResourceType{
				// ProviderID:   provDB.ID,
				ResourceType: srcRes.Type(),
			}); err != nil {
				return nil, nil // skip unknown resources (e.g. need to populate more resource types)
			}
			resourceTypeByNameCache[srcRes.Type()] = srcResDB
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

		dstResDB, ok := resourceTypeByNameCache[dstRes.Type]
		if !ok {
			dstResDB = &db.TFResourceType{}
			if err := g.GetTFResourceType(dstResDB, &db.TFResourceType{
				// ProviderID:   provDB.ID,
				ResourceType: dstRes.Type,
			}); err != nil {
				return nil, nil // skip unknown resources (e.g. need to populate more resource types)
			}
			resourceTypeByNameCache[dstRes.Type] = dstResDB
		}

		dstAttrDB := &db.TFResourceAttribute{}
		if err := g.GetTFResourceAttribute(dstAttrDB, &db.TFResourceAttribute{
			ResourceTypeID: dstResDB.ID,
			ProviderID:     dstResDB.ProviderID,
			AttributePath:  dstResInputName,
		}); err != nil {
			return nil, eris.Wrapf(err, "unknown resource-attribute record %s", fmtAttrMeta(dstRes.Type, dstRes.Name, dstResInputName, dstRes.Pos.Filename, dstRes.Pos.Line))
		}

		mappingDB := &db.TFResourceAttributesMapping{
			InputAttributeID:  dstAttrDB.ID,
			OutputAttributeID: srcAttrDB.ID,
		}
		if _, err := g.CreateTFResourceAttributesMapping(mappingDB); err != nil {
			return nil, eris.Wrap(err, "error creating attribut-mapping record")
		}
		return mappingDB, nil
	}
	return nil, nil // skip unresolvable resources
}

func fmtAttrMeta(resType string, resName string, resAttr string, resFile string, resLine int) string {
	return fmt.Sprintf("[resource='%s.%s'; attribute='%s'; file='%s'; line=%d]", resType, resName, resAttr, resFile, resLine)
}

func createMappingsForResources(g db.DB, parent *tfconfig.Module, resources map[string]*tfconfig.Resource, created *[]*db.TFResourceAttributesMapping, resourceTypeByNameCache map[string]*db.TFResourceType) (resourceCount int, err error) {
	for dstResName, dstRes := range resources {
		log.Infof("Processing resource declaration '%s'...", dstResName)
		resourceMappingCount := 0
		for dstResInputName, inputValueReferences := range dstRes.References {
			for _, item := range inputValueReferences {
				mapping, err := createMappingRecord(g, parent, dstRes, dstResInputName, item, resourceTypeByNameCache)
				if err != nil {
					return 0, err
				} else if mapping != nil {
					*created = append(*created, mapping)
					resourceMappingCount += 1
				}
			}
		}
		log.Infof("Created %d mappings for resource declaration '%s'", resourceMappingCount, dstResName)
	}
	return len(resources), nil
}

func createMappingsForModule(g db.DB, config *tfconfig.Module, resourceTypeByNameCache map[string]*db.TFResourceType) (mappings []*db.TFResourceAttributesMapping, resourceCount int, err error) {
	log.Infof("Processing module '%s'...", config.Path)
	mappings = make([]*db.TFResourceAttributesMapping, 0)
	count, err := createMappingsForResources(g, config, config.ManagedResources, &mappings, resourceTypeByNameCache)
	if err != nil {
		return nil, 0, err
	}
	resourceCount += count

	count, err = createMappingsForResources(g, config, config.DataResources, &mappings, resourceTypeByNameCache)
	if err != nil {
		return nil, 0, err
	}
	resourceCount += count

	// process sub-modules
	for _, dstMod := range config.ModuleCalls {
		if dstMod.Module != nil {
			subMappings, subResCount, err := createMappingsForModule(g, dstMod.Module, resourceTypeByNameCache)
			if err != nil {
				return subMappings, subResCount, err
			}
			mappings = append(mappings, subMappings...)
			resourceCount += subResCount
		}
	}
	return
}
