// Copyright (c) Ollion
// SPDX-License-Identifier: Apache-2.0

package generate

import (
	"context"
	"os"
	"path"
	"testing"

	"github.com/cldcvr/terrarium/src/pkg/testutils/clitesting"
	"github.com/stretchr/testify/assert"
)

func TestCmd(t *testing.T) {
	os.RemoveAll("./testdata/.terrarium")
	testSetup := clitesting.CLITest{
		CmdToTest: NewCmd,
		TeardownTestCase: func(ctx context.Context, t *testing.T, tc clitesting.CLITestCase) {
			os.RemoveAll("./testdata/.terrarium")
		},
	}

	testSetup.RunTests(t, []clitesting.CLITestCase{
		{
			Name:     "No components provided",
			Args:     []string{},
			WantErr:  true,
			ExpError: "No Apps provided. use -a flag to set apps",
		},
		{
			Name:     "Invalid app path",
			Args:     []string{"-a", "./invalid-path"},
			WantErr:  true,
			ExpError: "invalid file path: ./invalid-path",
		},
		{
			Name: "Success (no env files)",
			Args: []string{"-p", "../../../../examples/platform/", "-a", "../../../../examples/apps/voting-be", "-a", "../../../../examples/apps/voting-worker", "-o", "./testdata/.terrarium", "--skip-env-file"},
			ValidateOutput: func(ctx context.Context, t *testing.T, cmdOpts clitesting.CmdOpts, output []byte) bool {
				pass := assert.Equal(t, "Successfully pulled 13 of 22 terraform blocks at: ./testdata/.terrarium\n", string(output))
				pass = assertFilesExists(t,
					"./testdata/.terrarium",
					[]string{ // shouldExist
						"component_redis.tf",
						"outputs.tf",
						"tr_base_backend.tf",
						"tr_gen_locals.tf",
						"vpc.tf",
					},
					[]string{ // shouldNotExist
						"app_voting_be.env.mustache",
						"app_voting_worker.env.mustache",
						"tr_gen_profile.auto.tfvars",
						"component_postgres.tf",
					},
				) && pass
				return pass
			},
		},
		{
			Name: "Success (no profile)",
			Args: []string{"-p", "../../../../examples/platform/", "-a", "../../../../examples/apps/voting-be", "-a", "../../../../examples/apps/voting-worker", "-o", "./testdata/.terrarium"},
			ValidateOutput: func(ctx context.Context, t *testing.T, cmdOpts clitesting.CmdOpts, output []byte) bool {
				pass := assert.Equal(t, "Successfully pulled 13 of 22 terraform blocks at: ./testdata/.terrarium\n", string(output))
				pass = assertFilesExists(t,
					"./testdata/.terrarium",
					[]string{ // shouldExist
						"app_voting_be.env.mustache",
						"app_voting_worker.env.mustache",
						"component_redis.tf",
						"outputs.tf",
						"tr_base_backend.tf",
						"tr_gen_locals.tf",
						"vpc.tf",
					},
					[]string{ // shouldNotExist
						"tr_gen_profile.auto.tfvars",
						"component_postgres.tf",
					},
				) && pass
				return pass
			},
		},
		{
			Name: "Success (with profile)",
			Args: []string{"-p", "../../../../examples/platform/", "-a", "../../../../examples/apps/voting-be", "-a", "../../../../examples/apps/voting-worker", "-o", "./testdata/.terrarium", "-c", "dev"},
			ValidateOutput: func(ctx context.Context, t *testing.T, cmdOpts clitesting.CmdOpts, output []byte) bool {
				pass := assert.Equal(t, "Successfully pulled 13 of 22 terraform blocks at: ./testdata/.terrarium\n", string(output))
				pass = assertFilesExists(t,
					"./testdata/.terrarium",
					[]string{ // shouldExist
						"app_voting_be.env.mustache",
						"app_voting_worker.env.mustache",
						"component_redis.tf",
						"outputs.tf",
						"tr_base_backend.tf",
						"tr_gen_locals.tf",
						"tr_gen_profile.auto.tfvars",
						"vpc.tf",
					},
					[]string{ // shouldNotExist
						"component_postgres.tf",
					},
				) && pass
				return pass
			},
		},
		{
			Name:     "Invalid profile name",
			Args:     []string{"-p", "../../../../examples/platform/", "-a", "../../../../examples/apps/voting-be", "-a", "../../../../examples/apps/voting-worker", "-o", "./testdata/.terrarium", "-c", "Isle"},
			WantErr:  true,
			ExpError: "could not retrieve configuration file for platform profile 'Isle'",
		},
	})
}

func assertFilesExists(t *testing.T, dir string, shouldExist, shouldNotExist []string) bool {
	t.Helper()

	pass := true
	for _, f := range shouldExist {
		filePath := path.Join(dir, f)
		pass = assert.FileExists(t, filePath) && pass
	}

	for _, f := range shouldNotExist {
		filePath := path.Join(dir, f)
		pass = assert.NoFileExists(t, filePath) && pass
	}

	return pass
}
