// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfconfig

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcl/v2"
)

const AttributePathSeparator = "."

type Metadata struct {
	Name    string
	Source  string
	Version string
}

// Module is the top-level type representing a parsed and processed Terraform
// module.
type Module struct {
	// Path is the local filesystem directory where the module was loaded from.
	Path string `json:"path"`

	Metadata *Metadata

	Locals map[string]hcl.Expression

	Variables map[string]*Variable `json:"variables"`
	Outputs   map[string]*Output   `json:"outputs"`

	// All resolved input variables of this module mapped to underlying resource inputs (as <var-name>: <var-attr-path>: <res-input-refs>)
	// The inputs are resolved recursively (e.g. input variable of A -> module A -> module-call to B -> module B -> resource).
	// resource "res" {
	//	input-name = var.module_input_variable.variable_attribute
	// }
	// resolved as module_input_variable: variable_attribute: [res.input-name] (list since one variable may get passed to multiple resources)
	Inputs map[string]map[string][]ResourceAttributeReference `json:"inputs"`

	RequiredCore      []string                        `json:"required_core,omitempty"`
	RequiredProviders map[string]*ProviderRequirement `json:"required_providers"`

	ProviderConfigs  map[string]*ProviderConfig `json:"provider_configs,omitempty"`
	ManagedResources map[string]*Resource       `json:"managed_resources"`
	DataResources    map[string]*Resource       `json:"data_resources"`
	ModuleCalls      map[string]*ModuleCall     `json:"module_calls"`

	// Diagnostics records any errors and warnings that were detected during
	// loading, primarily for inclusion in serialized forms of the module
	// since this slice is also returned as a second argument from LoadModule.
	Diagnostics Diagnostics `json:"diagnostics,omitempty"`
}

// AddResourceReference maps a given dstResource attribute from this module to
// another srcResource attribute within this module or its sub-modules.
//
//	resource "destination_resource" {
//		attr = source_resource.attribute.path
//	}
func (m *Module) AddResourceReference(srcResource AttributeReference, dstResource ResourceAttributeReference) bool {
	return m.addResourceReference(m.ManagedResources, srcResource, dstResource) || m.addResourceReference(m.DataResources, srcResource, dstResource)
}

func (m *Module) addResourceReference(resources map[string]*Resource, reference AttributeReference, referencedBy ResourceAttributeReference) bool {
	for _, item := range resources {
		if item.Name == referencedBy.ResourceName && item.Type == referencedBy.ResourceType {
			attrPath := referencedBy.Path()
			current, ok := item.References[attrPath]
			if !ok {
				current = make([]AttributeReference, 0)
			}
			item.References[attrPath] = append(current, reference)
			return true
		}
	}
	return false
}

// GetResourceAttributeReferences recursively resolves a given variable name to underlying resource reference.
// It follows the variable to any sub-module calls and returns all resources it gets passed into.
func (m Module) GetResourceAttributeReferences(varName string) []RelativeAttributeReference {
	references := make([]RelativeAttributeReference, 0)

	for _, mc := range m.ModuleCalls {
		if mc.Module != nil {
			for inpName, inp := range mc.Inputs {
				if inp.ResourceType == "var" && inp.ResourceName == varName {
					found := mc.Module.GetResourceAttributeReferences(inpName)
					for i := range found {
						found[i].RelativePath = append(found[i].RelativePath, inp.AttributePath...)
					}
					references = append(references, found...)
				}
			}
		}
	}

	for _, res := range m.ManagedResources {
		for inpName, inp := range res.Inputs {
			if inp.ResourceType == "var" && inp.ResourceName == varName {
				references = append(references, RelativeAttributeReference{
					ResourceAttributeReference: ResourceAttributeReference{
						Module:        &m,
						ResourceType:  res.Type,
						ResourceName:  res.Name,
						AttributePath: []string{inpName},
					},
					RelativePath: inp.AttributePath,
				})
			}
		}
	}
	return references
}

type AttributeReference interface {
	fmt.Stringer
	Type() string
	Name() string
	Attribute() string
	Path() string
	Pos() (file string, line int)
}

type RelativeAttributeReference struct {
	ResourceAttributeReference
	RelativePath []string `json:"relative_path"` // appended to the base ResourceReference
}

func (r RelativeAttributeReference) Path() string {
	return buildAttributePath(append(append([]string{}, r.AttributePath...), r.RelativePath...)...)
}

type ResourceAttributeReference struct {
	Expression    hcl.Expression `json:"expression"`
	Module        *Module
	ResourceType  string   `json:"resource_type"`
	ResourceName  string   `json:"resource_name"`
	AttributePath []string `json:"attribute_path"`
}

func (r ResourceAttributeReference) Pos() (string, int) {
	return r.Expression.StartRange().Filename, r.Expression.StartRange().Start.Line
}

func (r ResourceAttributeReference) Type() string {
	return r.ResourceType
}

func (r ResourceAttributeReference) Name() string {
	return r.ResourceName
}

func (r ResourceAttributeReference) String() string {
	return buildAttributePath(r.ResourceType, r.Attribute())
}

func buildAttributePath(tokens ...string) string {
	return strings.Join(tokens, AttributePathSeparator)
}

func (r ResourceAttributeReference) Attribute() string {
	return buildAttributePath(r.ResourceType, r.Path())
}

func (r ResourceAttributeReference) Path() string {
	return buildAttributePath(r.AttributePath...)
}

func (r ResourceAttributeReference) MakeRelative(relativePath []string) RelativeAttributeReference {
	return RelativeAttributeReference{
		ResourceAttributeReference: r,
		RelativePath:               relativePath,
	}
}

// ProviderConfig represents a provider block in the configuration
type ProviderConfig struct {
	Name  string `json:"name"`
	Alias string `json:"alias,omitempty"`
}

// NewModule creates new Module representing Terraform module at the given path
func NewModule(path string) *Module {
	return &Module{
		Path:              path,
		Locals:            make(map[string]hcl.Expression),
		Variables:         make(map[string]*Variable),
		Outputs:           make(map[string]*Output),
		Inputs:            make(map[string]map[string][]ResourceAttributeReference),
		RequiredProviders: make(map[string]*ProviderRequirement),
		ProviderConfigs:   make(map[string]*ProviderConfig),
		ManagedResources:  make(map[string]*Resource),
		DataResources:     make(map[string]*Resource),
		ModuleCalls:       make(map[string]*ModuleCall),
	}
}
