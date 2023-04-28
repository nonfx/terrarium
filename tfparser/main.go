package main

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-config-inspect/tfconfig"
)

func main() {
	config, _ := tfconfig.LoadModule("./modules/terraform-aws-vpc")
	// if diag.HasErrors() {
	// 	log.Fatalf(diag.Error())
	// }
	fmt.Printf("Module Input, Resource Input Attribute, Resource, Resource Output Attribute, Module Output\n")
	attrByResource := make(map[string]map[string]map[string]string) // resource-type.resource-name: attribute_kind: resource-attribute-name: attribute_path
	// fmt.Print("Resources:\n")
	for resourceKey, resource := range config.ManagedResources {
		inputs := make(map[string]string)
		for name, item := range resource.Inputs {
			inputs[name] = fmt.Sprintf("%s.%s.%s", item.ResourceType, item.ResourceName, strings.Join(item.AttributePath, ".")) // resource-attribute: module-input-var
			// fmt.Printf("%s.%s: %s '%s' . %s\n", resourceName, name, item.ResourceType, item.ResourceName, strings.Join(item.AttributePath, "."))
		}
		attrByResource[resourceKey] = map[string]map[string]string{
			"inputs": inputs,
		}
	}
	// fmt.Print("Outputs:\n")
	for name, item := range config.Outputs {
		attrsForResource, ok := attrByResource[item.Value.ResourceType]
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

		// fmt.Printf("%s: %s '%s' . %s\n", name, item.Value.ResourceType, item.Value.ResourceName, strings.Join(item.Value.AttributePath, "."))
	}

	for resourceKey, r1 := range attrByResource {
		for resourceAttributeName, moduleInputName := range r1["inputs"] {
			fmt.Printf("%s, %s, %s, %s, %s\n", moduleInputName, resourceAttributeName, resourceKey, "", "")
		}
		for resourceAttributeName, moduleOutputName := range r1["outputs"] {
			fmt.Printf("%s, %s, %s, %s, %s\n", "", "", resourceKey, resourceAttributeName, moduleOutputName)
		}
	}

	// b, err := os.ReadFile("./modules/terraform-aws-vpc/main.tf")
	// if err != nil {
	// 	log.Fatalf(err.Error())
	// }

	// config := make(map[string]interface{})
	// if err := hclsimple.Decode(
	// 	"example.hcl", b,
	// 	nil, &config,
	// ); err != nil {
	// 	log.Fatalf("Failed to load configuration: %s", err)
	// }
	// fmt.Printf("Configuration is %v\n", config)

	// config, diag := hclwrite.ParseConfig(b, "", hcl.Pos{Line: 1, Column: 1})
	// if diag.HasErrors() {
	// 	log.Fatalf(diag.Error())
	// }
	// for _, b := range config.Body().Blocks() {
	// 	fmt.Printf("> %s block %s\n", b.Type(), b.Labels())
	// 	for _, a := range b.Body().Attributes() {
	// 		fmt.Printf("  >> attribute %s\n", a.Expr().BuildTokens(nil).Bytes())
	// 		// for _, v := range a.Expr().Variables() {
	// 		// 	fmt.Printf("  >> %v\n", v.BuildTokens(nil))
	// 		// }
	// 	}
	// }

	// _, err := json.Marshal(config)
	// if err != nil {
	// 	log.Fatalf(err.Error())
	// }
	// fmt.Print(string(jb))
}
