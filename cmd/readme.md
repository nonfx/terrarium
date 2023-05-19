# Commands Module

This Go module contains various entry points for Go commands.

## Usage

To understand the usage and functionality of each command, refer to the dedicated `readme.md` file in the command directory.

## Setup

This module relies on dependencies from other modules within the Go workspace. Therefore, it cannot be imported separately into other repositories.

To set up this module properly, follow these steps:

1. Run `go work sync` in the root directory to update `go.mod` and `go.sync`.
2. Avoid running `go mod tidy` within this module, as it is not aware of the workspace and will fail to recognize other modules present. For more information, refer to [golang/go/issues#50750](https://github.com/golang/go/issues/50750).
