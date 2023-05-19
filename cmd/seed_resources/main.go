package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/cldcvr/terrarium/api/db"
	"github.com/cldcvr/terrarium/api/pkg/tf/schema"
)

func main() {
	// Connect to the database
	db, err := db.Connect()
	mustNotErr(err, "Error connecting to the database")

	// Load providers schema from file
	providersSchema, err := loadProvidersSchema("cache_data/tf_resources.json")
	mustNotErr(err, "Error loading providers schema")

	pushProvidersSchemaToDB(providersSchema, db)
}

func pushProvidersSchemaToDB(providersSchema *schema.ProvidersSchema, dbConn db.DB) {
	// Process each provider in the schema
	for providerName, resources := range providersSchema.ProviderSchemas {
		// Create a new TFProvider instance
		provider := &db.TFProvider{
			Name: providerName,
		}
		providerID, err := dbConn.CreateTFProvider(provider)
		mustNotErr(err, "Error creating provider: %s", providerName)
		log.Printf("Provider created: %s\t%s\n", providerID, providerName)

		// Process each resource type in the provider
		for resourceType, resourceSchema := range resources.ResourceSchemas {
			// Create a new TFResourceType instance
			resource := &db.TFResourceType{
				ProviderID:   providerID,
				ResourceType: resourceType,
			}
			resourceID, err := dbConn.CreateTFResourceType(resource)
			mustNotErr(err, "Error creating resource type: %s", resourceType)
			log.Printf("\tResource type created: %s\t%s\n", resourceID, resourceType)

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
				mustNotErr(err, "Error creating attribute: %s", attributePath)
				log.Printf("\t\tAttribute created: %s\t%s\n", attrID, attributePath)
			}
		}
	}
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

// mustNotErr checks if the error is non-nil and panics if it is, logging the provided message
func mustNotErr(err error, msg string, args ...interface{}) {
	if err != nil {
		log.Printf(msg, args...)
		panic(err)
	}
}
