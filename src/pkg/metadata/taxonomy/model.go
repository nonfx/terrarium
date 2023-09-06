// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package taxonomy

import "strings"

// Taxon represents a single taxonomic unit. It is a string type that can represent multiple levels of a taxonomy,
// separated by a `/`. For example, a Taxon could be "storage/database/rdbms/postgres", representing a hierarchical
// structure of categories.
type Taxon string

// separator is a constant that defines the character used to separate different levels within a Taxon.
const separator = "/"

// Split method of a Taxon splits the Taxon into its constituent levels based on the separator and returns them as a slice of strings.
// For example, if the Taxon is "storage/database/rdbms/postgres", Split would return ["storage", "database", "rdbms", "postgres"].
func (t Taxon) Split() (levels []string) {
	return strings.Split(string(t), separator)
}

// NewTaxonomy function takes a variadic parameter of strings representing levels and joins them into a single Taxon using the separator.
// This function is used to create a new Taxon. For example, NewTaxonomy("storage", "database", "rdbms", "postgres") would return a Taxon
// "storage/database/rdbms/postgres".
func NewTaxonomy(levels ...string) Taxon {
	return Taxon(strings.Join(levels, separator))
}
