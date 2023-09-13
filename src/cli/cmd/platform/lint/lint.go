// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package lint

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/cldcvr/terraform-config-inspect/tfconfig"
	"github.com/cldcvr/terrarium/src/cli/internal/constants"
	"github.com/cldcvr/terrarium/src/pkg/metadata/platform"
	"github.com/cldcvr/terrarium/src/pkg/tf/parser"
	"github.com/hashicorp/hcl/v2"
	"github.com/rotisserie/eris"
	"gopkg.in/yaml.v3"
)

const terrariumComponentModulePrefix = "tr_component_"
const terrariumComponentModuleEnabledSuffix = "_enabled"
const terrariumTaxonEnabledPrefix = "tr_taxon_"
const terrariumTaxonEnabledSuffix = "_enabled"

func lintPlatform(dir string) error {
	log.Info("Linting terrarium platform template...")

	log.Infof("Loading Terraform modules to lint from '%s'...", dir)
	module, _ := tfconfig.LoadModule(dir, &tfconfig.ResolvedModulesSchema{})

	log.Info("Validating Terraform modules...")
	if err := validatePlatformTerraform(module); err != nil {
		log.Infof("Following Terraform issues were found: %v", err)
		return eris.Wrap(err, "platform lint")
	}
	log.Info("Platform is valid.")

	metadataFile := filepath.Join(dir, "terrarium.yaml")
	log.Infof("Loading Terrarium metadata file '%s'...", metadataFile)

	fileData, err := os.ReadFile(metadataFile)
	if os.IsNotExist(err) {
		// ignore not exists error since we create the metadata file anyway in this case.
		err = nil
	}

	pm, err := platform.NewPlatformMetadata(module, fileData)
	if err != nil {
		return eris.Wrap(err, "erro parsing platform metadata")
	}

	pmYAML, err := yaml.Marshal(pm)
	if err != nil {
		return eris.Wrap(err, "erro serializing platform metadata")
	}

	if string(fileData) == string(pmYAML) {
		log.Info("No change in metadata.")
		return nil
	}

	log.Infof("Updating metadata file at: %s", metadataFile)
	os.WriteFile(metadataFile, pmYAML, constants.ReadWritePermissions)

	log.Info("Metadata updated.")

	return nil
}

func validatePlatformTerraform(module *tfconfig.Module) error {
	requiredModuleNames := []string{}
	for name, expr := range module.Locals {
		// Find all auto-generated inputs and assert they are iterable.
		if strings.HasPrefix(name, terrariumComponentModulePrefix) && !strings.HasSuffix(name, terrariumComponentModuleEnabledSuffix) {
			if !parser.IsObject(expr.Expression) {
				return eris.Errorf("dependency input declaration '%s' %s must be iterable", name, fmtExpressionPosition(expr.Expression))
			}
			requiredModuleNames = append(requiredModuleNames, name)
		}

		// Assert taxon switch variables are boolean.
		if strings.HasPrefix(name, terrariumTaxonEnabledPrefix) && strings.HasSuffix(name, terrariumTaxonEnabledSuffix) {
			if !parser.IsBool(expr.Expression) {
				return eris.Errorf("terraform variable '%s' %s must evaluate to a boolean", name, fmtExpressionPosition(expr.Expression))
			}
		}
	}

	for _, name := range requiredModuleNames {
		switchVarName := name + terrariumComponentModuleEnabledSuffix
		if expr, found := module.Locals[switchVarName]; found && !parser.IsBool(expr.Expression) {
			return eris.Errorf("terraform variable '%s' %s must evaluate to a boolean: %s = length(local.%s) > 0", switchVarName, fmtExpressionPosition(expr.Expression), switchVarName, name)
		}

		// Verify that a module exists for each input map,
		if _, ok := module.ModuleCalls[name]; !ok {
			return eris.Errorf("terraform must implement '%s' component by declaring a module call with matching label: module \"%s\" { for_each = local.%s }", name, name, name)
		}
	}

	for name, output := range module.Outputs {
		if strings.HasPrefix(name, terrariumComponentModulePrefix) {
			// Ensure the output is an iterable map object.
			// The map will contain an output value for each instance of the dependency created.
			if !parser.IsCollection(output.Value.Expression) {
				return eris.Errorf("terraform output '%s' %s be a map", name, fmtExpressionPosition(output.Value.Expression))
			}
		}
	}

	return nil
}

func fmtExpressionPosition(expr hcl.Expression) string {
	r := expr.StartRange()
	loc := r.Start
	return fmtSourcePosition(r.Filename, loc.Line, loc.Column)
}

func fmtModuleCallPosition(mc *tfconfig.ModuleCall) string {
	return fmtSourcePosition(mc.Pos.Filename, mc.Pos.Line, 0)
}

func fmtSourcePosition(filename string, line int, column int) string {
	tokens := make([]string, 0)
	if filename != "" {
		tokens = append(tokens, fmt.Sprintf("filename = %s", filename))
	}
	if line != 0 {
		tokens = append(tokens, fmt.Sprintf("line = %d", line))
	}
	if column != 0 {
		tokens = append(tokens, fmt.Sprintf("column = %d", column))
	}

	return strings.Join(tokens, "; ")
}
