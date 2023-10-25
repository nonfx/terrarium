// Copyright (c) Ollion
// SPDX-License-Identifier: Apache-2.0

package testutils

import (
	"fmt"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func ValidateListOutput(t *testing.T, expected string, actual string, skipFields []int) (isValid bool) {
	expFields := splitListOutput(expected)
	actFields := splitListOutput(actual)
	sort.Ints(skipFields)
	skipIndex := 0
	isValid = assert.Equal(t, len(expFields), len(actFields), "List output mismatch: expected %d fields, actual %d fields. Actual ouptut is: %s", len(expFields), len(actFields), actual)
	if !isValid {
		return
	}
	for i, val := range expFields {
		if skipIndex < len(skipFields) && (i+1) == skipFields[skipIndex] {
			skipIndex++
		} else {
			isValid = assert.Equal(t, val, actFields[i], fmt.Sprintf("List output mismatch: expected '%s' - actual '%s' for field %d", val, actFields[i], i+1))
			if !isValid {
				return
			}
		}
	}
	return true
}

func splitListOutput(output string) []string {
	lastCharSpace := false
	splitStr := strings.FieldsFunc(output, func(c rune) bool {
		rv := (c == rune(' ') && lastCharSpace)
		lastCharSpace = (c == rune(' '))
		return rv
	})

	result := []string{}
	for _, val := range splitStr {
		v := strings.TrimSpace(val)
		if len(v) > 0 {
			result = append(result, v)
		}
	}
	return result
}
