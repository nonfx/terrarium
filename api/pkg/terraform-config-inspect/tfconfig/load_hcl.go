// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfconfig

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/hcl/v2/hclsyntax"
	"golang.org/x/exp/slices"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/zclconf/go-cty/cty"
	ctyjson "github.com/zclconf/go-cty/cty/json"
)

const (
	passOneLoggerPrefix = "PASS-1"
	passTwoLoggerPrefix = "PASS-2"
)

func loadModule(fs FS, dir string, resolvedModuleRefs *ResolvedModulesSchema) (*Module, Diagnostics) {
	mod := NewModule(dir)
	primaryPaths, diags := dirFiles(fs, dir)

	parser := hclparse.NewParser()

	if meta := resolvedModuleRefs.Find(mod.Path); meta != nil {
		mod.Metadata = &Metadata{
			Name:    meta.GetNormalizedKey(),
			Source:  meta.Source,
			Version: meta.Version,
		}
	}

	for _, filename := range primaryPaths {
		var file *hcl.File
		var fileDiags hcl.Diagnostics

		b, err := fs.ReadFile(filename)
		if err != nil {
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Failed to read file",
				Detail:   fmt.Sprintf("The configuration file %q could not be read.", filename),
			})
			continue
		}
		if strings.HasSuffix(filename, ".json") {
			file, fileDiags = parser.ParseJSON(b, filename)
		} else {
			file, fileDiags = parser.ParseHCL(b, filename)
		}
		diags = append(diags, fileDiags...)
		if file == nil {
			continue
		}

		contentDiags := LoadModuleFromFile(file, mod, resolvedModuleRefs)
		diags = append(diags, contentDiags...)
	}

	// second-pass: resolve remaining references
	// that could not be evaluated before all files are loaded
	runSecondPassResolution(mod, mod.ManagedResources)
	runSecondPassResolution(mod, mod.DataResources)

	return mod, diagnosticsHCL(diags)
}

func runSecondPassResolution(parentModule *Module, resources map[string]*Resource) {
	for i := range resources {
		for qualifiedAttrName, attrValue := range resources[i].Inputs {
			if slices.Contains([]string{"local", "module"}, attrValue.ResourceType) {
				resolveResourceInputReference(parentModule, resources[i], qualifiedAttrName, attrValue.Expression, passTwoLoggerPrefix)
			}
		}
	}
}

