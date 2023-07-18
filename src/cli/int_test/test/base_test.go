package test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/cldcvr/api-automation-framework/logger"
	customReport "github.com/cldcvr/api-automation-framework/report"
	"github.com/cldcvr/terrarium/src/cli/int_test/helpers"
	"github.com/stretchr/testify/suite"
)

var report *customReport.Report

type SmokeTestSuite struct {
	suite.Suite
}

func (suite *SmokeTestSuite) SetupSuite() {
	// Initialize the report
	report = customReport.NewReport(helpers.TestSuiteName, helpers.Env, os.Getenv("GIT_REF"), os.Getenv("AZURE_BUILD_URL"))
}

func (suite *SmokeTestSuite) BeforeTest(suiteName, testName string) {
	logger.LogStep(fmt.Sprintf("%s Start", testName))
}

func (suite *SmokeTestSuite) AfterTest(suiteName, testName string) {
	logger.LogStep(fmt.Sprintf("%s End", testName))
	if !(suite.T().Failed() || suite.T().Skipped()) {
		report.SetModulePass()
	}
}

func (suite *SmokeTestSuite) TearDownSuite() {
	logger.LogStep(fmt.Sprintf("TearDownSuite Start"))

	saveReportPath := filepath.Join(helpers.ReportPath, "CLIReport.html")
	logger.LogInfo("Report path : ", saveReportPath)
	if os.Getenv("SLACK_NOTIFICATION") == "true" {

		// Create and Upload report in gcs
		url, err := report.UploadReport(saveReportPath)
		if err != nil {
			logger.LogStep("Report upload to google cloud is failed.")
		}
		logger.LogStep("Report link: ", url)

		// Send slack notification
		customReport.SendSlackMessage(suite.T(), report.GetReport(), url)
	} else {
		// Create the report and store it in the root directory
		report.CreateReport(saveReportPath)
	}
}

func TestSmokeSuite(t *testing.T) {
	suite.Run(t, new(SmokeTestSuite))
}
