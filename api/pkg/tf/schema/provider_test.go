package schema

import (
	"reflect"
	"testing"
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
			name: "test with simple attribute",
			fields: fields{
				Attributes: map[string]AttributeRepresentation{
					"identity": {
						Type: map[string]interface{}{
							"id": "string",
						},
					},
				},
				BlockTypes: map[string]BlockTypeRepresentation{},
			},
			want: map[string]AttributeRepresentation{
				"identity.id": {
					Type: "string",
				},
			},
		},
		{
			name: "test with nested attribute types",
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
			if got := btr.ListLeafNodes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BlockRepresentation.ListLeafNodes() = %v, want %v", got, tt.want)
			}
		})
	}
}
