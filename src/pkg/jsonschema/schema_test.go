// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package jsonschema

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