// LoadModuleFromFile reads given file, interprets it and stores in given Module
// This is useful for any caller which does tokenization/parsing on its own
// e.g. because it will reuse these parsed files later for more detailed
// interpretation.
func LoadModuleFromFile(file *hcl.File, mod *Module, resolvedModuleRefs *ResolvedModulesSchema) hcl.Diagnostics {
	var diags hcl.Diagnostics
	content, _, contentDiags := file.Body.PartialContent(rootSchema)
	diags = append(diags, contentDiags...)

	for _, block := range content.Blocks {
		switch block.Type {

		case "terraform":
			content, _, contentDiags := block.Body.PartialContent(terraformBlockSchema)
			diags = append(diags, contentDiags...)

			if attr, defined := content.Attributes["required_version"]; defined {
				var version string
				valDiags := gohcl.DecodeExpression(attr.Expr, nil, &version)
				diags = append(diags, valDiags...)
				if !valDiags.HasErrors() {
					mod.RequiredCore = append(mod.RequiredCore, version)
				}
			}

			for _, innerBlock := range content.Blocks {
				switch innerBlock.Type {
				case "required_providers":
					reqs, reqsDiags := decodeRequiredProvidersBlock(innerBlock)
					diags = append(diags, reqsDiags...)
					for name, req := range reqs {
						if _, exists := mod.RequiredProviders[name]; !exists {
							mod.RequiredProviders[name] = req
						} else {
							if req.Source != "" {
								source := mod.RequiredProviders[name].Source
								if source != "" && source != req.Source {
									diags = append(diags, &hcl.Diagnostic{
										Severity: hcl.DiagError,
										Summary:  "Multiple provider source attributes",
										Detail:   fmt.Sprintf("Found multiple source attributes for provider %s: %q, %q", name, source, req.Source),
										Subject:  &innerBlock.DefRange,
									})
								} else {
									mod.RequiredProviders[name].Source = req.Source
								}
							}

							mod.RequiredProviders[name].VersionConstraints = append(mod.RequiredProviders[name].VersionConstraints, req.VersionConstraints...)
						}
					}
				}
			}
		case "locals":
			content, contentDiags := block.Body.JustAttributes()
			diags = append(diags, contentDiags...)
			if !contentDiags.HasErrors() {
				for _, attr := range content {
					mod.Locals[attr.Name] = attr.Expr
				}
			}
		case "variable":
			content, _, contentDiags := block.Body.PartialContent(variableSchema)
			diags = append(diags, contentDiags...)

			name := block.Labels[0]
			v := &Variable{
				Name: name,
				Pos:  sourcePosHCL(block.DefRange),
			}

			mod.Variables[name] = v

			if attr, defined := content.Attributes["type"]; defined {
				// We handle this particular attribute in a somewhat-tricky way:
				// since Terraform may evolve its type expression syntax in
				// future versions, we don't want to be overly-strict in how
				// we handle it here, and so we'll instead just take the raw
				// source provided by the user, using the source location
				// information in the expression object.
				//
				// However, older versions of Terraform expected the type
				// to be a string containing a keyword, so we'll need to
				// handle that as a special case first for backward compatibility.

				var typeExpr string

				var typeExprAsStr string
				valDiags := gohcl.DecodeExpression(attr.Expr, nil, &typeExprAsStr)
				if !valDiags.HasErrors() {
					typeExpr = typeExprAsStr
				} else {
					rng := attr.Expr.Range()
					typeExpr = string(rng.SliceBytes(file.Bytes))
				}

				v.Type = typeExpr
			}

			if attr, defined := content.Attributes["description"]; defined {
				var description string
				valDiags := gohcl.DecodeExpression(attr.Expr, nil, &description)
				diags = append(diags, valDiags...)
				v.Description = description
			}

			if attr, defined := content.Attributes["default"]; defined {
				// To avoid the caller needing to deal with cty here, we'll
				// use its JSON encoding to convert into an
				// approximately-equivalent plain Go interface{} value
				// to return.
				val, valDiags := attr.Expr.Value(nil)
				diags = append(diags, valDiags...)
				if val.IsWhollyKnown() { // should only be false if there are errors in the input
					valJSON, err := ctyjson.Marshal(val, val.Type())
					if err != nil {
						// Should never happen, since all possible known
						// values have a JSON mapping.
						panic(fmt.Errorf("failed to serialize default value as JSON: %s", err))
					}
					var def interface{}
					err = json.Unmarshal(valJSON, &def)
					if err != nil {
						// Again should never happen, because valJSON is
						// guaranteed valid by ctyjson.Marshal.
						panic(fmt.Errorf("failed to re-parse default value from JSON: %s", err))
					}
					v.Default = def
				}
			} else {
				v.Required = true
			}

			if attr, defined := content.Attributes["sensitive"]; defined {
				var sensitive bool
				valDiags := gohcl.DecodeExpression(attr.Expr, nil, &sensitive)
				diags = append(diags, valDiags...)
				v.Sensitive = sensitive
			}

		case "output":

			content, _, contentDiags := block.Body.PartialContent(outputSchema)
			diags = append(diags, contentDiags...)

			name := block.Labels[0]
			o := &Output{
				Name: name,
				Pos:  sourcePosHCL(block.DefRange),
			}

			mod.Outputs[name] = o

			if attr, defined := content.Attributes["description"]; defined {
				var description string
				valDiags := gohcl.DecodeExpression(attr.Expr, nil, &description)
				diags = append(diags, valDiags...)
				o.Description = description
			}

			if attr, defined := content.Attributes["sensitive"]; defined {
				var sensitive bool
				valDiags := gohcl.DecodeExpression(attr.Expr, nil, &sensitive)
				diags = append(diags, valDiags...)
				o.Sensitive = sensitive
			}

			// PETMAL: extract the output value expression
			if attr, defined := content.Attributes["value"]; defined {
				o.Value = ResourceAttributeReference{}
				parseOutputReference(mod, attr.Expr, &o.Value)
			}

		case "provider":

			content, _, contentDiags := block.Body.PartialContent(providerConfigSchema)
			diags = append(diags, contentDiags...)

			name := block.Labels[0]
			// Even if there isn't an explicit version required, we still
			// need an entry in our map to signal the unversioned dependency.
			if _, exists := mod.RequiredProviders[name]; !exists {
				mod.RequiredProviders[name] = &ProviderRequirement{}
			}
			if attr, defined := content.Attributes["version"]; defined {
				var version string
				valDiags := gohcl.DecodeExpression(attr.Expr, nil, &version)
				diags = append(diags, valDiags...)
				if !valDiags.HasErrors() {
					mod.RequiredProviders[name].VersionConstraints = append(mod.RequiredProviders[name].VersionConstraints, version)
				}
			}

			providerKey := name
			var alias string
			if attr, defined := content.Attributes["alias"]; defined {
				valDiags := gohcl.DecodeExpression(attr.Expr, nil, &alias)
				diags = append(diags, valDiags...)
				if !valDiags.HasErrors() && alias != "" {
					providerKey = fmt.Sprintf("%s.%s", name, alias)
				}
			}

			mod.ProviderConfigs[providerKey] = &ProviderConfig{
				Name:  name,
				Alias: alias,
			}

		case "resource", "data":

			content, _, contentDiags := block.Body.PartialContent(resourceSchema)
			diags = append(diags, contentDiags...)

			typeName := block.Labels[0]
			name := block.Labels[1]

			r := &Resource{
				Type:       typeName,
				Name:       name,
				Pos:        sourcePosHCL(block.DefRange),
				References: make(map[string][]AttributeReference),
			}

			var resourcesMap map[string]*Resource

			switch block.Type {
			case "resource":
				r.Mode = ManagedResourceMode
				resourcesMap = mod.ManagedResources
			case "data":
				r.Mode = DataResourceMode
				resourcesMap = mod.DataResources
			}

			key := r.MapKey()

			resourcesMap[key] = r

			if attr, defined := content.Attributes["provider"]; defined {
				// New style here is to provide this as a naked traversal
				// expression, but we also support quoted references for
				// older configurations that predated this convention.
				traversal, travDiags := hcl.AbsTraversalForExpr(attr.Expr)
				if travDiags.HasErrors() {
					traversal = nil // in case we got any partial results

					// Fall back on trying to parse as a string
					var travStr string
					valDiags := gohcl.DecodeExpression(attr.Expr, nil, &travStr)
					if !valDiags.HasErrors() {
						var strDiags hcl.Diagnostics
						traversal, strDiags = hclsyntax.ParseTraversalAbs([]byte(travStr), "", hcl.Pos{})
						if strDiags.HasErrors() {
							traversal = nil
						}
					}
				}

				// If we get out here with a nil traversal then we didn't
				// succeed in processing the input.
				if len(traversal) > 0 {
					providerName := traversal.RootName()
					alias := ""
					if len(traversal) > 1 {
						if getAttr, ok := traversal[1].(hcl.TraverseAttr); ok {
							alias = getAttr.Name
						}
					}
					r.Provider = ProviderRef{
						Name:  providerName,
						Alias: alias,
					}
				} else {
					diags = append(diags, &hcl.Diagnostic{
						Severity: hcl.DiagError,
						Summary:  "Invalid provider reference",
						Detail:   "Provider argument requires a provider name followed by an optional alias, like \"aws.foo\".",
						Subject:  attr.Expr.Range().Ptr(),
					})
				}
			} else {
				// If provider _isn't_ set then we'll infer it from the
				// resource type.
				r.Provider = ProviderRef{
					Name: resourceTypeDefaultProviderName(r.Type),
				}
			}

			switch block.Type {
			case "resource":
				switch t := block.Body.(type) {
				case *hclsyntax.Body:
					r.Inputs = make(map[string]ResourceAttributeReference)
					pathRoot := ""
					for _, attr := range t.Attributes {
						parseResourceAttribute(mod, r, pathRoot, attr)
					}
					parseNestedResourceBlocks(mod, r, pathRoot, t.Blocks)
				}
			}

		case "module":

			content, _, contentDiags := block.Body.PartialContent(moduleCallSchema)
			diags = append(diags, contentDiags...)

			name := block.Labels[0]
			mc := &ModuleCall{
				Name: block.Labels[0],
				Pos:  sourcePosHCL(block.DefRange),
			}

			// check if this is overriding an existing module
			var origSource string
			if origMod, exists := mod.ModuleCalls[name]; exists {
				origSource = origMod.Source
			}

			mod.ModuleCalls[name] = mc

			if attr, defined := content.Attributes["source"]; defined {
				var source string
				valDiags := gohcl.DecodeExpression(attr.Expr, nil, &source)
				diags = append(diags, valDiags...)
				mc.Source = source
			}

			if mc.Source == "" {
				mc.Source = origSource
			}

			if attr, defined := content.Attributes["version"]; defined {
				var version string
				valDiags := gohcl.DecodeExpression(attr.Expr, nil, &version)
				diags = append(diags, valDiags...)
				mc.Version = version
			}

			// recursively parse local module references
			if subModPath := resolvedModuleRefs.Get(mc.Source, mc.Version); subModPath != "" {
				mc.Module, _ = LoadModule(subModPath, resolvedModuleRefs)
			}

			contentAttrs, contentDiags := block.Body.JustAttributes()
			diags = append(diags, contentDiags...)
			mc.Inputs = make(map[string]ResourceAttributeReference)
			for _, attr := range contentAttrs {
				parseModuleCallAttribute(mod, mc, attr)
			}

		default:
			// Should never happen because our cases above should be
			// exhaustive for our schema.
			panic(fmt.Errorf("unhandled block type %q", block.Type))
		}
	}

	return diags
}

