// Copyright (c) Ollion
// SPDX-License-Identifier: Apache-2.0

package localstate

import (
	"os"
	"testing"

	"github.com/cldcvr/terrarium/src/pkg/testutils"
	"github.com/stretchr/testify/assert"
)

func TestLocalState(t *testing.T) {
	testStateFileName, err := testutils.GetTempFileName("/tmp", "lstest", "yaml")
	if err != nil {
		t.Error(err)
	}
	t.Cleanup(func() { os.Remove(testStateFileName) })

	SetStateFileName(testStateFileName)
	Set("key1", "value1")
	assert.Equal(t, "value1", Get("key1"))
	Set("key2", "value2")
	assert.Equal(t, "value2", Get("key2"))
	Unset("key1")
	assert.Equal(t, "", Get("key1"))
	Clear()
	assert.Equal(t, "", Get("key2"))
}
