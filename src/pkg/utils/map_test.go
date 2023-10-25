// Copyright (c) Ollion
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"testing"

	"github.com/rotisserie/eris"
	"github.com/stretchr/testify/assert"
)

func TestGetKeys(t *testing.T) {
	t.Run("EmptyMap", func(t *testing.T) {
		// Test when input map is empty.
		m := make(map[int]string)
		keys := GetKeys(m)
		assert.Empty(t, keys, "Resulting keys slice should be empty")
	})

	t.Run("IntStringMap", func(t *testing.T) {
		// Test when input map contains integers as keys and strings as values.
		m := map[int]string{
			1: "one",
			2: "two",
			3: "three",
		}
		keys := GetKeys(m)
		assert.Len(t, keys, len(m), "Length of resulting keys slice should match input map")

		for k := range m {
			assert.Contains(t, keys, k, "Keys slice should contain all keys from the input map")
		}
	})

	t.Run("StringIntMap", func(t *testing.T) {
		// Test when input map contains strings as keys and integers as values.
		m := map[string]int{
			"one":   1,
			"two":   2,
			"three": 3,
		}
		keys := GetKeys(m)
		assert.Len(t, keys, len(m), "Length of resulting keys slice should match input map")

		for k := range m {
			assert.Contains(t, keys, k, "Keys slice should contain all keys from the input map")
		}
	})
}

func TestMapEachSortedKeys(t *testing.T) {
	t.Run("EmptyMap", func(t *testing.T) {
		// Test when input map is empty.
		m := make(map[int]string)
		var processedKeys []int
		err := MapEachSortedKeys(m, func(k int, v string) error {
			processedKeys = append(processedKeys, k)
			return nil
		})
		assert.NoError(t, err, "No error should be returned")
		assert.Empty(t, processedKeys, "No keys should be processed")
	})

	t.Run("IntStringMap", func(t *testing.T) {
		// Test when input map contains integers as keys and strings as values.
		m := map[int]string{
			3: "three",
			1: "one",
			2: "two",
		}
		var processedKeys []int
		err := MapEachSortedKeys(m, func(k int, v string) error {
			processedKeys = append(processedKeys, k)
			return nil
		})
		assert.NoError(t, err, "No error should be returned")
		assert.Equal(t, []int{1, 2, 3}, processedKeys, "Keys should be processed in sorted order")
	})

	t.Run("ErrorCallback", func(t *testing.T) {
		// Test when the callback function returns an error.
		m := map[int]string{
			1: "one",
			2: "two",
		}
		err := MapEachSortedKeys(m, func(k int, v string) error {
			if k == 2 {
				return eris.New("mocked error")
			}
			return nil
		})
		assert.Error(t, err)
	})
}
