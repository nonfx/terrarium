// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package generate

import (
	"testing"

	"github.com/cldcvr/terrarium/src/pkg/testutils/clitesting"
)

func TestCmd(t *testing.T) {
	testSetup := clitesting.CLITest{
		CmdToTest: NewCmd,
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
			Name:           "Success (no profile)",
			Args:           []string{"-p", "../../../../examples/platform/", "-a", "../../../../examples/apps/voting-be", "-a", "../../../../examples/apps/voting-worker", "-o", "./testdata/.terrarium"},
			ValidateOutput: clitesting.ValidateOutputMatch("Successfully pulled 13 of 22 terraform blocks at: ./testdata/.terrarium\n"),
		},
		{
			Name:           "Success (with profile)",
			Args:           []string{"-p", "../../../../examples/platform/", "-a", "../../../../examples/apps/voting-be", "-a", "../../../../examples/apps/voting-worker", "-o", "./testdata/.terrarium", "-c", "low-cost"},
			ValidateOutput: clitesting.ValidateOutputMatch("Successfully pulled 13 of 22 terraform blocks at: ./testdata/.terrarium\n"),
		},
		{
			Name:     "Invalid profile name",
			Args:     []string{"-p", "../../../../examples/platform/", "-a", "../../../../examples/apps/voting-be", "-a", "../../../../examples/apps/voting-worker", "-o", "./testdata/.terrarium", "-c", "Isle"},
			WantErr:  true,
			ExpError: "could not retrieve configuration file for platform profile 'Isle'",
		},
	})
}
