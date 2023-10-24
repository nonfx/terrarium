// Copyright (c) Ollion
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"fmt"
	"sort"
	"testing"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/cldcvr/terrarium/src/pkg/jsonschema"
	"github.com/cldcvr/terrarium/src/pkg/metadata/app"
	"github.com/cldcvr/terrarium/src/pkg/metadata/platform"
	"github.com/hoisie/mustache"
	"github.com/stretchr/testify/assert"
)

func TestGetAppEnvTemplate(t *testing.T) {
	type args struct {
		pm  *platform.PlatformMetadata
		app app.App
	}
	tests := []struct {
		name string
		args args
		want EnvVars
	}{
		{
			name: "default output no prefix",
			args: args{
				app: app.App{
					Dependencies: app.Dependencies{
						{
							ID:      "mydep1",
							Use:     "comp1",
							Outputs: map[string]string{},
						},
					},
				},
				pm: &platform.PlatformMetadata{
					Components: platform.Components{
						{
							ID: "comp1",
							Outputs: &jsonschema.Node{
								Properties: map[string]*jsonschema.Node{
									"outp1": {},
									"outp2": {},
								},
							},
						},
					},
				},
			},
			want: EnvVars{
				{"OUTP1", `{{ tr_component_comp1_outp1.value.mydep1 }}`},
				{"OUTP2", `{{ tr_component_comp1_outp2.value.mydep1 }}`},
			},
		},
		{
			name: "templated output with prefix",
			args: args{
				app: app.App{
					EnvPrefix: "APP",
					Dependencies: app.Dependencies{
						{
							ID:      "mydep1",
							Use:     "comp1",
							Outputs: map[string]string{},
						},
						{
							ID:        "mydep2",
							Use:       "comp1",
							EnvPrefix: "MYDEP2",
							Outputs: map[string]string{
								"COMB": `combination of {{outp1}} and {{outp2}}`,
							},
						},
						{
							ID:      "mydep3",
							Use:     "comp2",
							Outputs: map[string]string{},
						},
					},
				},
				pm: &platform.PlatformMetadata{
					Components: platform.Components{
						{
							ID: "comp1",
							Outputs: &jsonschema.Node{
								Properties: map[string]*jsonschema.Node{
									"outp1": {},
									"outp2": {},
								},
							},
						},
					},
				},
			},
			want: EnvVars{
				{"APP_MYDEP2_COMB", `combination of {{ tr_component_comp1_outp1.value.mydep2 }} and {{ tr_component_comp1_outp2.value.mydep2 }}`},
				{"APP_OUTP1", `{{ tr_component_comp1_outp1.value.mydep1 }}`},
				{"APP_OUTP2", `{{ tr_component_comp1_outp2.value.mydep1 }}`},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetAppEnvTemplate(tt.args.pm, tt.args.app)
			sort.Sort(got)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestEnvVars(t *testing.T) {
	tests := []struct {
		name             string
		vars             EnvVars
		wantRender       string
		wantRenderQuoted string
	}{
		{
			vars: EnvVars{
				{"APP_OUTP1", `{{ tr_component_comp1_outp1.value.mydep1 }}`},
				{"APP_OUTP2", `{{ tr_component_comp1_outp2.value.mydep1 }}`},
			},
			wantRender: heredoc.Doc(`
			APP_OUTP1={{ tr_component_comp1_outp1.value.mydep1 }}
			APP_OUTP2={{ tr_component_comp1_outp2.value.mydep1 }}
			`),
			wantRenderQuoted: heredoc.Doc(`
			APP_OUTP1="{{ tr_component_comp1_outp1.value.mydep1 }}"
			APP_OUTP2="{{ tr_component_comp1_outp2.value.mydep1 }}"
			`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantRender, tt.vars.Render())
			assert.Equal(t, tt.wantRenderQuoted, tt.vars.RenderWithQuotes())
		})
	}
}

func TestTemplate(t *testing.T) {
	stateOut := map[string]interface{}{
		"tr_component_postgres_host": map[string]interface{}{
			"sensitive": false,
			"type": []interface{}{
				"object",
				map[string]interface{}{"ledgerdb": "string"},
			},
			"value": map[string]interface{}{
				"ledgerdb": "the value!",
			},
		},
	}

	tests := []struct {
		name                     string
		compID, compInp, depName string
		want                     string
	}{
		{
			name:    "string value",
			compID:  "postgres",
			compInp: "host",
			depName: "ledgerdb",
			want:    "the value!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			template := fmt.Sprintf(envValTemp, tt.compID, tt.compInp, tt.depName)
			got := mustache.Render(template, stateOut)
			assert.Equal(t, tt.want, got)
		})
	}
}
