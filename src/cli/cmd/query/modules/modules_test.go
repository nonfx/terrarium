package modules

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/cldcvr/terrarium/src/cli/internal/config"
	"github.com/cldcvr/terrarium/src/pkg/db"
	"github.com/cldcvr/terrarium/src/pkg/db/mocks"
	"github.com/cldcvr/terrarium/src/pkg/testutils/clitesting"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_CmdModules(t *testing.T) {
	config.LoadDefaults()
	os.Setenv("TR_DB_TYPE", "mock")
	clitest := clitesting.CLITest{
		CmdToTest: GetCmd(),
	}
	mockUuid1 := uuid.New()
	configdb, _ := config.DBConnect()
	mockDB := configdb.(*mocks.DB)
	mockDB.On("QueryTFModules", mock.AnythingOfType("db.FilterOption"), mock.AnythingOfType("db.FilterOption"), mock.AnythingOfType("db.FilterOption"), mock.AnythingOfType("db.FilterOption")).
		Return(nil, fmt.Errorf("mock error")).Once()
	mockDB.On("QueryTFModules", mock.AnythingOfType("db.FilterOption"), mock.AnythingOfType("db.FilterOption"), mock.AnythingOfType("db.FilterOption"), mock.AnythingOfType("db.FilterOption")).
		Return(db.TFModules{
			{
				Model:       db.Model{ID: mockUuid1},
				ModuleName:  "Rds",
				Version:     "1",
				Source:      "/Users/xyz/abc/tf-dir",
				Description: "",
				Namespace:   "farm_repo",
			},
		}, nil)
	tests := []clitesting.CLITestCase{
		{
			Name:    "db failure",
			WantErr: true,
		},
		{
			Name: "list local modules with namespace in json format",
			PreExecute: func(ctx context.Context, t *testing.T, cmd *cobra.Command, cmdOpts clitesting.CmdOpts) {

				args := []string{"-o", "json", "-n", "/Users/xyz/abc/tf-dir"}

				cmd.SetArgs(args)
			},
			ValidateOutput: func(ctx context.Context, t *testing.T, cmdOpts clitesting.CmdOpts, output []byte) bool {
				expectedOutput := "{\n \"modules\": [\n  {\n   \"id\": \"" + mockUuid1.String() + "\",\n   \"taxonomy_id\": \"00000000-0000-0000-0000-000000000000\",\n   \"module_name\": \"Rds\",\n   \"source\": \"/Users/xyz/abc/tf-dir\",\n   \"version\": \"1\",\n   \"namespace\": \"farm_repo\"\n  }\n ],\n \"page\": {\n  \"size\": 100\n }\n}"
				return assert.Equal(t, expectedOutput, string(output))
			},
		},
		{
			Name: "list local modules with namespace in tabular format",
			PreExecute: func(ctx context.Context, t *testing.T, cmd *cobra.Command, cmdOpts clitesting.CmdOpts) {
				args := []string{"-o", "table", "-n", "/Users/xyz/abc/tf-dir"}
				cmd.SetArgs(args)
			},
			ValidateOutput: func(ctx context.Context, t *testing.T, cmdOpts clitesting.CmdOpts, output []byte) bool {
				expectedOutput := "  ID                                    MODULE NAME  SOURCE                 VERSION  NAMESPACE  \n  " + mockUuid1.String() + "  Rds          /Users/xyz/abc/tf-dir  1        farm_repo  \n"
				return assert.Equal(t, expectedOutput, string(output))
			},
		},
	}
	clitest.RunTests(t, tests)
}