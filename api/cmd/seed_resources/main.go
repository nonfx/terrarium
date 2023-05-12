package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"log"

	"github.com/cldcvr/terrarium/api/db"
	"github.com/cldcvr/terrarium/api/pkg/tf/schema"
)

// generate using `terraform providers schema -json > tf_aws_resources.json`
//
//go:embed cache_data/tf_aws_resources.json
var tf_aws_resources embed.FS

func main() {
	g, err := db.Connect()
	if err != nil {
		panic(err)
	}

	var config schema.ProvidersSchema

	data, err := tf_aws_resources.ReadFile("cache_data/tf_aws_resources.json")
	if err != nil {
		log.Println("Error reading config file:", err)
		return
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Println("Error unmarshaling config data:", err)
		return
	}

	// load resources
	for providerName, resources := range config.ProviderSchemas {
		provider := &db.TFProvider{
			Name: providerName,
		}
		provider.Create(g)
		log.Printf("%d\t%s\n", provider.ID, providerName)

		for resType, block := range resources.ResourceSchemas {
			// resName := "aws_elastic_beanstalk_environment"
			// block := resources.ResourceSchemas[resName]
			res := &db.TFResourceType{
				ProviderID:   provider.ID,
				ResourceType: resType,
			}
			res.Create(g)
			log.Printf("\tResource: %d\t%s\n", res.ID, resType)

			attrs := block.Block.ListLeafNodes()
			for attrPath, attr := range attrs {
				attrDB := &db.TFResourceAttribute{
					ResourceTypeID: res.ID,
					ProviderID:     provider.ID,
					AttributePath:  attrPath,
					DataType:       fmt.Sprintf("%v", attr.Type),
					Description:    attr.Description,
					Optional:       attr.Optional,
					Computed:       attr.Computed,
				}
				attrDB.Create(g)
				log.Printf("\t\tAttribute: %d\t%s\n", attrDB.ID, attrPath)
			}

			// break
		}

		for resType, block := range resources.DataSourceSchemas {
			// resName := "aws_elastic_beanstalk_environment"
			// block := resources.ResourceSchemas[resName]
			res := &db.TFResourceType{
				ProviderID:   provider.ID,
				ResourceType: resType,
			}
			res.Create(g)
			log.Printf("\tResource: %d\t%s\n", res.ID, resType)

			attrs := block.Block.ListLeafNodes()
			for attrPath, attr := range attrs {
				attrDB := &db.TFResourceAttribute{
					ResourceTypeID: res.ID,
					ProviderID:     provider.ID,
					AttributePath:  attrPath,
					DataType:       fmt.Sprintf("%v", attr.Type),
					Description:    attr.Description,
					Optional:       attr.Optional,
					Computed:       attr.Computed,
				}
				attrDB.Create(g)
				log.Printf("\t\tAttribute: %d\t%s\n", attrDB.ID, attrPath)
			}

			// break
		}

	}

}
