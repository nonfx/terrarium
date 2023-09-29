// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package taxonomy

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplit(t *testing.T) {
	testCases := []struct {
		name     string
		taxon    Taxon
		expected []string
	}{
		{
			name:     "Single level taxonomy",
			taxon:    "storage",
			expected: []string{"storage"},
		},
		{
			name:     "Multiple level taxonomy",
			taxon:    "storage/database/rdbms/postgres",
			expected: []string{"storage", "database", "rdbms", "postgres"},
		},
		{
			name:     "Empty taxonomy",
			taxon:    "",
			expected: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.taxon.Split()
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestNewTaxonomy(t *testing.T) {
	testCases := []struct {
		name     string
		levels   []string
		expected Taxon
	}{
		{
			name:     "Single level taxonomy",
			levels:   []string{"storage"},
			expected: "storage",
		},
		{
			name:     "Multiple level taxonomy",
			levels:   []string{"storage", "database", "rdbms", "postgres"},
			expected: "storage/database/rdbms/postgres",
		},
		{
			name:     "Empty taxonomy",
			levels:   []string{},
			expected: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := NewTaxonomy(tc.levels...)
			assert.Equal(t, tc.expected, result)
		})
	}
}
