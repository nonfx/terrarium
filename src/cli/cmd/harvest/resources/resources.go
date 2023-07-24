package resources

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/cldcvr/terrarium/src/cli/internal/config"
	"github.com/cldcvr/terrarium/src/pkg/db"
	"github.com/cldcvr/terrarium/src/pkg/tf/schema"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var moduleSchemaFilePathFlag string

const DefaultSchemaPath = ".terraform/providers/schema.json"

var resourcesCmd = &cobra.Command{
	Use:   "resources",
	Short: "Harvests Terraform resources and attributes using the provider schema json",
	Long: heredoc.Docf(`
		This command requires terraform provider schema already generated. To do that, run:
			terraform providers schema -json > %s
	`, DefaultSchemaPath),
	Run: func(cmd *cobra.Command, args []string) {
		main()
	},
}

func GetCmd() *cobra.Command {
	addFlags()
	return resourcesCmd
}

func addFlags() {
	resourcesCmd.Flags().StringVarP(&moduleSchemaFilePathFlag, "file", "f", DefaultSchemaPath, "terraform provider schema json file path")
}

func main() {
	// Connect to the database
	db, err := config.DBConnect()
	mustNotErr(err, "Error connecting to the database")

	// Load providers schema from file
	providersSchema, err := loadProvidersSchema(moduleSchemaFilePathFlag)
	mustNotErr(err, "Error loading providers schema")

	pushProvidersSchemaToDB(providersSchema, db)
}

func pushProvidersSchemaToDB(providersSchema *schema.ProvidersSchema, dbConn db.DB) {
	// Process each provider in the schema
	for providerName, resources := range providersSchema.ProviderSchemas {
		providerID, isNew, err := dbConn.GetOrCreateTFProvider(&db.TFProvider{Name: providerName})
		mustNotErr(err, "Error creating provider: %s", providerName)

		if !isNew {
			log.Printf("Provider already exists, skipping resource seed: %s | %s\n", providerID, providerName)
			continue
		}

		log.Printf("Provider created: %s | %s\n", providerID, providerName)

		// Process each resource and data-resource type in the provider
		pushSchemasToDB(dbConn, providerID, resources.ResourceSchemas)
		pushSchemasToDB(dbConn, providerID, resources.DataSourceSchemas)
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

func pushSchemasToDB(dbConn db.DB, providerID uuid.UUID, schemas map[string]schema.SchemaRepresentation) {
	for resourceType, resourceSchema := range schemas {
		// Create a new TFResourceType instance
		resource := &db.TFResourceType{
			ProviderID:   providerID,
			ResourceType: resourceType,
		}
		resourceID, err := dbConn.CreateTFResourceType(resource)
		mustNotErr(err, "Error creating resource type: %s", resourceType)
		log.Printf("\tResource type created: %s | %s\n", resourceID, resourceType)

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
			log.Printf("\t\tAttribute created: %s | %s\n", attrID, attributePath)
		}
	}
}
