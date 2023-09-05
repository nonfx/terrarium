package harvest

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
			ValidateOutput: clitesting.ValidateOutputContains("The 'harvest' command provides subcommands to scan Terraform providers, modules, mappings, and more from the specified farm directory.\n\nUsage:\n  harvest [command]"),
		},
	})
}
