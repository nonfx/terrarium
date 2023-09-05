//go:build mock
// +build mock

package mappings

import (
	"context"
	"testing"

	"github.com/cldcvr/terrarium/src/cli/internal/config"
	"github.com/cldcvr/terrarium/src/pkg/commander"
	"github.com/cldcvr/terrarium/src/pkg/commander/mocks"
	dbmocks "github.com/cldcvr/terrarium/src/pkg/db/mocks"
	"github.com/cldcvr/terrarium/src/pkg/testutils/clitesting"
	"github.com/rotisserie/eris"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/mock"
)

func TestCmd(t *testing.T) {
	testSetup := clitesting.CLITest{
		CmdToTest: NewCmd,
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
			Name: "error loading modules list YAML file",
			PreExecute: func(ctx context.Context, t *testing.T, cmd *cobra.Command, cmdOpts clitesting.CmdOpts) {
				mockedDB := &dbmocks.DB{}
				config.SetDBMocks(mockedDB)
			},
			Args:     []string{"-f", "invalid-file-path"},
			WantErr:  true,
			ExpError: "failed to load farm module list file",
		},
		{
			Name: "error running tf command",
			PreExecute: func(ctx context.Context, t *testing.T, cmd *cobra.Command, cmdOpts clitesting.CmdOpts) {
				mockedDB := &dbmocks.DB{}
				config.SetDBMocks(mockedDB)

				mockedCmdr := &mocks.Commander{}
				mockedCmdr.On("Run", mock.Anything).Return(eris.New("mocked error")).Once()

				commander.SetCommander(mockedCmdr)
			},
			Args:     []string{"-f", "./testdata/modules.yaml"},
			WantErr:  true,
			ExpError: "mocked error",
		},
	})
}
