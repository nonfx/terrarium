package platform

import (
	"github.com/cldcvr/terraform-config-inspect/tfconfig"
	"github.com/cldcvr/terrarium/src/pkg/jsonschema"
	"gopkg.in/yaml.v3"
)

const (
	ComponentPrefix = "tr_component_"
)

type Component struct {
	ID          string           `yaml:",omitempty"`
	Taxonomy    []string         `yaml:",omitempty"`
	Title       string           `yaml:",omitempty"`
	Description string           `yaml:",omitempty"`
	Inputs      *jsonschema.Node `yaml:",omitempty"`
	Outputs     *jsonschema.Node `yaml:",omitempty"`
}

type Components []Component

type PlatformMetadata struct {
	Components Components
	Graph      Graph
}

type BlockID string

type GraphNode struct {
	ID           BlockID
	Requirements []BlockID
}

type Graph []GraphNode

type BlockType string

const (
	BlockType_Undefined  BlockType = ""
	BlockType_ModuleCall BlockType = "module"
	BlockType_Resource   BlockType = "resource"
	BlockType_Data       BlockType = "data"
	BlockType_Local      BlockType = "local"
	BlockType_Variable   BlockType = "var"
	BlockType_Output     BlockType = "output"
	BlockType_Provider   BlockType = "provider"
)

func NewPlatformMetadata(platformModule *tfconfig.Module, existingYaml []byte) (*PlatformMetadata, error) {
	p := PlatformMetadata{Components: Components{}, Graph: Graph{}}

	if len(existingYaml) > 0 {
		err := yaml.Unmarshal(existingYaml, &p)
		if err != nil {
			return nil, err
		}
	}

	p.Components.Parse(platformModule)
	p.Graph.Parse(platformModule)
	return &p, nil
}
