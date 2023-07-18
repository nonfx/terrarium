package test

import (
	"strings"

	"github.com/cldcvr/api-automation-framework/assert"
	"github.com/cldcvr/terrarium/src/cli/int_test/helpers"
)

func (suite *SmokeTestSuite) Test1GetVersion_Smoke() {
	out, _, err := helpers.GetVersion(suite.T())
	assert.Nil(suite.T(), err, err)
	assert.NotNil(suite.T(), out, "Verify version command returns an output")
	assert.True(suite.T(), strings.Contains(out, "terrarium"))
}
