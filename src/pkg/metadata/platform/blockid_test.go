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

func TestParseComponent(t *testing.T) {
	tests := []struct {
		name     string
		blockID  BlockID
		wantName string
		wantType BlockType
	}{
		{
			name:     "is component module call",
			blockID:  "module.tr_component_test",
			wantName: "test",
			wantType: BlockType_ModuleCall,
		},
		{
			name:     "is not component",
			blockID:  "module.test",
			wantName: "",
			wantType: BlockType_ModuleCall,
		},
		{
			name:     "is component local input",
			blockID:  "local.tr_component_test",
			wantName: "test",
			wantType: BlockType_Local,
		},
		{
			name:     "is component invalid type",
			blockID:  "resource.tr_component_test",
			wantName: "",
			wantType: BlockType_Resource,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bt, compType := tt.blockID.ParseComponent()
			assert.Equal(t, tt.wantName, compType)
			assert.Equal(t, tt.wantType, bt)
		})
	}
}
