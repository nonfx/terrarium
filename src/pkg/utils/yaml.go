// Copyright (c) Ollion
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"path/filepath"

	"golang.org/x/exp/slices"
)

// IsYaml checks if the file name or path contains a valid yaml extension
func IsYaml(filePath string) bool {
	return slices.Contains([]string{".yml", ".yaml"}, filepath.Ext(filePath))
}
