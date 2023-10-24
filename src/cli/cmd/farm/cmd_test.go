// Copyright (c) Ollion
// SPDX-License-Identifier: Apache-2.0

package farm

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
			Name:           "farm help",
			ValidateOutput: clitesting.ValidateOutputContains("commands to access farm repository\n\nUsage:\n  farm [command]\n\nAvailable Commands:\n  completion  Generate the autocompletion script for the specified shell\n  help        Help about any command\n  update      farm update updates the databse with latest farm release\n\nFlags:\n  -h, --help   help for farm\n\nUse \"farm [command] --help\" for more information about a command.\n"),
		},
	})
}
