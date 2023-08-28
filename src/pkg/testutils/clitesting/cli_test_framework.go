package clitesting

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"time"

	"github.com/Netflix/go-expect"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/h2non/gock.v1"

	"github.com/cldcvr/terrarium/src/pkg/localstate"
	"github.com/cldcvr/terrarium/src/pkg/testutils"
	"github.com/cldcvr/terrarium/src/pkg/utils"
)

const (
	KeyAuthToken         = "auth_token"
	KeyRefreshToken      = "refresh_token"
	KeyConnectedEndpoint = "connected_endpoint"
	APIHost              = "https://"

	failTestError = "some generic error"
)

var (
	ErrNotLoggedIn = errors.New("auth token isn't present in local cache - please run auth login")
)

type CmdOpts interface {
	SetTTY(tty *os.File)
}

var EmptyOpts CmdOpts

var failTestResponse = map[string]interface{}{
	"code":    16,
	"message": failTestError,
	"details": []string{},
}

type CLITest struct {
	// A function that will run once before any tests execute.
	SetupTest func(ctx context.Context, t *testing.T)
	// A function that will run once after all tests have run (optional)
	TeardownTest func(ctx context.Context, t *testing.T)
	// A function that will run before each test case (i.e. before each element in the []CLITestCase passed to RunTests)
	SetupTestCase func(ctx context.Context, t *testing.T, tc CLITestCase)
	// A function that will run after each test case
	TeardownTestCase func(ctx context.Context, t *testing.T, tc CLITestCase)
	// A Context object that can be used to maintain context between tests/test cases
	Ctx context.Context
	// newCmd function to test
	CmdToTest  *cobra.Command
	ParentCmds []func() *cobra.Command
	// function to initialize the command specific CmdOptions structure
	CmdOptionsInit func() CmdOpts
	// This is a function that will run before the command's Execute() for every test (i.e. a global PreExecute)
	PreExecute func(ctx context.Context, t *testing.T, tc CLITestCase, cmd *cobra.Command, cmdOpts CmdOpts)
	// should the Unauthenticated test be automatically added to the run
	AddUnauthenticatedTest bool
	// Args to set for unauthenticated test
	UnauthTestArgs []string
	// should the API Failure test be automatically added to the run
	AddFailTest      bool
	FailTestMethod   string
	FailTestEndpoint string
	FailTestPath     string
	FailTestArgs     []string
	DumpMocks        bool

	cmdOptions        CmdOpts
	testStateFileName string
	console           *expect.Console
	doneConsole       chan struct{}
	outputBuffer      *bytes.Buffer
}

// CLITestCase is a structure used to define the specific test cases to be run via RunTests.
type CLITestCase struct {
	// Description of the test case
	Name string
	// Arguments to pass while executing the command for this test case
	Args []string
	// Whether to inject a token into the localstate before starting the test
	TokenPresent    bool
	UseInvalidToken bool
	UseExpiredToken bool
	// Setup pseudo-tty for go-expect testing and call this function for test
	PseudoTTY func(ctx context.Context, t *testing.T, console *expect.Console)
	// Default timeout for go-expect responses in seconds
	ExpectTimeout int64
	// Whether an error is expected
	WantErr bool
	// String to compare error to
	ExpError string
	// A function that will run before the test and will setup Gock for the test
	GockSetup func(ctx context.Context, t *testing.T)
	// A function that will run before the command's Execute() method is called - a good place to call SetArgs
	PreExecute func(ctx context.Context, t *testing.T, cmd *cobra.Command, cmdOpts CmdOpts)
	// A function that will run to validate the output of the test
	ValidateOutput      ValidateOutputFunc
	ValidateAuditFields bool
}

type ValidateOutputFunc func(ctx context.Context, t *testing.T, cmdOpts CmdOpts, output []byte) bool

