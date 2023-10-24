// Copyright (c) Ollion
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

func TestTrimEmpty(t *testing.T) {
	tests := []struct {
		name string
		inp  []string
		outp []string
	}{
		{
			name: "nil input",
			inp:  nil,
			outp: nil,
		},
		{
			name: "empty input",
			inp:  []string{},
			outp: []string{},
		},
		{
			name: "one empty",
			inp:  []string{""},
			outp: []string{},
		},
		{
			name: "array with no empty",
			inp:  []string{"a", "b"},
			outp: []string{"a", "b"},
		},
		{
			name: "array with no empty at ends",
			inp:  []string{"a", "", "b"},
			outp: []string{"a", "", "b"},
		},
		{
			name: "array with empty at ends",
			inp:  []string{"", "", "a", "", "b", "", ""},
			outp: []string{"a", "", "b"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOutp := TrimEmpty(tt.inp)
			assert.Equal(t, tt.outp, gotOutp)
		})
	}
}

func TestToIfaceArr(t *testing.T) {
	tests := []struct {
		name string
		inp  []string
		outp []interface{}
	}{
		{
			name: "nil",
			inp:  nil,
			outp: []interface{}{},
		},
		{
			name: "empty",
			inp:  []string{},
			outp: []interface{}{},
		},
		{
			name: "regular array",
			inp:  []string{"a", "b", "c"},
			outp: []interface{}{"a", "b", "c"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOutp := ToIfaceArr(tt.inp)
			assert.Equal(t, tt.outp, gotOutp)
		})
	}
}
