package platform

import (
	"testing"

	"github.com/cldcvr/terraform-config-inspect/tfconfig"
	"github.com/hashicorp/hcl/v2"
	"github.com/stretchr/testify/assert"
)

func TestGraph_GetByID(t *testing.T) {
	// Create a sample graph
	graph := Graph{
		{ID: NewBlockID(BlockType_ModuleCall, "module1")},
		{ID: NewBlockID(BlockType_Resource, "resource1")},
		{ID: NewBlockID(BlockType_Output, "output1")},
	}

	// Test GetByID with existing IDs
	assert.Equal(t, &graph[0], graph.GetByID(NewBlockID(BlockType_ModuleCall, "module1")))
	assert.Equal(t, &graph[1], graph.GetByID(NewBlockID(BlockType_Resource, "resource1")))
	assert.Equal(t, &graph[2], graph.GetByID(NewBlockID(BlockType_Output, "output1")))

	// Test GetByID with non-existing ID
	assert.Nil(t, graph.GetByID(NewBlockID(BlockType_ModuleCall, "non_existing_module")))
}

func TestGraph_Append(t *testing.T) {
	graph := Graph{}

	// Test Append with a single node
	node := graph.Append(NewBlockID(BlockType_ModuleCall, "module1"), nil)
	assert.Equal(t, 1, len(graph))
	assert.Equal(t, node, &graph[0])
	assert.Equal(t, NewBlockID(BlockType_ModuleCall, "module1"), graph[0].ID)
	assert.Nil(t, graph[0].Requirements)

	// Test Append with multiple nodes
	node2 := graph.Append(NewBlockID(BlockType_Resource, "resource1"), []BlockID{NewBlockID(BlockType_ModuleCall, "module1")})
	assert.Equal(t, 2, len(graph))
	assert.Equal(t, node2, &graph[1])
	assert.Equal(t, NewBlockID(BlockType_Resource, "resource1"), graph[1].ID)
	assert.Equal(t, []BlockID{NewBlockID(BlockType_ModuleCall, "module1")}, graph[1].Requirements)
}

func TestGraph_Parse(t *testing.T) {
	// Create a sample tfconfig.Module
	module := &tfconfig.Module{
		ModuleCalls: map[string]*tfconfig.ModuleCall{
			"tr_component_module1": {},
		},
		Outputs: map[string]*tfconfig.Output{
			"output1": {},
		},
	}

	// Test Parse with a simple module and output
	graph := NewGraph(module)
	assert.Equal(t, 2, len(graph))
	assert.Equal(t, NewBlockID(BlockType_ModuleCall, "tr_component_module1"), graph[0].ID)
	assert.Equal(t, NewBlockID(BlockType_Output, "output1"), graph[1].ID)
}

func TestParseWithNestedModules(t *testing.T) {
	// Create a sample tfconfig.Module with nested modules
	module := &tfconfig.Module{
		ModuleCalls: map[string]*tfconfig.ModuleCall{
			"tr_component_module1": {
				Inputs: map[string]tfconfig.AttributeReference{
					"data_ref":    tfconfig.ResourceAttributeReference{ResourceType: "resource_type", ResourceName: "label2"},
					"module_ref":  tfconfig.ResourceAttributeReference{ResourceType: "module", ResourceName: "module2"},
					"local_ref":   tfconfig.ResourceAttributeReference{ResourceType: "local", ResourceName: "local1"},
					"var_ref":     tfconfig.ResourceAttributeReference{ResourceType: "var", ResourceName: "var1"},
					"unknown_ref": tfconfig.ResourceAttributeReference{ResourceType: "unknown_type", ResourceName: "unknown_name"},
				},
			},
			"module2": {
				Inputs: map[string]tfconfig.AttributeReference{
					"resource_ref": tfconfig.ResourceAttributeReference{ResourceType: "resource_type", ResourceName: "label1"},
				},
			},
		},
		ManagedResources: map[string]*tfconfig.Resource{
			"resource_type.label1": {Mode: tfconfig.ManagedResourceMode, Type: "resource_type", Name: "label1"},
		},
		DataResources: map[string]*tfconfig.Resource{
			"data.resource_type.label2": {Mode: tfconfig.DataResourceMode, Type: "resource_type", Name: "label2"},
		},
		Locals: map[string]hcl.Expression{
			"local1": nil,
		},
		Variables: map[string]*tfconfig.Variable{
			"var1": {},
		},
		Outputs: map[string]*tfconfig.Output{
			"output1": {},
		},
	}

	// Test Parse with nested modules
	graph := NewGraph(module)
	assert.Equal(t, Graph{
		{"data.data.resource_type.label2", []BlockID{}},
		{"local.local1", []BlockID{}},
		{"module.module2", []BlockID{"resource.resource_type.label1"}},
		{"module.tr_component_module1", []BlockID{"data.data.resource_type.label2", "local.local1", "module.module2", "var.var1"}},
		{"output.output1", []BlockID{}},
		{"resource.resource_type.label1", []BlockID{}},
		{"var.var1", []BlockID{}},
	}, graph)
}
