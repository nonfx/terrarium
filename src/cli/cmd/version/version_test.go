package version

import (
	"context"
	"testing"

	"github.com/cldcvr/terrarium/src/cli/internal/build"
	"github.com/cldcvr/terrarium/src/pkg/testutils/clitesting"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_CmdVersion(t *testing.T) {
	clitest := clitesting.CLITest{
		CmdToTest: GetCmd(),
	}

	tests := []clitesting.CLITestCase{
		{
			Name: "Version DEV",
			ValidateOutput: func(ctx context.Context, t *testing.T, cmdOpts clitesting.CmdOpts, output []byte) bool {
				return assert.Equal(t, "terrarium version DEV\n", string(output))
			},
		},
		{
			Name: "Version info set",
			PreExecute: func(ctx context.Context, t *testing.T, cmd *cobra.Command, cmdOpts clitesting.CmdOpts) {
				build.Date = "2020-12-01"
				build.Version = "5e08f18"
			},
			ValidateOutput: func(ctx context.Context, t *testing.T, cmdOpts clitesting.CmdOpts, output []byte) bool {
				return assert.Equal(t, "terrarium version 5e08f18 (2020-12-01)\n", string(output))
			},
		},
	}
	clitest.RunTests(t, tests)
}