func parseNestedResourceBlocks(parentModule *Module, parentResource *Resource, pathRoot string, blocks hclsyntax.Blocks) {
	for _, block := range blocks {
		blockName := block.Type
		if blockName == "dynamic" {
			blockName = block.Labels[0]
		} else if blockName == "content" {
			blockName = "" // ignore this level; defines body for iterators
		}

		newPathRoot := buildPath(pathRoot, blockName)
		for _, attr := range block.Body.Attributes {
			parseResourceAttribute(parentModule, parentResource, newPathRoot, attr)
		}

		parseNestedResourceBlocks(parentModule, parentResource, newPathRoot, block.Body.Blocks)
	}
}

func buildPath(segments ...string) string {
	nonEmptySegments := make([]string, 0)
	for _, segment := range segments {
		if segment != "" {
			nonEmptySegments = append(nonEmptySegments, segment)
		}
	}
	return strings.Join(nonEmptySegments, ".")
}

func parseResourceAttribute(parentModule *Module, parentResource *Resource, pathRoot string, attr *hclsyntax.Attribute) {
	switch attr.Name {
	case "tags", "count", "depends_on", "for_each":
		return // ignore common attributes
	default:
		qualifiedAttrName := buildPath(pathRoot, attr.Name)
		resolveResourceInputReference(parentModule, parentResource, qualifiedAttrName, attr.Expr, passOneLoggerPrefix)
	}
}

