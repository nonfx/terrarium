// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToRefArr(t *testing.T) {
	t.Run("EmptyArray", func(t *testing.T) {
		// Test when input array is empty.
		arr := []string{}
		refArr := ToRefArr(arr)
		assert.Empty(t, refArr, "Resulting reference array should be empty")
	})

	t.Run("StringArray", func(t *testing.T) {
		// Test when input array contains strings.
		arr := []string{"apple", "banana", "cherry"}
		refArr := ToRefArr(arr)
		assert.Len(t, refArr, len(arr))

		for i := 0; i < len(arr); i++ {
			assert.Equal(t, &arr[i], refArr[i])
		}
	})

	t.Run("IntArray", func(t *testing.T) {
		// Test when input array contains integers.
		arr := []int{1, 2, 3}
		refArr := ToRefArr(arr)
		assert.Len(t, refArr, len(arr))

		for i := 0; i < len(arr); i++ {
			assert.Equal(t, &arr[i], refArr[i])
		}
	})

	t.Run("MixedTypeArray", func(t *testing.T) {
		// Test when input array contains mixed types (string and int).
		arr := []interface{}{"apple", 2, "cherry"}
		refArr := ToRefArr(arr)
		assert.Len(t, refArr, len(arr))

		for i := 0; i < len(arr); i++ {
			assert.Equal(t, &arr[i], refArr[i])
		}
	})
}
