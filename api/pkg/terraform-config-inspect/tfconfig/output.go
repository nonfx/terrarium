// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfconfig

// Output represents a single output from a Terraform module.
type Output struct {
	Name        string                     `json:"name"`
	Value       ResourceAttributeReference `json:"value"` // PETMAL: store the value expression
	Description string                     `json:"description,omitempty"`
	Sensitive   bool                       `json:"sensitive,omitempty"`
	Pos         SourcePos                  `json:"pos"`
}

func (v Output) GetName() string {
	return v.Name
}

func (v Output) GetDescription() string {
	return v.Description
}

func (v Output) IsRequired() bool {
	return false
}

func (v Output) IsComputed() bool {
	return true
}
