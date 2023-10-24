// Copyright (c) Ollion
// SPDX-License-Identifier: Apache-2.0

package platform

import (
	"github.com/cldcvr/terraform-config-inspect/tfconfig"
	"github.com/cldcvr/terrarium/src/pkg/jsonschema"
	"gopkg.in/yaml.v3"
)

const (
	ComponentDocPrefix = "component"     // Documentation comment prefix (i.e. "component[<name-override>]: <description>")
	ComponentPrefix    = "tr_component_" // Prefix for component identifiers in terraform code
)

// Profile represents a set of pre-set configuration variables that can be applied to generated Terraform code.
type Profile struct {
	ID          string `yaml:",omitempty"` // Unique identifier for the profile
	Title       string `yaml:",omitempty"` // Descriptive title for the profile
	Description string `yaml:",omitempty"` // Detailed description of the profile's properties
}

// Profiles is a slice of Profile objects.
type Profiles []Profile

// Component represents an implementation of a dependency in the Terrarium platform.
type Component struct {
	ID          string           `yaml:",omitempty"` // Unique identifier for the component
	Title       string           `yaml:",omitempty"` // Descriptive title for the component
	Description string           `yaml:",omitempty"` // Detailed description of the component's functionality
	Inputs      *jsonschema.Node `yaml:",omitempty"` // Input parameters required by the component
	Outputs     *jsonschema.Node `yaml:",omitempty"` // Output properties produced by the component
}

// Components is a slice of Component objects.
type Components []Component

// PlatformMetadata represents the metadata for the Terrarium platform.
// It includes the components and the graph that defines the relationships between terraform blocks.
type PlatformMetadata struct {
	Profiles   Profiles   // Configuration profiles in the platform
	Components Components // Components in the Terrarium platform
	Graph      Graph      // Graph defining the relationships between terraform blocks
}

// BlockID represents a unique identifier for a terraform block.
type BlockID string

// GraphNode represents a single node in the graph.
// It includes the node's ID and the IDs of the nodes it depends on.
type GraphNode struct {
	ID           BlockID   // Unique identifier for the graph node
	Requirements []BlockID // IDs of other graph nodes that the current node depends on
}

// Graph defines the relationships between terraform blocks.
type Graph []GraphNode

// BlockType represents the type of a terraform block.
type BlockType string

const (
	BlockType_Undefined  BlockType = ""         // Undefined block type
	BlockType_ModuleCall BlockType = "module"   // Module call block type
	BlockType_Resource   BlockType = "resource" // Resource block type
	BlockType_Data       BlockType = "data"     // Data block type
	BlockType_Local      BlockType = "local"    // Local block type
	BlockType_Variable   BlockType = "var"      // Variable block type
	BlockType_Output     BlockType = "output"   // Output block type
	BlockType_Provider   BlockType = "provider" // Provider block type
)

type ParsedBlock interface {
	GetPos() tfconfig.SourcePos
}

type BlockDependencyGetter interface {
	GetDependencies() map[string]tfconfig.AttributeReference
}

type BlockParentPosGetter interface {
	GetParentPos() *tfconfig.SourcePos
}

type BlockProviderGetter interface {
	GeProviderName() string
}

// NewPlatformMetadata creates a new PlatformMetadata object.
// It parses the platform module and existing YAML to create the components and the graph.
func NewPlatformMetadata(platformModule *tfconfig.Module, existingYaml []byte) (*PlatformMetadata, error) {
	p := PlatformMetadata{Components: Components{}, Graph: Graph{}}

	// If there is existing YAML, unmarshal it into the PlatformMetadata object
	if len(existingYaml) > 0 {
		err := yaml.Unmarshal(existingYaml, &p)
		if err != nil {
			return nil, err
		}
	}

	// Parse the platform module to create the components and the graph
	p.Components.Parse(platformModule)
	p.Graph.Parse(platformModule)
	p.Profiles.Parse(platformModule)
	return &p, nil
}
