//go:build mock
// +build mock

package resources

import (
	"context"
	"os"
	"os/exec"
	"testing"

	"github.com/cldcvr/terrarium/src/cli/internal/config"
	"github.com/cldcvr/terrarium/src/pkg/commander"
	"github.com/cldcvr/terrarium/src/pkg/commander/mocks"
	dbmocks "github.com/cldcvr/terrarium/src/pkg/db/mocks"
	"github.com/cldcvr/terrarium/src/pkg/testutils/clitesting"
	"github.com/google/uuid"
	"github.com/rotisserie/eris"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
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
			Name: "error loading provider schema file",
			PreExecute: func(ctx context.Context, t *testing.T, cmd *cobra.Command, cmdOpts clitesting.CmdOpts) {
				mockedDB := &dbmocks.DB{}
				config.SetDBMocks(mockedDB)
			},
			Args:     []string{"-s", "invalid-file-path"},
			WantErr:  true,
			ExpError: "error loading providers schema file",
		},
		{
			Name: "error updating db",
			PreExecute: func(ctx context.Context, t *testing.T, cmd *cobra.Command, cmdOpts clitesting.CmdOpts) {
				mockedDB := &dbmocks.DB{}
				mockedDB.On("GetOrCreateTFProvider", mock.Anything).Return(uuid.New(), false, eris.New("mock error"))
				config.SetDBMocks(mockedDB)
			},
			Args:     []string{"-s", "./testdata/example_schema.json"},
			WantErr:  true,
			ExpError: "error writing data to db",
		},
		{
			Name: "success with schema file",
			PreExecute: func(ctx context.Context, t *testing.T, cmd *cobra.Command, cmdOpts clitesting.CmdOpts) {
				mockedDB := &dbmocks.DB{}
				mockedDB.On("GetOrCreateTFProvider", mock.Anything).Return(uuid.New(), false, nil)
				config.SetDBMocks(mockedDB)
			},
			Args:           []string{"-s", "./testdata/example_schema.json"},
			ValidateOutput: clitesting.ValidateOutputContains("Successfully added 0 Providers, 0 Resources, and 0 Attributes.\n"),
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
				mockedDB.On("GetOrCreateTFProvider", mock.Anything).Return(uuid.New(), false, nil)
				config.SetDBMocks(mockedDB)

				mockedCmdr := &mocks.Commander{}
				mockedCmdr.On("Run", mock.Anything).Return(eris.New("mocked error")).Once()

				commander.SetCommander(mockedCmdr)
			},
			Args:     []string{"-f", "./testdata/modules.yaml"},
			WantErr:  true,
			ExpError: "mocked error",
		},
		{
			Name: "error loading provider-schema from module YAML",
			PreExecute: func(ctx context.Context, t *testing.T, cmd *cobra.Command, cmdOpts clitesting.CmdOpts) {
				mockedDB := &dbmocks.DB{}
				mockedDB.On("GetOrCreateTFProvider", mock.Anything).Return(uuid.New(), false, nil)
				config.SetDBMocks(mockedDB)

				mockedCmdr := &mocks.Commander{}
				mockedCmdr.On("Run", mock.Anything).Return(nil).Times(3)

				commander.SetCommander(mockedCmdr)
			},
			Args:     []string{"-f", "./testdata/modules.yaml"},
			WantErr:  true,
			ExpError: "error loading providers schema file",
		},
		{
			Name: "success with modules list YAML file",
			PreExecute: func(ctx context.Context, t *testing.T, cmd *cobra.Command, cmdOpts clitesting.CmdOpts) {
				mockedDB := &dbmocks.DB{}
				mockedDB.On("GetOrCreateTFProvider", mock.Anything).Return(uuid.New(), false, nil)
				config.SetDBMocks(mockedDB)

				mockedCmdr := &mocks.Commander{}
				mockedCmdr.On("Run", mock.Anything).Return(func(cmd *exec.Cmd) error {
					assert.Contains(t, cmd.String(), "terraform version")
					return nil
				}).Once()
				mockedCmdr.On("Run", mock.Anything).Return(func(cmd *exec.Cmd) error {
					assert.Contains(t, cmd.String(), "terraform init")
					return nil
				}).Once()
				mockedCmdr.On("Run", mock.Anything).Return(func(cmd *exec.Cmd) error {
					assert.Contains(t, cmd.String(), "terraform providers schema -json")
					schema, _ := os.ReadFile("./testdata/example_schema.json")
					cmd.Stdout.Write(schema)
					return nil
				}).Once()

				commander.SetCommander(mockedCmdr)
			},
			Args:           []string{"-f", "./testdata/modules.yaml"},
			ValidateOutput: clitesting.ValidateOutputContains("Successfully added 0 Providers, 0 Resources, and 0 Attributes"),
		},
	})
}
