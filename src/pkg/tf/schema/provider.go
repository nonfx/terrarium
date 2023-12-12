// Copyright (c) Ollion
// SPDX-License-Identifier: Apache-2.0

package schema

import "fmt"

// Reference:
//
// - https://developer.hashicorp.com/terraform/cli/commands/providers/schema
//

// ProvidersSchema is the schema representation for the Providers Schema Representation.
type ProvidersSchema struct {
	// FormatVersion is the version of the schema format.
	FormatVersion string `json:"format_version"`

	// ProviderSchemas describes the provider schemas for all providers throughout the configuration tree.
	// Keys in this map are the provider type, such as "random".
	ProviderSchemas map[string]ProviderSchema `json:"provider_schemas"`
}

// ProviderSchema describes the schema for a Terraform provider.
type ProviderSchema struct {
	// Provider is the schema for the provider configuration.
	Provider SchemaRepresentation `json:"provider"`

	// ResourceSchemas maps the resource type name to the resource's schema.
	ResourceSchemas map[string]SchemaRepresentation `json:"resource_schemas,omitempty"`

	// DataSourceSchemas maps the data source type name to the data source's schema.
	DataSourceSchemas map[string]SchemaRepresentation `json:"data_source_schemas,omitempty"`
}

// SchemaRepresentation represents a provider or resource schema paired with
// that schema's version.
type SchemaRepresentation struct {
	// Version is the schema version, not the provider version.
	Version int64 `json:"version"`

	// Block contains the block representation of the schema.
	Block BlockRepresentation `json:"block"`
}

// BlockRepresentation represents the schema definition for a Terraform provider block.
type BlockRepresentation struct {
	// Attributes describes any attributes that appear directly inside the block. Keys in this map are the attribute names.
	Attributes map[string]AttributeRepresentation `json:"attributes"`

	// BlockTypes describes any nested blocks that appear directly inside the block. Keys in this map are the names of the block_type.
	BlockTypes map[string]BlockTypeRepresentation `json:"block_types,omitempty"`
}

// AttributeRepresentation represents the schema definition for a Terraform provider attribute.
type AttributeRepresentation struct {
	// Type is a representation of a type specification that the attribute's value must conform to.
	// This can be string or array. Inconsistencies have been observed
	Type interface{} `json:"type,omitempty"`

	// Description is an English-language description of the purpose and usage of the attribute.
	Description string `json:"description,omitempty"`

	// Required, if set to true, specifies that an omitted or null value is not permitted.
	Required bool `json:"required,omitempty"`

	// Optional, if set to true, specifies that an omitted or null value is permitted.
	Optional bool `json:"optional,omitempty"`

	// Computed, if set to true, indicates that the value comes from the provider rather than the configuration.
	Computed bool `json:"computed,omitempty"`

	// Sensitive, if set to true, indicates that the attribute may contain sensitive information.
	Sensitive bool `json:"sensitive,omitempty"`
}

// BlockTypeRepresentation represents the schema definition for a Terraform provider nested block type.
type BlockTypeRepresentation struct {
	// NestingMode describes the nesting mode for the child block, and can be one of the following:
	//    single
	//    list
	//    set
	//    map
	NestingMode string `json:"nesting_mode,omitempty"`

	// Block is a BlockRepresentation that represents the nested block.
	Block BlockRepresentation `json:"block,omitempty"`

	// MinItems and MaxItems set lower and upper limits on the number of child blocks allowed for the list and set modes. These are omitted for other modes.
	MinItems int `json:"min_items,omitempty"`
	MaxItems int `json:"max_items,omitempty"`
}

// ListLeafNodes returns a map of leaf nodes in the BlockRepresentation.
// Leaf nodes are the attributes directly present in the block and attributes within nested blocks.
// The keys of the map are the attribute paths, and the values are the corresponding AttributeRepresentation.
//
// Example:
//
//	Input: { a: string, b: {c: string, d: string}}
//	Output: { a: string, b: map, "b.c": string, "b.d": string }
func (btr BlockRepresentation) ListLeafNodes() map[string]AttributeRepresentation {
	leafNodes := map[string]AttributeRepresentation{}

	// Iterate over the attributes directly present in the block
	for name, attr := range btr.Attributes {
		for k, v := range getAttributePaths(name, attr) {
			leafNodes[k] = v
		}
	}

	// Iterate over the nested blocks and recursively collect leaf nodes
	for path, nestedBlock := range btr.BlockTypes {
		leafNodes[path] = AttributeRepresentation{
			Type: nestedBlock.NestingMode,
		}
		nestedLeafNodes := nestedBlock.Block.ListLeafNodes()
		for name, attr := range nestedLeafNodes {
			// Append the nested block's attribute name to the path
			newPath := fmt.Sprintf("%s.%s", path, name)
			leafNodes[newPath] = attr
		}
	}

	return leafNodes
}

func getAttributePaths(attrName string, attr AttributeRepresentation) map[string]AttributeRepresentation {
	flattened := make(map[string]AttributeRepresentation)
	buildAttributePaths(attrName, attr, flattened)
	return flattened
}

func buildAttributePaths(attrName string, attr AttributeRepresentation, flattened map[string]AttributeRepresentation) {
	attrType := attr.Type
	switch t := attr.Type.(type) {
	case map[string]interface{}:
		attrType = "object"
		for k, v := range t {
			buildAttributePaths(fmt.Sprintf("%s.%s", attrName, k), copyWithType(attr, v), flattened)
		}
	case []interface{}:
		attrType = t[0]
		buildAttributePaths(attrName, copyWithType(attr, t[1]), flattened)
	}

	flattened[attrName] = copyWithType(attr, attrType)
}

func copyWithType(src AttributeRepresentation, ty interface{}) AttributeRepresentation {
	return AttributeRepresentation{
		Type:        ty,
		Description: src.Description,
		Required:    src.Required,
		Optional:    src.Optional,
		Computed:    src.Computed,
		Sensitive:   src.Sensitive,
	}
}
