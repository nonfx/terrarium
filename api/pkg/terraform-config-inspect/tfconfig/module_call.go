// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfconfig

// ModuleCall represents a "module" block within a module. That is, a
// declaration of a child module from inside its parent.
type ModuleCall struct {
	Name    string `json:"name"`
	Source  string `json:"source"`
	Version string `json:"version,omitempty"`

	Pos SourcePos `json:"pos"`

	// All resolved input values to this module call (as <module-input-name>: <value-ref>).
	// The <module-input-name> maps to an input variable in the called module.
	// module "mod" {
	//	input-name = value-ref (i.e. var.variable_ref, module.mod.out_ref, other_res.attr)
	// }
	Inputs map[string]ResourceAttributeReference `json:"inputs"`

	Module *Module `json:"-"`
}
