// Copyright (c) Ollion
// SPDX-License-Identifier: Apache-2.0

package utils

import "golang.org/x/exp/slices"

// ToRefArr convert array elements to reference
func ToRefArr[T any](arr []T) (refArr []*T) {
	refArr = make([]*T, len(arr))
	for i := range arr {
		refArr[i] = &arr[i]
	}
	return
}

// TrimEmpty removes all empty elements from start and end of the array.
func TrimEmpty[T comparable](arr []T) []T {
	var empty T
	newArr := arr

	// Trim empty elements from the start
	for len(newArr) > 0 && newArr[0] == empty {
		newArr = newArr[1:]
	}

	// Trim empty elements from the end
	for len(newArr) > 0 && newArr[len(newArr)-1] == empty {
		newArr = newArr[:len(newArr)-1]
	}

	slices.Clip(newArr)
	return newArr
}

// ToIfaceArr converts given T type array to interface array.
func ToIfaceArr[T any](arr []T) []interface{} {
	ifaceArr := make([]interface{}, len(arr))
	for i := range arr {
		ifaceArr[i] = arr[i]
	}
	return ifaceArr
}
