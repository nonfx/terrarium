// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package runner

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/cldcvr/terrarium/src/pkg/commander"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/rotisserie/eris"
	"github.com/spf13/viper"
)

func NewTerraformRunner() *terraformRunner {
	cacheDir := viper.GetString("terraform.plugin_cache_dir")
	terraformDefaultEnv := []string{
		"TF_IN_AUTOMATION=1",
		"TF_INPUT=0",
	}
	if cacheDir != "" {
		resolvedCacheDir, err := homedir.Expand(cacheDir)
		if err != nil {
			log.Fatalf("could not open user home directory: %v", err)
		} else if err := os.MkdirAll(resolvedCacheDir, os.ModePerm); err != nil {
			log.Fatalf("could not create Terraform plugin_cache_dir '%s': %v", resolvedCacheDir, err)
		} else {
			terraformDefaultEnv = append(terraformDefaultEnv, fmt.Sprintf("TF_PLUGIN_CACHE_DIR=%s", resolvedCacheDir))
		}
	}

	return &terraformRunner{
		terraformDefaultEnv: terraformDefaultEnv,
	}
}

type terraformRunner struct {
	terraformDefaultEnv []string
}

func (tr *terraformRunner) RunTerraformVersion() error {
	return tr.runTerraformCommandWithDefaultEnv("", []string{"version"}, nil)
}

func (tr *terraformRunner) RunTerraformInit(dir string) error {
	return tr.runTerraformCommandWithDefaultEnv(dir, []string{"init"}, nil)
}

func (tr *terraformRunner) RunTerraformProviders(dir string) error {
	return tr.runTerraformCommandWithDefaultEnv(dir, []string{"providers"}, nil)
}

func (tr *terraformRunner) RunTerraformProvidersSchema(tfDir string, outFilePath string) error {
	fileDir := path.Dir(outFilePath)
	err := os.MkdirAll(fileDir, os.ModePerm)
	if err != nil {
		return eris.Wrapf(err, "failed to create directory path '%s'", fileDir)
	}

	out, err := os.Create(outFilePath)
	if err != nil {
		return eris.Wrapf(err, "could not create file '%s'", outFilePath)
	}

	return tr.runTerraformCommandWithDefaultEnv(tfDir, []string{"providers", "schema", "-json"}, out)
}

func (tr *terraformRunner) runTerraformCommandWithDefaultEnv(dir string, args []string, outWriter io.Writer) error {
	return tr.runTerraformCommand(dir, args, tr.terraformDefaultEnv, outWriter)
}

func (tr *terraformRunner) runTerraformCommand(dir string, args []string, env []string, outWriter io.Writer) error {
	cmd := exec.Command("terraform", args...)
	cmd.Dir = dir
	cmd.Env = append(cmd.Env, os.Environ()...) // inherit the user's ENV
	cmd.Env = append(cmd.Env, env...)
	if outWriter == nil {
		outWriter = os.Stdout
	}
	cmd.Stdout = outWriter
	cmd.Stderr = os.Stderr
	log.Printf("[%s] %s", strings.Join(env, ", "), cmd.String())
	if err := commander.GetCommander().Run(cmd); err != nil {
		return eris.Wrapf(err, "command '%v' finished with error", cmd)
	}
	return nil
}
