// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package platform

import (
	"os"
	"testing"

	"github.com/cldcvr/terraform-config-inspect/tfconfig"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestParse(t *testing.T) {
	dir := "../../../../examples/platform"
	m, diags := tfconfig.LoadModule(dir, &tfconfig.ResolvedModulesSchema{})
	if diags.HasErrors() {
		panic(diags.Err())
	}

	expected, err := os.ReadFile("./test-data.yaml")
	if err != nil {
		panic(err)
	}

	pm, _ := NewPlatformMetadata(m, nil)

	actual, err := yaml.Marshal(pm)
	if err != nil {
		panic(err)
	}

	// update test data when code changes
	// os.WriteFile("./test-data.yaml", actual, 0644)

	assert.YAMLEq(t, string(expected), string(actual))
}
