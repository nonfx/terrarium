//go:build mock
// +build mock

package dependencies

import (
	"context"
	"testing"

	"github.com/cldcvr/terrarium/src/cli/internal/config"
	"github.com/cldcvr/terrarium/src/pkg/db/mocks"
	"github.com/cldcvr/terrarium/src/pkg/testutils/clitesting"
	"github.com/spf13/cobra"
)

func TestCmd(t *testing.T) {
	testSetup := clitesting.CLITest{
		CmdToTest: GetCmd(),
		SetupTest: func(ctx context.Context, t *testing.T) {
			t.Setenv("TR_LOG_LEVEL", "error")
			config.LoadDefaults()
		},
	}

	testSetup.RunTests(t, []clitesting.CLITestCase{
		{
			Name: "error connecting to db",
			PreExecute: func(ctx context.Context, t *testing.T, cmd *cobra.Command, cmdOpts clitesting.CmdOpts) {
				config.SetDBMocks(nil)
			},
			Args:     []string{},
			WantErr:  true,
			ExpError: "error connecting to the database: mocked err: connection failed",
		},
		{
			Name: "error connecting to db",
			PreExecute: func(ctx context.Context, t *testing.T, cmd *cobra.Command, cmdOpts clitesting.CmdOpts) {
				config.SetDBMocks(&mocks.DB{})
			},
			Args:     []string{"-d", "./non-existing"},
			WantErr:  true,
			ExpError: "lstat ./non-existing: no such file or directory",
		},
	})
}
