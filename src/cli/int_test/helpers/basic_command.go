package helpers

import (
	"fmt"
	"testing"

	"github.com/cldcvr/api-automation-framework/report"
	"github.com/rendon/testcli"
)

func GetVersion(t *testing.T) (output, stdErr string, err error) {
	c := testcli.Command(t, COMMAND_NAME, "version")
	c.Run()
	c.Stdout()
	report.SetCliReqRes(FARM_MODULE, COMMAND_NAME, c.Stdout(), c.Stderr(), "version")
	if !c.Success() {
		return c.Stdout(), c.Stderr(), fmt.Errorf("error version commmand %v %v", c.Stdout(), c.Stderr())
	}
	return c.Stdout(), c.Stderr(), nil
}
