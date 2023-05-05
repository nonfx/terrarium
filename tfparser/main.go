package main

import (
	"fmt"
	"strings"

	"github.com/cldcvr/terrarium/tfparser/lib/terraform-config-inspect/tfconfig"
)

func loadAttributesByResource(config *tfconfig.Module, attrByResource map[string]map[string]map[string]string) {
	for name, items := range config.Inputs {
		for _, item := range items {
			attrsForResource, ok := attrByResource[item.ResourceType+"."+item.ResourceName]
			if !ok {
				attrsForResource = make(map[string]map[string]string)
				attrByResource[item.ResourceType+"."+item.ResourceName] = attrsForResource
			}

			inputs, ok := attrsForResource["inputs"]
			if !ok {
				inputs = make(map[string]string)
				attrsForResource["inputs"] = inputs // resource-attribute: input-variable
			}

			inputs[strings.Join(item.AttributePath, ".")] = name
		}
	}
	for name, item := range config.Outputs {
		attrsForResource, ok := attrByResource[item.Value.ResourceType+"."+item.Value.ResourceName]
		if !ok {
			attrsForResource = make(map[string]map[string]string)
			attrByResource[item.Value.ResourceType+"."+item.Value.ResourceName] = attrsForResource
		}

		outputs, ok := attrsForResource["outputs"]
		if !ok {
			outputs = make(map[string]string)
			attrsForResource["outputs"] = outputs // resource-attribute: module-output
		}

		outputs[strings.Join(item.Value.AttributePath, ".")] = name
	}
}

func parseModule(modulePath string) {
	config, _ := tfconfig.LoadModule(modulePath)
	// if diag.HasErrors() {
	// 	log.Fatalf(diag.Error())
	// }
	attrByResource := make(map[string]map[string]map[string]string) // resource-type.resource-name: attribute_kind: resource-attribute-name: attribute_path
	// fmt.Print("Resources:\n")
	loadAttributesByResource(config, attrByResource)
	// for _, subModuleCall := range config.ModuleCalls {
	// 	loadAttributesByResource(subModuleCall.Module, attrByResource)
	// }

	fmt.Printf("Module Input, Resource Input Attribute, Resource, Resource Output Attribute, Module Output\n")
	for resourceKey, r1 := range attrByResource {
		for resourceAttributeName, moduleInputName := range r1["inputs"] {
			fmt.Printf("%s, %s, %s, %s, %s\n", moduleInputName, resourceAttributeName, resourceKey, "", "")
		}
		for resourceAttributeName, moduleOutputName := range r1["outputs"] {
			fmt.Printf("%s, %s, %s, %s, %s\n", "", "", resourceKey, resourceAttributeName, moduleOutputName)
		}
	}
}

func main() {
	for _, p := range []string{
		"./modules/terraform-aws-vpc",
		"./modules/terraform-aws-rds",
		"./modules/terraform-aws-security-group",
		"./modules/terraform-aws-eks",
		"./modules/terraform-aws-s3-bucket",
	} {
		fmt.Printf("MODULE: %s\n", p)
		parseModule(p)
		fmt.Printf("\n\n")
	}
}
