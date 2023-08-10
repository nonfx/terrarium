package constants

import "io/fs"

const (
	ModuleSchemaFilePath      = ".terraform/modules/modules.json"
	DefaultExampleTFDirectory = "examples/farm/modules"
)

// fs.FileMode constants
const (
	ReadWritePermissions fs.FileMode = 0644
)
