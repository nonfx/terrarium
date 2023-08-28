package generate

import (
	"testing"

	"github.com/cldcvr/terrarium/src/pkg/testutils/clitesting"
)

func TestCmd(t *testing.T) {
	testSetup := clitesting.CLITest{
		CmdToTest: GetCmd(),
	}

	testSetup.RunTests(t, []clitesting.CLITestCase{
		{
			Name:     "No components provided",
			Args:     []string{},
			WantErr:  true,
			ExpError: "No Apps provided. use -a flag to set apps",
		},
		{
			Name:           "Success",
			Args:           []string{"-p", "../../../../examples/platform/", "-a", "../../../../examples/apps/voting-be", "-a", "../../../../examples/apps/voting-worker", "-o", "./testdata/.terrarium"},
			ValidateOutput: clitesting.ValidateOutputMatch("Successfully pulled 13 of 22 terraform blocks at: ./testdata/.terrarium\n"),
		},
	})
}
