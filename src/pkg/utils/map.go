// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
)

func GetKeys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	return keys
}

func MapEachSortedKeys[K constraints.Ordered, V any](m map[K]V, fu func(K, V) error) (err error) {
	keys := GetKeys(m)
	slices.Sort(keys)
	for _, k := range keys {
		err = fu(k, m[k])
		if err != nil {
			return err
		}
	}

	return
}
