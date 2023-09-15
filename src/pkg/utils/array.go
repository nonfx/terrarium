// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package utils

// ToRefArr convert array elements to reference
func ToRefArr[T any](arr []T) (refArr []*T) {
	refArr = make([]*T, len(arr))
	for i := range arr {
		refArr[i] = &arr[i]
	}
	return
}
