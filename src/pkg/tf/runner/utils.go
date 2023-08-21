package runner

func TerraformProviderSchema(runner *terraformRunner, dir string, outFilePath string) error {
	if err := TerraformInit(runner, dir); err != nil {
		return err
	}
	if err := runner.RunTerraformProvidersSchema(dir, outFilePath); err != nil {
		return err
	}

	return nil
}

func TerraformInit(runner *terraformRunner, dir string) error {
	if err := runner.RunTerraformVersion(); err != nil {
		return err
	}
	if err := runner.RunTerraformInit(dir); err != nil {
		return err
	}

	return nil
}
