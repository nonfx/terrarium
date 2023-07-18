package test

import (
	"strings"

	"github.com/cldcvr/api-automation-framework/assert"
	"github.com/cldcvr/terrarium/src/cli/int_test/constant"
	"github.com/cldcvr/terrarium/src/cli/int_test/helpers"
)

// Testify.suite() currently orders the test alphabetically,
// so I'm forcing the desired order by altering these three test case names.
func (suite *SmokeTestSuite) Test2FarmResources_Smoke() {
	outputResources, stderrResources, errResources := helpers.FarmResource(suite.T())
	assert.Nil(suite.T(), errResources, errResources)
	assert.NotNil(suite.T(), outputResources, "Verify farm resource command returns an output")
	assert.True(suite.T(), strings.Contains(outputResources, constant.SEED_RESOURCE) || strings.Contains(stderrResources, constant.SEED_RESOURCE) || strings.Contains(stderrResources, constant.SKIPPING_RESOURCE), "Verify farm resources command returns an output")
}

func (suite *SmokeTestSuite) Test3FarmModules_Smoke() {
	outputModules, stderrModules, errModules := helpers.FarmModules(suite.T())
	assert.Nil(suite.T(), errModules, errModules)
	assert.NotNil(suite.T(), outputModules, "Verify farm modules command returns an output")
	assert.True(suite.T(), strings.Contains(outputModules, constant.MODULE_RESOURCE_FOUND) || strings.Contains(stderrModules, constant.MODULE_RESOURCE_FOUND) || strings.Contains(stderrModules, constant.SKIPPING_RESOURCE), "Verify farm modules command returns an output")
}

func (suite *SmokeTestSuite) Test4FarmMappings_Smoke() {
	outputMappings, stderrMappings, errMappings := helpers.FarmMappings(suite.T())
	assert.Nil(suite.T(), errMappings, errMappings)
	assert.NotNil(suite.T(), outputMappings, "Verify farm mapings command returns an output")
	assert.True(suite.T(), strings.Contains(outputMappings, constant.SEED_RESOURCE) || strings.Contains(stderrMappings, constant.SEED_RESOURCE) || strings.Contains(stderrMappings, constant.MODULE_MAPPING), "Verify farm mappings command returns an output")
}
