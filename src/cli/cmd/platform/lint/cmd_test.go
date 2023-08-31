// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package lint

import (
	"testing"

	"github.com/cldcvr/terrarium/src/pkg/testutils/clitesting"
)

func TestCmd(t *testing.T) {
	clitest := clitesting.CLITest{
		CmdToTest: NewCmd,
	}

	clitest.RunTests(t, []clitesting.CLITestCase{
		{
			Name:           "render help",
			Args:           []string{"-h"},
			ValidateOutput: clitesting.ValidateOutputContains("Analyze the directory and verify it constitutes a valid platform definition.\n\nUsage:\n  lint [flags]"),
		},
		{
			Name:     "invalid directory",
			Args:     []string{"-d", "testdata/invalid-path"},
			WantErr:  true,
			ExpError: "could not open given directory",
		},
		{
			Name:           "valid platform",
			Args:           []string{"-d", "testdata/valid-terraform-1"},
			ValidateOutput: clitesting.ValidateOutputMatch("Platform parse and lint completed\n"),
		},
	})
}
