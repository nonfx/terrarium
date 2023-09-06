// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBlockRepresentation_ListLeafNodes(t *testing.T) {
	type fields struct {
		Attributes map[string]AttributeRepresentation
		BlockTypes map[string]BlockTypeRepresentation
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]AttributeRepresentation
	}{
		{
			name: "NoNestedBlocks",
			fields: fields{
				Attributes: map[string]AttributeRepresentation{
					"attribute1": {Type: "string", Description: "Attr 1", Required: true},
					"attribute2": {Type: "int", Description: "Attr 2", Optional: true},
				},
				BlockTypes: map[string]BlockTypeRepresentation{},
			},
			want: map[string]AttributeRepresentation{
				"attribute1": {Type: "string", Description: "Attr 1", Required: true},
				"attribute2": {Type: "int", Description: "Attr 2", Optional: true},
			},
		},
		{
			name: "NestedBlocks",
			fields: fields{
				Attributes: map[string]AttributeRepresentation{
					"attribute1": {Type: "string", Description: "Attr 1", Required: true},
				},
				BlockTypes: map[string]BlockTypeRepresentation{
					"nestedBlock1": {
						NestingMode: "single",
						Block: BlockRepresentation{
							Attributes: map[string]AttributeRepresentation{
								"nestedAttribute1": {Type: "bool", Description: "Nested Attr 1", Optional: true},
							},
						},
					},
				},
			},
			want: map[string]AttributeRepresentation{
				"attribute1":                    {Type: "string", Description: "Attr 1", Required: true},
				"nestedBlock1":                  {Type: "single"},
				"nestedBlock1.nestedAttribute1": {Type: "bool", Description: "Nested Attr 1", Optional: true},
			},
		},
		{
			name: "NestedAttributeTypes",
			fields: fields{
				Attributes: map[string]AttributeRepresentation{
					"identity": {
						Type: []interface{}{
							"list", []interface{}{
								"object", map[string]interface{}{
									"oidc": []interface{}{
										"list", []interface{}{
											"object", map[string]interface{}{
												"issuer": "string",
											},
										},
									},
								},
							},
						},
					},
				},
				BlockTypes: map[string]BlockTypeRepresentation{},
			},
			want: map[string]AttributeRepresentation{
				"identity": {
					Type: "list",
				},
				"identity.oidc": {
					Type: "list",
				},
				"identity.oidc.issuer": {
					Type: "string",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			btr := BlockRepresentation{
				Attributes: tt.fields.Attributes,
				BlockTypes: tt.fields.BlockTypes,
			}

			got := btr.ListLeafNodes()
			assert.Equal(t, tt.want, got)
		})
	}
}
