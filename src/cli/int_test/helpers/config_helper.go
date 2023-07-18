package helpers

import (
	"os"

	"github.com/cldcvr/api-automation-framework/logger"
)

var (
	AutomationPath string // Worspace path
	Env            string // Environment on which all test case will be triggered
	TestSuiteName  string
	ReportPath     string
)

func init() {
	ReportPath = os.Getenv("REPORT_PATH")

	if AutomationPath = os.Getenv("AUTOMATION_PATH"); AutomationPath == "" && ReportPath != "" {
		logger.LogError("AUTOMATION_PATH environment variable must be set.")
	}

	if TestSuiteName = os.Getenv("REPORT_TITLE"); TestSuiteName == "" {
		TestSuiteName = "Terrarium CLI Test Result"
	}

	if Env = os.Getenv("ENVIRONMENT"); Env == "" {
		Env = "Dev"
	}
}