func resolveResourceInputReference(parentModule *Module, parentResource *Resource, qualifiedAttrName string, expr hcl.Expression, loggerPrefix string) {
	in := ResourceAttributeReference{}
	parseOutputReference(parentModule, expr, &in)
	parentResource.Inputs[qualifiedAttrName] = in

	// resource "res" {
	//	 attr_1 = var.input_variable (input variable whose resolution is not known at this level - i.e. from user or a higher-level module-call)
	//	 attr_2 = module.other_mod.out_val (reference to other module's output that has not been resolved yet - i.e. the referenced module has not been parsed yet)
	//	 attr_3 = other_resource.attribute.path (direct resource reference - other-resource to this-resource mapping)
	// }
	if in.ResourceType == "" || in.ResourceName == "" {
		// unknown reference type or literal
		return
	}

	thisResourceReference := ResourceAttributeReference{
		Module:        parentModule,
		ResourceType:  parentResource.Type,
		ResourceName:  parentResource.Name,
		AttributePath: []string{qualifiedAttrName},
	}

	switch in.ResourceType {
	case "var":
		current, ok := parentModule.Inputs[in.ResourceName]
		if !ok {
			current = make(map[string][]ResourceAttributeReference)
			parentModule.Inputs[in.ResourceName] = current
		}

		attrPath := strings.Join(in.AttributePath, ".")
		current[attrPath] = append(current[attrPath], thisResourceReference)
	case "local", "module":
		// this may not be available before all module calls get resolved (finish during 2nd pass)
		fmt.Printf("WARN [%s]: unresolved module-output or local-variable reference in resource %s.%s.%s'\n", loggerPrefix, parentResource.Type, parentResource.Name, qualifiedAttrName)
	case "each":
		// relative reference to for_each attribute value
	default:
		if ok := parentModule.AddResourceReference(in, thisResourceReference); ok {
			fmt.Printf("INFO [%s]: resource-to-resource found: %s -> %s\n", loggerPrefix, in.Attribute(), thisResourceReference.Attribute())
		} else {
			fmt.Printf("WARN [%s]: resource-to-resource : %s -> %s\n", loggerPrefix, in.Attribute(), thisResourceReference.Attribute())
		}
	}
}