func (clitest *CLITest) setupTest(t *testing.T) func(t *testing.T) {
	var err error

	// Override the default config file so that we don't overwrite or pick up any config
	// present in the local dev environment
	viper.SetConfigFile("/tmp/some_unknown_file")

	// Override the default state file so we don't overwrite any existing state
	// present in the local dev environment
	clitest.testStateFileName, err = testutils.GetTempFileName("/tmp", "clitest", "yaml")
	if err != nil {
		t.Error(err)
	}

	if clitest.SetupTest != nil {
		clitest.SetupTest(clitest.Ctx, t)
	}

	return func(t *testing.T) {
		if clitest.TeardownTest != nil {
			clitest.TeardownTest(clitest.Ctx, t)
		}
	}
}

func (clitest *CLITest) setupTestCase(t *testing.T, tc CLITestCase) func(t *testing.T, tc CLITestCase) {
	if clitest.CmdOptionsInit != nil {
		clitest.cmdOptions = clitest.CmdOptionsInit()
	}
	if clitest.cmdOptions == nil {
		clitest.cmdOptions = EmptyOpts
	}

	localstate.SetStateFileName(clitest.testStateFileName)
	clitest.outputBuffer = new(bytes.Buffer)

	if tc.TokenPresent {
		localstate.Set(KeyConnectedEndpoint, "api.endpoint")
		if tc.UseInvalidToken {
			localstate.Set(KeyAuthToken, "sometoken")
		} else {
			if tc.UseExpiredToken {
				localstate.Set(KeyAuthToken, testutils.GetExpiredToken())
			} else {
				localstate.Set(KeyAuthToken, testutils.GetTestToken())
			}
			localstate.Set(KeyRefreshToken, testutils.GetTestRefreshToken())
		}
	}
	if tc.GockSetup != nil {
		tc.GockSetup(clitest.Ctx, t)
		if clitest.DumpMocks {
			fmt.Printf(">>> Mock dump: %s <<<\n", tc.Name)
			for _, g := range gock.GetAll() {
				fmt.Printf("   %s: %s -> %d\n", g.Request().Method, g.Request().URLStruct, g.Request().Response.StatusCode)
			}
		}

	}
	if tc.PseudoTTY != nil {
		var err error
		// create a fake TTY for the test that prompts for password
		var consoleOpts = []expect.ConsoleOpt{expect.WithStdout(clitest.outputBuffer)}
		if tc.ExpectTimeout > 0 {
			consoleOpts = append(consoleOpts, expect.WithDefaultTimeout(time.Duration(tc.ExpectTimeout)*time.Second))
		}
		clitest.console, err = utils.NewVT10XConsole(consoleOpts...)
		require.Nil(t, err)
		clitest.cmdOptions.SetTTY(clitest.console.Tty())
	}

	if clitest.SetupTestCase != nil {
		clitest.SetupTestCase(clitest.Ctx, t, tc)
	}
	return func(t *testing.T, tc CLITestCase) {
		if clitest.TeardownTestCase != nil {
			clitest.TeardownTestCase(clitest.Ctx, t, tc)
		}
		if clitest.console != nil {
			clitest.console.Close()
		}
		os.Remove(clitest.testStateFileName)
	}
}

