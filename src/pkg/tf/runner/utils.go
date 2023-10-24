// Copyright (c) Ollion
// SPDX-License-Identifier: Apache-2.0

package runner

import "github.com/rotisserie/eris"

func TerraformProviderSchema(runner *terraformRunner, dir string, outFilePath string) error {
	if err := TerraformInit(runner, dir); err != nil {
		return eris.Wrap(err, "error initilizing terraform working directory")
	}
	if err := runner.RunTerraformProvidersSchema(dir, outFilePath); err != nil {
		return eris.Wrap(err, "error getting terraform providers schema")
	}

	return nil
}

func TerraformInit(runner *terraformRunner, dir string) error {
	if err := runner.RunTerraformVersion(); err != nil {
		return eris.Wrap(err, "error getting terraform version")
	}
	if err := runner.RunTerraformInit(dir); err != nil {
		return eris.Wrap(err, "error running terraform init")
	}

	return nil
}
