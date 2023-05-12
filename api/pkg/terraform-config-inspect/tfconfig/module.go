// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfconfig

import (
	"strings"

	"github.com/hashicorp/hcl/v2"
)

// Module is the top-level type representing a parsed and processed Terraform
// module.
type Module struct {
	// Path is the local filesystem directory where the module was loaded from.
	Path string `json:"path"`

	Locals map[string]hcl.Expression

	Variables map[string]*Variable                               `json:"variables"`
	Outputs   map[string]*Output                                 `json:"outputs"`
	Inputs    map[string]map[string][]ResourceAttributeReference `json:"inputs"` // var-name: var-attr-path: res-refs

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

func (m Module) GetResourceAttributeReferences(varName string) []InputReference {
	references := make([]InputReference, 0)

	for _, mod := range m.ModuleCalls {
		if mod.Module != nil {
			for inpName, inp := range mod.Inputs {
				if inp.ResourceType == "var" && inp.ResourceName == varName {
					found := mod.Module.GetResourceAttributeReferences(inpName)
					for i := range found {
						found[i].InputPath = append(found[i].InputPath, inp.AttributePath...)
					}
					references = append(references, found...)
				}
			}
		}
	}

	for _, res := range m.ManagedResources {
		for inpName, inp := range res.Inputs {
			if inp.ResourceType == "var" && inp.ResourceName == varName {
				references = append(references, InputReference{
					InputPath: inp.AttributePath,
					ResourceReference: ResourceAttributeReference{
						ResourceType:  res.Type,
						ResourceName:  res.Name,
						AttributePath: []string{inpName},
					},
				})
			}
		}
	}
	return references
}

type InputReference struct {
	InputPath         []string
	ResourceReference ResourceAttributeReference
}

type ResourceAttributeReference struct {
	ResourceType  string   `json:"resource_type"`
	ResourceName  string   `json:"resource_name"`
	AttributePath []string `json:"attribute_path"`
}

func (r ResourceAttributeReference) ToKey() string {
	return strings.Join(append([]string{r.ResourceType, r.ResourceName}, r.AttributePath...), ".")
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
