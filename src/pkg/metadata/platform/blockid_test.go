package platform

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBlockID(t *testing.T) {
	tests := []struct {
		name      string
		blockType BlockType
		blockName string
		expected  BlockID
	}{
		{
			name:      "valid block type and name",
			blockType: BlockType_ModuleCall,
			blockName: "test",
			expected:  "module.test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewBlockID(tt.blockType, tt.blockName)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetBlockType(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected BlockType
	}{
		{
			name:     "valid block type",
			input:    "module",
			expected: BlockType_ModuleCall,
		},
		{
			name:     "invalid block type",
			input:    "invalid",
			expected: BlockType_Undefined,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetBlockType(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBlockIDParse(t *testing.T) {
	tests := []struct {
		name         string
		blockID      BlockID
		expectedType BlockType
		expectedName string
	}{
		{
			name:         "valid block ID",
			blockID:      "module.test",
			expectedType: BlockType_ModuleCall,
			expectedName: "test",
		},
		{
			name:         "invalid block ID",
			blockID:      "invalid",
			expectedType: BlockType_Undefined,
			expectedName: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultType, resultName := tt.blockID.Parse()
			assert.Equal(t, tt.expectedType, resultType)
			assert.Equal(t, tt.expectedName, resultName)
		})
	}
}

func TestIsComponent(t *testing.T) {
	tests := []struct {
		name     string
		blockID  BlockID
		expected bool
	}{
		{
			name:     "is component",
			blockID:  "module.tr_component_test",
			expected: true,
		},
		{
			name:     "is not component",
			blockID:  "module.test",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.blockID.IsComponent()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetComponentName(t *testing.T) {
	tests := []struct {
		name     string
		blockID  BlockID
		expected string
	}{
		{
			name:     "valid component name",
			blockID:  "module.tr_component_test",
			expected: "test",
		},
		{
			name:     "invalid component name",
			blockID:  "module.test",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.blockID.GetComponentName()
			assert.Equal(t, tt.expected, result)
		})
	}
}
