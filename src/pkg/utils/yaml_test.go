// Copyright (c) Ollion
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsYaml(t *testing.T) {
	tests := []struct {
		name string
		inp  string
		outp bool
	}{
		{
			name: "empty",
			inp:  "",
			outp: false,
		},
		{
			name: "non yaml file name",
			inp:  "mock-file.txt",
			outp: false,
		},
		{
			name: "non yaml file path",
			inp:  "./a/b/mock-file.txt",
			outp: false,
		},
		{
			name: "yaml file name",
			inp:  "mock-file.yaml",
			outp: true,
		},
		{
			name: "yaml file path",
			inp:  "./a/b/mock-file.yaml",
			outp: true,
		},
		{
			name: "yml file name",
			inp:  "mock-file.yml",
			outp: true,
		},
		{
			name: "yml file path",
			inp:  "./a/b/mock-file.yml",
			outp: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsYaml(tt.inp)
			assert.Equal(t, tt.outp, got)
		})
	}
}