func parseModuleCallAttribute(parentModule *Module, parentModuleCall *ModuleCall, attr *hcl.Attribute) {
	switch attr.Name {
	case "tags", "count", "depends_on", "for_each":
		return // ignore common attributes
	default:
		resolveModuleCallInputReference(parentModule, parentModuleCall, attr.Name, attr.Expr, passOneLoggerPrefix)
	}
}

func resolveModuleCallInputReference(parentModule *Module, parentModuleCall *ModuleCall, qualifiedAttrName string, expr hcl.Expression, loggerPrefix string) {
	in := ResourceAttributeReference{}
	parseOutputReference(parentModule, expr, &in)
	parentModuleCall.Inputs[qualifiedAttrName] = in

	// module "mod" {
	//	 attr_1 = var.input_variable (input variable whose resolution is not known at this level - i.e. from user or a higher-level module-call)
	//	 attr_2 = module.other_mod.out_val (reference to other module's output that has not been resolved yet - i.e. the referenced module has not been parsed yet)
	//	 attr_3 = other_resource.attribute.path (direct resource reference - knowing what attr_3 resolves to inside the called module we can compute resource-to-resource mapping)
	// }
	if parentModuleCall.Module != nil {
		if in.ResourceType == "" || in.ResourceName == "" {
			// unknown reference type, literal or special value (e.g. destroy, each)
			return
		}

		resolved := parentModuleCall.Module.GetResourceAttributeReferences(qualifiedAttrName) // attr.Name is the input variable's name inside the module

		switch in.ResourceType {
		case "var": // resolve this variable to the underlying resource in the called sub-module
			for _, item := range resolved {
				current, ok := parentModule.Inputs[in.ResourceName]
				if !ok {
					current = make(map[string][]ResourceAttributeReference)
					parentModule.Inputs[in.ResourceName] = current
				}

				// attrPath := strings.Join(append(append([]string{}, in.AttributePath...), item.RelativePath...), ".") // e.g. var.name.all.nested.attributes
				attrPath := in.MakeRelative(item.RelativePath).Path() // e.g. var.name.all.nested.attributes
				current[attrPath] = append(current[attrPath], item.ResourceAttributeReference)
			}
		case "local", "module":
			fmt.Printf("WARN [%s]: unresolved module-output or local-variable reference in module-call '%s.%s'\n", loggerPrefix, parentModuleCall.Name, qualifiedAttrName)
		case "each":
		// relative reference to for_each attribute value
		default:
			for _, item := range resolved {
				relInp := in.MakeRelative(item.RelativePath)
				if ok := item.ResourceAttributeReference.Module.AddResourceReference(relInp, item.ResourceAttributeReference); ok {
					fmt.Printf("INFO [%s]: resource-to-resource found: %s -> %s\n", loggerPrefix, relInp.Attribute(), item.ResourceAttributeReference.Attribute())
				} else {
					fmt.Printf("WARN [%s]: resource-to-resource : %s -> %s\n", loggerPrefix, relInp.Attribute(), item.ResourceAttributeReference.Attribute())
				}
			}
		}
	} else {
		fmt.Printf("WARN [%s]: unresolved module-call '%s' ('%s', '%s')\n", loggerPrefix, parentModuleCall.Name, parentModuleCall.Source, parentModuleCall.Version)
	}
}

