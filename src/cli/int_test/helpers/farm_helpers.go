package helpers

import (
	"fmt"
	"testing"

	"github.com/cldcvr/api-automation-framework/report"
	"github.com/rendon/testcli"
)

const (
	FARM_MODULE  = "farm"
	COMMAND_NAME = "terrarium"
)

func FarmResource(t *testing.T) (output, stdErr string, err error) {
	c := testcli.Command(t, COMMAND_NAME, FARM_MODULE, "resources")
	c.Run()
	c.Stdout()
	report.SetCliReqRes(FARM_MODULE, COMMAND_NAME, c.Stdout(), c.Stderr(), FARM_MODULE, "resources")
	if !c.Success() {
		return c.Stdout(), c.Stderr(), fmt.Errorf("error farm resource commmand %v %v", c.Stdout(), c.Stderr())
	}
	return c.Stdout(), c.Stderr(), nil
}

func FarmModules(t *testing.T) (output, stdErr string, err error) {
	c := testcli.Command(t, COMMAND_NAME, FARM_MODULE, "modules")
	c.Run()
	c.Stdout()
	report.SetCliReqRes(FARM_MODULE, COMMAND_NAME, c.Stdout(), c.Stderr(), FARM_MODULE, "modules")
	if !c.Success() {
		return c.Stdout(), c.Stderr(), fmt.Errorf("error farm modules commmand %v %v", c.Stdout(), c.Stderr())
	}
	return c.Stdout(), c.Stderr(), nil
}

func FarmMappings(t *testing.T) (output, stdErr string, err error) {
	c := testcli.Command(t, COMMAND_NAME, FARM_MODULE, "mappings")
	c.Run()
	c.Stdout()
	report.SetCliReqRes(FARM_MODULE, COMMAND_NAME, c.Stdout(), c.Stderr(), FARM_MODULE, "mappings")
	if !c.Success() {
		return c.Stdout(), c.Stderr(), fmt.Errorf("error farm mappings commmand %v %v", c.Stdout(), c.Stderr())
	}

	return c.Stdout(), c.Stderr(), nil
}
