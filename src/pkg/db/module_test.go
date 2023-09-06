// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersion_Compare(t *testing.T) {
	testCases := []struct {
		v1       Version
		v2       Version
		expected int
	}{
		{Version("1.0"), Version("1.0"), 0},
		{Version("1.0"), Version("1.1"), -1},
		{Version("1.1"), Version("1.0"), 1},
		{Version("1.0.0"), Version("1.0"), 0},
		{Version("1.0"), Version("1.0.0"), 0},
		{Version("1.0.1"), Version("1.0.2"), -1},
		{Version("1.0.2"), Version("1.0.1"), 1},
		{Version("1.10.2"), Version("1.2.5"), 1},
		{Version("2.0.0"), Version("1.9.9"), 1},
		{Version("1.9.9"), Version("2.0.0"), -1},
	}

	for _, testCase := range testCases {
		result := testCase.v1.Compare(testCase.v2)
		assert.Equal(t, testCase.expected, result, "Version(%s).Compare(Version(%s))", testCase.v1, testCase.v2)
	}
}
