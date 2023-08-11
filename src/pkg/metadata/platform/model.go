package platform

import (
	"github.com/cldcvr/terraform-config-inspect/tfconfig"
	"github.com/cldcvr/terrarium/src/pkg/jsonschema"
	"gopkg.in/yaml.v3"
)

const (
	ComponentPrefix = "tr_component_" // Prefix for component identifiers in terraform code
)

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
	return &p, nil
}
