// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfconfig

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/hcl/v2/hclsyntax"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	ctyjson "github.com/zclconf/go-cty/cty/json"
)

func loadModule(fs FS, dir string) (*Module, Diagnostics) {
	mod := NewModule(dir)
	primaryPaths, diags := dirFiles(fs, dir)

	parser := hclparse.NewParser()

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

		contentDiags := LoadModuleFromFile(file, mod)
		diags = append(diags, contentDiags...)
	}

	return mod, diagnosticsHCL(diags)
}

// LoadModuleFromFile reads given file, interprets it and stores in given Module
// This is useful for any caller which does tokenization/parsing on its own
// e.g. because it will reuse these parsed files later for more detailed
// interpretation.
func LoadModuleFromFile(file *hcl.File, mod *Module) hcl.Diagnostics {
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
				parseOutputReference(mod.Locals, attr.Expr, &o.Value)
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
				Type: typeName,
				Name: name,
				Pos:  sourcePosHCL(block.DefRange),
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
				content, contentDiags := block.Body.JustAttributes()
				diags = append(diags, contentDiags...)
				if !contentDiags.HasErrors() {
					r.Inputs = make(map[string]ResourceAttributeReference)
					for _, attr := range content {
						switch attr.Name {
						case "tags", "count", "depends_on":
							continue // ignore common attributes
						default:
							in := ResourceAttributeReference{}
							parseOutputReference(mod.Locals, attr.Expr, &in)
							r.Inputs[attr.Name] = in
						}
					}
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

		default:
			// Should never happen because our cases above should be
			// exhaustive for our schema.
			panic(fmt.Errorf("unhandled block type %q", block.Type))
		}
	}

	return diags
}

// func p(e hcl.Expression) {
// 	switch et := e.(type) {
// 	case *hclsyntax.SplatExpr:
// 		p(et.Each)
// 		p(et.Item)
// 	default:
// 		pp(et.Variables())
// 	}
// }

// func pp(vv []hcl.Traversal) {
// 	for _, qqq := range vv {
// 		if v, _ := qqq.TraverseAbs(nil); v.Type().IsPrimitiveType() {
// 			fmt.Print(v.AsString())
// 		}
// 		for _, ttt := range qqq {
// 			// switch t := ttt.(type) {
// 			// case hcl.TraverseRoot:
// 			// 	fmt.Printf("%v ->", t.Name)
// 			// case hcl.TraverseAttr:
// 			// 	fmt.Printf(".%v", t.Name)
// 			// case hcl.TraverseIndex:
// 			// 	// fmt.Printf(".%v", t.Key.GoString())
// 			// case hcl.TraverseSplat:
// 			// 	// fmt.Printf(".%v", t.Each.SimpleSplit().Rel.RootName())
// 			// default:
// 			// 	//
// 			// }
// 		}
// 		fmt.Printf("\n")
// 	}
// }

func parseOutputReference(locals map[string]hcl.Expression, expr hcl.Expression, out *ResourceAttributeReference) {
	switch t := expr.(type) {
	case *hclsyntax.ScopeTraversalExpr:
		out.ResourceType = t.Traversal.RootName()
		for i := range t.Traversal {
			switch t := t.Traversal[i].(type) {
			case hcl.TraverseAttr:
				if i == 1 {
					if out.ResourceType == "local" {
						if val, ok := locals[t.Name]; ok {
							parseOutputReference(locals, val, out)
							continue
						}
					}
					out.ResourceName = t.Name
				} else {
					out.AttributePath = append(out.AttributePath, t.Name)
				}
			}
		}
	case *hclsyntax.FunctionCallExpr:
		switch t.Name {
		case "lookup":
			parseOutputReference(locals, t.Args[0], out)
			parseOutputReference(locals, t.Args[1], out) // the lookup attribute
		default:
			parseOutputReference(locals, t.Args[0], out) // only do first arg
		}
	case *hclsyntax.SplatExpr:
		parseOutputReference(locals, t.Source, out)
		parseOutputReference(locals, t.Each, out)
	case *hclsyntax.IndexExpr:
		parseOutputReference(locals, t.Collection, out)
	case *hclsyntax.RelativeTraversalExpr:
		parseOutputReference(locals, t.Source, out)
		for _, traversal := range t.Traversal {
			switch t := traversal.(type) {
			case hcl.TraverseAttr:
				out.AttributePath = append(out.AttributePath, t.Name)
			case hcl.TraverseIndex:
				out.AttributePath = append(out.AttributePath, t.Key.AsString()) // e.g. map key
			}
		}
	case *hclsyntax.ConditionalExpr:
		parseOutputReference(locals, t.TrueResult, out) // only do True branch
	case *hclsyntax.ForExpr:
		tmp := ResourceAttributeReference{}
		parseOutputReference(locals, t.CollExpr, out)
		parseOutputReference(locals, t.ValExpr, &tmp)
		out.AttributePath = append(out.AttributePath, tmp.ResourceName)
		out.AttributePath = append(out.AttributePath, tmp.AttributePath...)
	case *hclsyntax.TemplateExpr:
		if t.IsStringLiteral() {
			v, _ := t.Value(nil)
			out.AttributePath = append(out.AttributePath, v.AsString())
		}
	}
}