func parseOutputReference(parent *Module, expr hcl.Expression, out *ResourceAttributeReference) {
	out.Expression = expr // save the original parsed expression
	out.Module = parent
	switch t := expr.(type) {
	case *hclsyntax.ScopeTraversalExpr:
		out.ResourceType = t.Traversal.RootName()
		traversals := t.Traversal
		if out.ResourceType == "data" {
			traversals = traversals[1:] // data resources are referenced as data.<type>.<name>.<attribute>
		}
		for i := range traversals {
			switch tr := traversals[i].(type) {
			case hcl.TraverseAttr:
				if i == 0 {
					out.ResourceType = tr.Name // data-resource case
				} else if i == 1 {
					switch out.ResourceType {
					case "local":
						if val, ok := parent.Locals[tr.Name]; ok {
							parseOutputReference(parent, val, out)
							continue
						}
					case "module":
						if submod, ok := parent.ModuleCalls[tr.Name]; ok {
							relAttrs := []string{}
							parseAttributes(traversals[2:], &relAttrs)
							if submod.Module != nil && len(relAttrs) > 0 {
								if resolvedOutput, ok := submod.Module.Outputs[relAttrs[0]]; ok {
									*out = resolvedOutput.Value
									out.AttributePath = append(append([]string{}, out.AttributePath...), relAttrs[1:]...) // append rest of the attribute path to the resolved reference
									return
								}
							}
						}
					}
					out.ResourceName = tr.Name
				} else {
					out.AttributePath = append(out.AttributePath, tr.Name)
				}
			}
		}
	case *hclsyntax.FunctionCallExpr:
		switch t.Name {
		case "lookup":
			parseOutputReference(parent, t.Args[0], out)
			parseOutputReference(parent, t.Args[1], out) // the lookup attribute
		default:
			parseOutputReference(parent, t.Args[0], out) // only do first arg
		}
	case *hclsyntax.SplatExpr:
		parseOutputReference(parent, t.Source, out)
		parseOutputReference(parent, t.Each, out)
	case *hclsyntax.IndexExpr:
		parseOutputReference(parent, t.Collection, out)
	case *hclsyntax.RelativeTraversalExpr:
		parseOutputReference(parent, t.Source, out)
		parseAttributes(t.Traversal, &out.AttributePath)
	case *hclsyntax.ConditionalExpr:
		parseOutputReference(parent, t.TrueResult, out) // only do True branch
	case *hclsyntax.ForExpr:
		// [for KeyVar, ValVar in <CollExpr (e.g. module.attr)> : <ValExpr (e.g. <ValVar>.attr)>]
		collectionExprRef := ResourceAttributeReference{}
		if parseOutputReference(parent, t.ValExpr, out); out.ResourceType == t.ValVar {
			// iterator has a value variable that refers to the CollExpr: [for k, v in module.module_name.ids : v.id]
			parseOutputReference(parent, t.CollExpr, &collectionExprRef)
		} else if parseOutputReference(parent, t.KeyExpr, out); out.ResourceType == t.KeyVar {
			// if not successful above, attempt the same for KeyExpr
			parseOutputReference(parent, t.CollExpr, &collectionExprRef)
		}

		// resolve the item expression by appending it's path to the CollExpr: [for k, v in module.module_name.ids : v.id] --> module.module_name.ids
		// if item expression does not refer to CollExpr; [for k, v in module.module_name.ids : local.vpc_cidr] --> local.vpc_cidr
		// if out.ResourceName is empty this is probably a collection of literals: [for k, v in module.module_name.ids : "hello"]
		if out.ResourceName != "" && collectionExprRef.ResourceName != "" {
			collectionExprRef.AttributePath = append(append(collectionExprRef.AttributePath, out.ResourceName), out.AttributePath...)
			*out = collectionExprRef
		}
	case *hclsyntax.TemplateExpr:
		if t.IsStringLiteral() {
			v, _ := t.Value(nil)
			out.AttributePath = append(out.AttributePath, v.AsString())
		}
	case *hclsyntax.TemplateWrapExpr:
		parseOutputReference(parent, t.Wrapped, out)
	case *hclsyntax.TupleConsExpr:
		if t.Exprs != nil { // empty tuple []
			parseOutputReference(parent, t.Exprs[0], out) // only do first item in the collection
		}
	}
}

func parseAttributes(tr hcl.Traversal, out *[]string) {
	for _, traversal := range tr {
		switch t := traversal.(type) {
		case hcl.TraverseAttr:
			*out = append(*out, t.Name)
		case hcl.TraverseIndex:
			if t.Key.Type() == cty.String {
				*out = append(*out, t.Key.AsString()) // e.g. map key
			}
		}
	}
}
