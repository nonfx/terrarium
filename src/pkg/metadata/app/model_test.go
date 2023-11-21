// Copyright (c) Ollion
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"testing"

	"github.com/cldcvr/terrarium/src/pkg/pb/terrariumpb"
	"github.com/stretchr/testify/assert"
)

func TestApp_Transformations(t *testing.T) {
	app := App{
		ID:        "af05842a-8fa2-4544-88a8-bddaf6f58daa",
		Name:      "Mission",
		EnvPrefix: "Division",
		Compute: Dependency{
			ID:        "7063b3bb-1326-45c1-9314-eaa967f7f8ee",
			Use:       "utilize",
			EnvPrefix: "Chips",
			Inputs: map[string]interface{}{
				"payment": 70795.25,
			},
			Outputs: map[string]string{
				"Austria": "Table",
			},
			NoProvision: true,
		},
		Dependencies: []Dependency{
			{
				ID:        "9b205eb2-394c-4a27-a5ce-adc2e17bb2d3",
				Use:       "Practical",
				EnvPrefix: "Global",
				Inputs: map[string]interface{}{
					"Internal": 60578.24,
				},
				Outputs: map[string]string{
					"utilisation": "circuit",
				},
				NoProvision: false,
			},
		},
	}

	anyMsg, err := app.WrapProtoMessage()
	assert.NoError(t, err)

	trMsg := &terrariumpb.App{}
	err = anyMsg.UnmarshalTo(trMsg)
	assert.NoError(t, err)

	other := App{}
	other.ScanProto(trMsg)

	assert.Equal(t, app, other)

	b, err := app.ToFileBytes()
	assert.NoError(t, err)

	other2, err := NewApp(b)
	assert.NoError(t, err)

	assert.Equal(t, app, *other2)

	v, err := app.Value()
	assert.NoError(t, err)

	other3 := App{}
	err = other3.Scan(v)
	assert.NoError(t, err)

	assert.Equal(t, app, other3)
}

