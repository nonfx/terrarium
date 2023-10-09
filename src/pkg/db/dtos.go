// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package db

import (
	"github.com/cldcvr/terrarium/src/pkg/jsonschema"
	"github.com/google/uuid"
)

type DependencyResult struct {
	DependencyID uuid.UUID
	InterfaceID  string
	Name         string
	Schema       *jsonschema.Node
	Computed     bool
}
