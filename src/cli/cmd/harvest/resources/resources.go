// Copyright (c) Ollion
// SPDX-License-Identifier: Apache-2.0

package resources

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/charmbracelet/log"
	"github.com/cldcvr/terrarium/src/pkg/db"
	"github.com/cldcvr/terrarium/src/pkg/tf/schema"
	"github.com/google/uuid"
	"github.com/rotisserie/eris"
)

func pushProvidersSchemaToDB(providersSchema *schema.ProvidersSchema, dbConn db.DB) (providerCount, allResCount, allAttrCount int, err error) {
	// Process each provider in the schema
	for providerName, resources := range providersSchema.ProviderSchemas {
		providerID, isNew, err := dbConn.GetOrCreateTFProvider(&db.TFProvider{Name: providerName})
		if err != nil {
			return providerCount, allResCount, allAttrCount, eris.Wrapf(err, "error creating provider: %s", providerName)
		}

		if !isNew {
			log.Info("Provider already exists, skipping", "name", providerName, "id", providerID)
			continue
		}

		log.Info("Provider created", "name", providerName, "id", providerID)
		providerCount++

		// Process each resource and data-resource type in the provider
		resCount, attrCount, err := pushSchemasToDB(dbConn, providerID, resources.ResourceSchemas)
		if err != nil {
			return providerCount, allResCount, allAttrCount, err
		}
		allResCount += resCount
		allAttrCount += attrCount

		resCount, attrCount, err = pushSchemasToDB(dbConn, providerID, resources.DataSourceSchemas)
		if err != nil {
			return providerCount, allResCount, allAttrCount, err
		}
		allResCount += resCount
		allAttrCount += attrCount
	}

	return
}

// loadProvidersSchema loads the providers schema from a file
func loadProvidersSchema(filename string) (*schema.ProvidersSchema, error) {
	// Read the file content
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// Unmarshal the JSON data into ProvidersSchema struct
	var providersSchema schema.ProvidersSchema
	err = json.Unmarshal(data, &providersSchema)
	if err != nil {
		return nil, err
	}

	return &providersSchema, nil
}

func pushSchemasToDB(dbConn db.DB, providerID uuid.UUID, schemas map[string]schema.SchemaRepresentation) (resCount, attrCount int, err error) {
	for resourceType, resourceSchema := range schemas {
		// Create a new TFResourceType instance
		resource := &db.TFResourceType{
			ProviderID:   providerID,
			ResourceType: resourceType,
		}

		resourceID, err := dbConn.CreateTFResourceType(resource)
		if err != nil {
			return resCount, attrCount, eris.Wrapf(err, "error creating resource type: %s", resourceType)
		}

		log.Debug("Resource type created", "name", resourceType, "id", resourceID)
		resCount++

		// Process each attribute in the resource type
		attributes := resourceSchema.Block.ListLeafNodes()
		for attributePath, attribute := range attributes {
			// Create a new TFResourceAttribute instance
			attr := &db.TFResourceAttribute{
				ResourceTypeID: resourceID,
				ProviderID:     providerID,
				AttributePath:  attributePath,
				DataType:       fmt.Sprintf("%v", attribute.Type),
				Description:    attribute.Description,
				Optional:       attribute.Optional,
				Computed:       attribute.Computed,
			}

			attrID, err := dbConn.CreateTFResourceAttribute(attr)
			if err != nil {
				return resCount, attrCount, eris.Wrapf(err, "error creating attribute: %s", attributePath)
			}

			log.Debug("Attribute created", "name", attributePath, "id", attrID)
			attrCount++
		}
	}

	return
}
