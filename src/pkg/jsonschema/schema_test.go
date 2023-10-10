// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package jsonschema

import (
	"errors"
	"testing"

	"github.com/cldcvr/terrarium/src/pkg/pb/terrariumpb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xeipuuv/gojsonschema"
)

func TestNode_Compile(t *testing.T) {
	tests := []struct {
		name    string
		node    *Node
		wantErr bool
	}{
		{
			name:    "valid node",
			node:    &Node{Type: "string"},
			wantErr: false,
		},
		{
			name:    "invalid node",
			node:    &Node{Type: "invalid"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.node.Compile()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, tt.node.compiled)
			}
		})
	}
}

func TestNode_Validate(t *testing.T) {
	strSchema := &Node{Type: "string"}
	tests := []struct {
		name    string
		node    *Node
		val     interface{}
		wantErr bool
	}{
		{
			name:    "invalid node",
			node:    &Node{Type: "invalid"},
			wantErr: true,
		},
		{
			name:    "invalid value",
			node:    &Node{Type: gojsonschema.TYPE_OBJECT},
			val:     map[bool]struct{}{true: struct{}{}},
			wantErr: true,
		},
		{
			name:    "valid string value",
			node:    strSchema,
			val:     "test",
			wantErr: false,
		},
		{
			name:    "invalid integer value",
			node:    strSchema,
			val:     123,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.node.Validate(tt.val)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNode_ApplyDefaultsToMSI(t *testing.T) {
	tests := []struct {
		name    string
		node    *Node
		inp     map[string]interface{}
		wantOut map[string]interface{}
	}{
		{
			name: "nil input",
			node: &Node{Properties: map[string]*Node{}},
			inp:  nil,
		},
		{
			name: "apply defaults to map",
			node: &Node{
				Properties: map[string]*Node{
					"key1": {Default: "defaultVal"},
					"key2": {Default: "defaultVal"},
				},
			},
			inp: map[string]interface{}{
				"key2": "value2",
			},
			wantOut: map[string]interface{}{
				"key1": "defaultVal",
				"key2": "value2",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.node.ApplyDefaultsToMSI(tt.inp)
			assert.Equal(t, tt.wantOut, tt.inp)
		})
	}
}

func TestNode_ApplyDefaultsToArr(t *testing.T) {
	tests := []struct {
		name    string
		node    *Node
		inp     []interface{}
		wantOut []interface{}
	}{
		{
			name: "nil input",
			node: &Node{Properties: map[string]*Node{}},
			inp:  nil,
		},
		{
			name: "apply defaults to array",
			node: &Node{
				Items: &Node{Default: "defaultItem"},
			},
			inp:     []interface{}{nil, "item2"},
			wantOut: []interface{}{"defaultItem", "item2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.node.ApplyDefaultsToArr(tt.inp)
			assert.Equal(t, tt.wantOut, tt.inp)
		})
	}
}

func TestNodeScan(t *testing.T) {
	n := &Node{}

	// Successful unmarshalling with a string
	validJSONStr := `{"type":"testType"}`
	err := n.Scan(validJSONStr)
	assert.Nil(t, err)
	assert.Equal(t, "testType", n.Type)

	// Successful unmarshalling with []byte
	validJSONBytes := []byte(`{"type":"anotherType"}`)
	err = n.Scan(validJSONBytes)
	assert.Nil(t, err)
	assert.Equal(t, "anotherType", n.Type)

	// Failure due to incorrect type assertion
	err = n.Scan(123)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "type assertion to string or []byte failed")

	// Failure due to invalid JSON
	invalidJSON := `{"type":}`
	err = n.Scan(invalidJSON)
	assert.NotNil(t, err)
}

func TestNodeValue(t *testing.T) {
	tests := []struct {
		name string
		node Node
		want error
	}{
		{
			name: "successful value",
			node: Node{
				Type: "string",
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.node.Value()
			if tt.want == nil && err != nil {
				t.Fatalf("expected no error, got: %v", err)
			}
			if tt.want != nil && !errors.Is(err, tt.want) {
				t.Fatalf("expected error %v, got: %v", tt.want, err)
			}
		})
	}
}

func TestNode_ToProto(t *testing.T) {
	tests := []struct {
		name      string
		jsn       *Node
		validator func(*testing.T, *terrariumpb.JSONSchema)
		wantErr   bool
	}{
		{
			name: "valid JSONSchema",
			jsn: &Node{
				Title:       "Test Title",
				Description: "Test Description",
				Type:        "object",
				Properties: map[string]*Node{
					"property1": {
						Title: "Property 1",
						Type:  "string",
					},
					"property2": {
						Title: "Property 2",
						Type:  "integer",
					},
				},
			},
			validator: func(t *testing.T, result *terrariumpb.JSONSchema) {
				require.NotNil(t, result)
				assert.Equal(t, "Test Title", result.Title)
				assert.Equal(t, "Test Description", result.Description)
				assert.Equal(t, "object", result.Type)
				require.Len(t, result.Properties, 2)
				assert.NotNil(t, result.Properties["property1"])
				assert.NotNil(t, result.Properties["property2"])
			},
			wantErr: false,
		},
		{
			name: "nil JSONSchema",
			jsn:  nil,
			validator: func(t *testing.T, result *terrariumpb.JSONSchema) {
				require.Nil(t, result)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.jsn.ToProto()
			tt.validator(t, result)
		})
	}
}