func TestDependency_IsEquivalent(t *testing.T) {
	base := Dependency{
		ID:        "eda5d4e1-98bd-438b-b477-c05d32cf79ea",
		Use:       "Security",
		EnvPrefix: "strategic",
		Inputs: map[string]interface{}{
			"Baby": "Automated",
		},
		Outputs: map[string]string{
			"withdrawal": "compressing",
		},
		NoProvision: true,
	}
	type args struct {
		other Dependency
	}
	tests := []struct {
		name string
		base Dependency
		args args
		want bool
	}{
		{
			name: "same",
			base: base,
			args: args{
				other: Dependency{
					ID:        "eda5d4e1-98bd-438b-b477-c05d32cf79ea",
					Use:       "Security",
					EnvPrefix: "strategic",
					Inputs: map[string]interface{}{
						"Baby": "Automated",
					},
					Outputs: map[string]string{
						"withdrawal": "compressing",
					},
					NoProvision: true,
				},
			},
			want: true,
		},
		{
			name: "different id",
			base: base,
			args: args{
				other: Dependency{
					ID:        "diff",
					Use:       "Security",
					EnvPrefix: "strategic",
					Inputs: map[string]interface{}{
						"Baby": "Automated",
					},
					Outputs: map[string]string{
						"withdrawal": "compressing",
					},
					NoProvision: true,
				},
			},
			want: false,
		},
		{
			name: "different use",
			base: base,
			args: args{
				other: Dependency{
					ID:        "eda5d4e1-98bd-438b-b477-c05d32cf79ea",
					Use:       "diff",
					EnvPrefix: "strategic",
					Inputs: map[string]interface{}{
						"Baby": "Automated",
					},
					Outputs: map[string]string{
						"withdrawal": "compressing",
					},
					NoProvision: true,
				},
			},
			want: false,
		},
		{
			name: "different env",
			base: base,
			args: args{
				other: Dependency{
					ID:        "eda5d4e1-98bd-438b-b477-c05d32cf79ea",
					Use:       "Security",
					EnvPrefix: "diff",
					Inputs: map[string]interface{}{
						"Baby": "Automated",
					},
					Outputs: map[string]string{
						"withdrawal": "compressing",
					},
					NoProvision: true,
				},
			},
			want: true,
		},
		{
			name: "different input",
			base: base,
			args: args{
				other: Dependency{
					ID:        "eda5d4e1-98bd-438b-b477-c05d32cf79ea",
					Use:       "Security",
					EnvPrefix: "strategic",
					Inputs: map[string]interface{}{
						"Baby": "diff",
					},
					Outputs: map[string]string{
						"withdrawal": "compressing",
					},
					NoProvision: true,
				},
			},
			want: false,
		},
		{
			name: "different output",
			base: base,
			args: args{
				other: Dependency{
					ID:        "eda5d4e1-98bd-438b-b477-c05d32cf79ea",
					Use:       "Security",
					EnvPrefix: "strategic",
					Inputs: map[string]interface{}{
						"Baby": "Automated",
					},
					Outputs: map[string]string{
						"withdrawal": "diff",
					},
					NoProvision: true,
				},
			},
			want: true,
		},
		{
			name: "different provision",
			base: base,
			args: args{
				other: Dependency{
					ID:        "eda5d4e1-98bd-438b-b477-c05d32cf79ea",
					Use:       "Security",
					EnvPrefix: "strategic",
					Inputs: map[string]interface{}{
						"Baby": "Automated",
					},
					Outputs: map[string]string{
						"withdrawal": "compressing",
					},
					NoProvision: false,
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.base.IsEquivalent(tt.args.other); got != tt.want {
				t.Errorf("Dependency.IsEquivalent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApp_IsEquivalent(t *testing.T) {
	mockDeps := []Dependency{
		{
			ID:        "23ea0e71-4951-45b2-b302-67d737cfcde9",
			Use:       "Chair",
			EnvPrefix: "Congolese",
			Inputs: map[string]interface{}{
				"Account": "Berkshire",
			},
			Outputs: map[string]string{
				"world-class": "calculating",
			},
			NoProvision: true,
		},
		{
			ID:        "5d6b206c-5c41-4e74-953a-946309897a4b",
			Use:       "Fresh",
			EnvPrefix: "schemas",
			Inputs: map[string]interface{}{
				"Fresh": "Mexico",
			},
			Outputs: map[string]string{
				"throughput": "hack",
			},
			NoProvision: true,
		},
		{
			ID:        "9bf332fe-41c7-4219-8b56-ccc66a7e4578",
			Use:       "back-end",
			EnvPrefix: "Director",
			Inputs: map[string]interface{}{
				"Car": "copy",
			},
			Outputs: map[string]string{
				"e-business": "Iraqi",
			},
			NoProvision: true,
		},
	}
	type args struct {
		other App
	}
	tests := []struct {
		name string
		base App
		args args
		want bool
	}{
		{
			name: "same",
			base: App{
				ID:        "08e250ac-09ba-4509-9278-facc31142a2d",
				Name:      "Synergistic",
				EnvPrefix: "efficient",
				Compute:   mockDeps[0],
				Dependencies: []Dependency{
					mockDeps[1],
					mockDeps[2],
				},
			},
			args: args{
				other: App{
					ID:        "a9d2beac-c5f3-489d-a871-a69292c78a49",
					Name:      "user",
					EnvPrefix: "SDD",
					Compute:   mockDeps[0],
					Dependencies: []Dependency{
						mockDeps[1],
						mockDeps[2],
					},
				},
			},
			want: true,
		},
		{
			name: "same (no dependencies)",
			base: App{
				ID:        "08e250ac-09ba-4509-9278-facc31142a2d",
				Name:      "Synergistic",
				EnvPrefix: "efficient",
			},
			args: args{
				other: App{
					ID:        "a9d2beac-c5f3-489d-a871-a69292c78a49",
					Name:      "user",
					EnvPrefix: "SDD",
				},
			},
			want: true,
		},
		{
			name: "different compute",
			base: App{
				ID:        "08e250ac-09ba-4509-9278-facc31142a2d",
				Name:      "Synergistic",
				EnvPrefix: "efficient",
				Compute:   mockDeps[0],
				Dependencies: []Dependency{
					mockDeps[1],
					mockDeps[2],
				},
			},
			args: args{
				other: App{
					ID:        "a9d2beac-c5f3-489d-a871-a69292c78a49",
					Name:      "user",
					EnvPrefix: "SDD",
					Compute:   mockDeps[1],
					Dependencies: []Dependency{
						mockDeps[1],
						mockDeps[2],
					},
				},
			},
			want: false,
		},
		{
			name: "different dependencies",
			base: App{
				ID:        "08e250ac-09ba-4509-9278-facc31142a2d",
				Name:      "Synergistic",
				EnvPrefix: "efficient",
				Compute:   mockDeps[0],
				Dependencies: []Dependency{
					mockDeps[1],
					mockDeps[2],
				},
			},
			args: args{
				other: App{
					ID:        "a9d2beac-c5f3-489d-a871-a69292c78a49",
					Name:      "user",
					EnvPrefix: "SDD",
					Compute:   mockDeps[0],
					Dependencies: []Dependency{
						mockDeps[0],
						mockDeps[2],
					},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotIsEquivalent := tt.base.IsEquivalent(tt.args.other); gotIsEquivalent != tt.want {
				t.Errorf("App.IsEquivalent() = %v, want %v", gotIsEquivalent, tt.want)
			}
		})
	}
}
