// Copyright (c) Ollion
// SPDX-License-Identifier: Apache-2.0

package query

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
			ValidateOutput: clitesting.ValidateOutputContains("commands to list matching modules as per the filter chosen. provides variety of filters to list desired modules\n\nUsage:\n  query [command]"),
		},
	})
}
