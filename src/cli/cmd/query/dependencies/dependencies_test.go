//go:build mock
// +build mock

package dependencies

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	"github.com/cldcvr/terrarium/src/cli/internal/config"
	"github.com/cldcvr/terrarium/src/pkg/db"
	"github.com/cldcvr/terrarium/src/pkg/db/mocks"
	"github.com/cldcvr/terrarium/src/pkg/testutils/clitesting"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_CmdDependencies(t *testing.T) {
	config.LoadDefaults()
	clitest := clitesting.CLITest{
		CmdToTest: NewCmd,
	}

	mockDB := &mocks.DB{}
	mockDB.On("FetchDependencyByInterfaceID", mock.AnythingOfType("string")).Return(nil, fmt.Errorf("mock error")).Once()
	mockDB.On("FetchDependencyByInterfaceID", mock.AnythingOfType("string")).Return(&db.Dependency{
		InterfaceID: "mockID",
		Title:       "mockTitle",
		Description: "mockDescription",
		Inputs:      []byte("mockInputs"),
		Outputs:     []byte("mockOutputs"),
	}, nil)

	config.SetDBMocks(mockDB)

	tests := []clitesting.CLITestCase{
		{
			Name:    "db failure",
			WantErr: true,
		},
		{
			Name: "list dependencies in json format",
			PreExecute: func(ctx context.Context, t *testing.T, cmd *cobra.Command, cmdOpts clitesting.CmdOpts) {
				args := []string{"-o", "json"}
				cmd.SetArgs(args)
			},
			ValidateOutput: func(ctx context.Context, t *testing.T, cmdOpts clitesting.CmdOpts, output []byte) bool {
				expectedOutput := "{\n \"InterfaceID\": \"mockID\",\n \"Title\": \"mockTitle\",\n \"Description\": \"mockDescription\",\n \"Inputs\": \"mockInputs\",\n \"Outputs\": \"mockOutputs\"\n}\n"
				return assert.Equal(t, expectedOutput, string(output))
			},
		},
		{
			Name: "list dependencies in tabular format",
			PreExecute: func(ctx context.Context, t *testing.T, cmd *cobra.Command, cmdOpts clitesting.CmdOpts) {
				args := []string{"-o", "table"}
				cmd.SetArgs(args)
			},
			ValidateOutput: func(ctx context.Context, t *testing.T, cmdOpts clitesting.CmdOpts, output []byte) bool {
				buffer := new(bytes.Buffer)
				table := tablewriter.NewWriter(buffer)
				table.SetHeader([]string{"Interface ID", "Title", "Description", "Inputs", "Outputs"})
				table.Append([]string{"mockID", "mockTitle", "mockDescription", "mockInputs", "mockOutputs"})
				table.Render()
				return assert.Equal(t, buffer.String(), string(output))
			},
		},
	}

	clitest.RunTests(t, tests)
}
