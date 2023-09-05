package platform

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
			ValidateOutput: clitesting.ValidateOutputContains("Commands to manage Terrarium platform template definitions\n\nUsage:\n  platform [command]"),
		},
	})
}
