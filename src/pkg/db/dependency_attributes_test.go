// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package db

import (
	"testing"

	"github.com/cldcvr/terrarium/src/pkg/jsonschema"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestToProtoSingleDependencyAttribute(t *testing.T) {
	// Given a DependencyAttribute with all fields filled
	name := "TestName"
	schemaDescription := "TestDescription"
	schemaType := "TestType"
	attr := DependencyAttribute{
		Name: &name,
		Schema: &jsonschema.Node{
			Description: schemaDescription,
			Type:        schemaType,
		},
	}

	// When ToProto is called
	protoResp := attr.ToProto()

	// Then the response should match the expected values
	assert.Equal(t, "TestName", protoResp.Title)
	assert.Equal(t, "TestDescription", protoResp.Description)
	assert.Equal(t, "TestType", protoResp.Type)
}

func TestToProtoMultipleDependencyAttributes(t *testing.T) {
	// Given a slice of DependencyAttributes
	name1 := "Name1"
	name2 := "Name2"
	desc2 := "Desc2"

	attrs := DependencyAttributes{
		&DependencyAttribute{Name: &name1},
		&DependencyAttribute{Name: &name2, Schema: &jsonschema.Node{Description: desc2}},
		// Add more if needed...
	}

	// When ToProto is called
	protoResps := attrs.ToProto()

	// Then the length of responses should match
	assert.Len(t, protoResps, 2)
}

func TestGetCondition(t *testing.T) {
	// Given a DependencyAttribute
	dependencyID := uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	name := "TestName"
	computed := true

	attr := &DependencyAttribute{
		DependencyID: &dependencyID,
		Name:         &name,
		Computed:     &computed,
	}

	// When GetCondition is called
	condition := attr.GetCondition()

	// Then the response should match the expected values
	expected := &DependencyAttribute{
		DependencyID: &dependencyID,
		Name:         &name,
		Computed:     &computed,
	}
	assert.Equal(t, expected, condition)
}