func (clitest *CLITest) RunTests(t *testing.T, testCases []CLITestCase) {
	teardownTest := clitest.setupTest(t)
	defer teardownTest(t)
	defer gock.Off()

	testCases = addAutoTests(clitest, testCases)

	for _, tt := range testCases {
		teardownTestCase := clitest.setupTestCase(t, tt)
		t.Run(tt.Name, func(t *testing.T) {
			var cmd *cobra.Command
			if clitest.ParentCmds != nil {
				subCmd := clitest.CmdToTest
				var nextCmd *cobra.Command
				for i, cmdFunc := range clitest.ParentCmds {
					if i == 0 {
						cmd = cmdFunc()
						nextCmd = cmd
					} else {
						newCmd := cmdFunc()
						nextCmd.AddCommand(newCmd)
						nextCmd = newCmd
					}
				}
				nextCmd.AddCommand(subCmd)
			} else {
				cmd = clitest.CmdToTest
			}

			cmd.SetArgs(tt.Args)

			if tt.PseudoTTY != nil {
				cmd.SetOut(clitest.console.Tty())
				cmd.SetErr(clitest.console.Tty())
				cmd.SetIn(clitest.console.Tty())
				clitest.doneConsole = make(chan struct{})
				go func() {
					defer close(clitest.doneConsole)
					tt.PseudoTTY(clitest.Ctx, t, clitest.console)
				}()
			} else {
				clitest.outputBuffer.Reset()
				cmd.SetOut(clitest.outputBuffer)
				cmd.SetErr(clitest.outputBuffer)
			}
			if clitest.PreExecute != nil {
				clitest.PreExecute(clitest.Ctx, t, tt, cmd, clitest.cmdOptions)
			}
			if tt.PreExecute != nil {
				tt.PreExecute(clitest.Ctx, t, cmd, clitest.cmdOptions)
			}

			err := cmd.Execute()
			if (err == nil) == tt.WantErr {
				t.Errorf("%s: error = %+v, wantErr %+v", tt.Name, err, tt.WantErr)
			} else {
				if tt.WantErr {
					assert.Contains(t, err.Error(), tt.ExpError, fmt.Sprintf("%+v", err))
					assert.True(t, gock.IsDone(), "Not all gocks were called")
				}

				if tt.PseudoTTY != nil {
					clitest.console.Tty().Close()
					<-clitest.doneConsole
				}
				out, err := io.ReadAll(clitest.outputBuffer)
				if err != nil {
					t.Errorf("%s: error occurred processing stdout: %s", tt.Name, err)
				}
				if tt.ValidateOutput != nil {
					if tt.ValidateAuditFields {
						assert.Contains(t, string(out), "createdAt")
						assert.Contains(t, string(out), "updatedAt")
						assert.Contains(t, string(out), "updatedBy")
					}
					assert.True(t, tt.ValidateOutput(clitest.Ctx, t, clitest.cmdOptions, out))
				}

				assert.True(t, gock.IsDone(), "Not all gocks were called")
			}
		})
		teardownTestCase(t, tt)
	}
}

func addAutoTests(clitest *CLITest, tcIn []CLITestCase) (tcOut []CLITestCase) {
	tcOut = tcIn
	if clitest.AddUnauthenticatedTest {
		tcOut = append(tcOut, CLITestCase{
			Name:         "Not authenticated",
			TokenPresent: false,
			WantErr:      true,
			ExpError:     ErrNotLoggedIn.Error(),
			PreExecute: func(ctx context.Context, t *testing.T, cmd *cobra.Command, cmdOpts CmdOpts) {
				if len(clitest.UnauthTestArgs) > 0 {
					cmd.SetArgs(clitest.UnauthTestArgs)
				}
			},
		})
	}
	if clitest.AddFailTest {
		tcOut = append(tcOut, CLITestCase{
			Name:         "Failure from API",
			TokenPresent: true,
			GockSetup: func(ctx context.Context, t *testing.T) {
				gockReq := gock.New(APIHost)
				gockReq.Method = strings.ToUpper(clitest.FailTestMethod)
				gockReq.Path(clitest.FailTestEndpoint + clitest.FailTestPath).
					Reply(http.StatusInternalServerError).
					JSON(failTestResponse)
			},
			PreExecute: func(ctx context.Context, t *testing.T, cmd *cobra.Command, cmdOpts CmdOpts) {
				if len(clitest.FailTestArgs) > 0 {
					cmd.SetArgs(clitest.FailTestArgs)
				}
			},
			WantErr:  true,
			ExpError: failTestError,
		})
	}
	return
}

// ValidateOutputContains helper to assert if the output contains the given string
func ValidateOutputContains(expectedString string) ValidateOutputFunc {
	return validateOutputAsserter(expectedString, false)
}

// ValidateOutputContains helper to assert if the output is exactly same as the given string
func ValidateOutputMatch(expectedString string) ValidateOutputFunc {
	return validateOutputAsserter(expectedString, true)
}

func validateOutputAsserter(expectedString string, exactMatch bool) ValidateOutputFunc {
	return func(ctx context.Context, t *testing.T, cmdOpts CmdOpts, output []byte) bool {
		if exactMatch {
			return assert.Equal(t, expectedString, string(output))
		}

		return assert.Contains(t, string(output), expectedString)
	}
}
