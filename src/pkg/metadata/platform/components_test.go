// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package platform

import (
	"testing"

	"github.com/cldcvr/terraform-config-inspect/tfconfig"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/stretchr/testify/assert"
	"github.com/zclconf/go-cty/cty"
)

// Helper function to create a test platformModule for testing.
func createTestPlatformModule() *tfconfig.Module {
	return &tfconfig.Module{
		ModuleCalls: map[string]*tfconfig.ModuleCall{
			"tr_component_test": {},
		},
		Locals: map[string]*tfconfig.Local{
			"tr_component_test": &tfconfig.Local{
				Name: "tr_component_test",
				Expression: &hclsyntax.LiteralValueExpr{Val: cty.MapVal(map[string]cty.Value{
					"default": cty.ObjectVal(map[string]cty.Value{
						"input1": cty.StringVal("value1"),
						"input2": cty.NumberFloatVal(42.0),
						"input3": cty.BoolVal(true),
						"input4": cty.TupleVal([]cty.Value{
							cty.StringVal("list_value1"),
						}),
						"input5": cty.ObjectVal(map[string]cty.Value{
							"obj_key1": cty.StringVal("obj_value1"),
						}),
					}),
				})},
			},
			"random_local": nil,
		},
		Outputs: map[string]*tfconfig.Output{
			"tr_component_test_output1": {
				Description: "Output description 1",
			},
			"tr_component_test_output2": {
				Description: "Output description 2",
			},
			"output3": {},
		},
	}
}

func TestNewComponents(t *testing.T) {
	module := createTestPlatformModule()
	components := NewComponents(module)

	assert.Len(t, components, 1, "Number of components should be 1")
	assert.Equal(t, "test", components[0].ID, "Component ID should be 'test'")
	assert.Equal(t, "value1", components[0].Inputs.Properties["input1"].Default, "Default value for input1 should be 'value1'")
	assert.Equal(t, 42.0, components[0].Inputs.Properties["input2"].Default, "Default value for input2 should be 42")
	assert.Equal(t, true, components[0].Inputs.Properties["input3"].Default, "Default value for input3 should be true")
	assert.Equal(t, "Output description 1", components[0].Outputs.Properties["output1"].Description, "Description for output1 should match")
	assert.Equal(t, "Output description 2", components[0].Outputs.Properties["output2"].Description, "Description for output2 should match")
}

func TestComponents_GetByID(t *testing.T) {
	module := createTestPlatformModule()
	components := NewComponents(module)

	component := components.GetByID("test")
	assert.NotNil(t, component, "Component with ID 'test' should be found")

	component = components.GetByID("nonexistent")
	assert.Nil(t, component, "Non-existent component should return nil")
}

func TestComponents_Append(t *testing.T) {
	components := Components{}
	c := Component{ID: "test"}

	component := components.Append(c)

	assert.Len(t, components, 1, "Number of components should be 1 after appending")
	assert.Equal(t, &c, component, "Appended component should match the one returned")
}

func TestComponents_Parse(t *testing.T) {
	module := createTestPlatformModule()
	components := Components{}
	components.Parse(module)

	assert.Len(t, components, 1, "Number of components should be 1 after parsing")
	assert.Equal(t, "test", components[0].ID, "Component ID should be 'test'")
	assert.Equal(t, "value1", components[0].Inputs.Properties["input1"].Default, "Default value for input1 should be 'value1'")
	assert.Equal(t, 42.0, components[0].Inputs.Properties["input2"].Default, "Default value for input2 should be 42")
	assert.Equal(t, true, components[0].Inputs.Properties["input3"].Default, "Default value for input3 should be true")
	assert.Equal(t, "Output description 1", components[0].Outputs.Properties["output1"].Description, "Description for output1 should match")
	assert.Equal(t, "Output description 2", components[0].Outputs.Properties["output2"].Description, "Description for output2 should match")
}

func TestComponent_fetchInputs(t *testing.T) {
	module := createTestPlatformModule()
	components := Components{}
	components.Parse(module)

	component := components.GetByID("test")
	assert.NotNil(t, component, "Component with ID 'test' should be found")

	component.fetchInputs(module)
	assert.NotNil(t, component.Inputs, "Inputs should not be nil after fetching")
	assert.Equal(t, "value1", component.Inputs.Properties["input1"].Default, "Default value for input1 should be 'value1'")
	assert.Equal(t, 42.0, component.Inputs.Properties["input2"].Default, "Default value for input2 should be 42")
	assert.Equal(t, true, component.Inputs.Properties["input3"].Default, "Default value for input3 should be true")
	assert.Equal(t, "array", component.Inputs.Properties["input4"].Type, "Type value for input4 should be array")
}

func TestComponent_fetchOutputs(t *testing.T) {
	module := createTestPlatformModule()
	components := Components{}
	components.Parse(module)

	component := components.GetByID("test")
	assert.NotNil(t, component, "Component with ID 'test' should be found")

	component.fetchOutputs(module)
	assert.NotNil(t, component.Outputs, "Outputs should not be nil after fetching")
	assert.NotNil(t, component.Outputs.Properties["output1"], "output1 property should be present in Outputs")
	assert.NotNil(t, component.Outputs.Properties["output2"], "output2 property should be present in Outputs")
}
