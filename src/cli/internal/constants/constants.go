// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package constants

import "io/fs"

const (
	ModuleSchemaFilePath      = ".terraform/modules/modules.json"
	DefaultExampleTFDirectory = "examples/farm/modules"
)

// fs.FileMode constants
const (
	ReadWritePermissions        fs.FileMode = 0644
	ReadWriteExecutePermissions fs.FileMode = 0755
)
