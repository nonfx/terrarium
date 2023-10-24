// Copyright (c) Ollion
// SPDX-License-Identifier: Apache-2.0

package jsonschema

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/cldcvr/terrarium/src/pkg/pb/terrariumpb"
	"github.com/rotisserie/eris"
	"github.com/xeipuuv/gojsonschema"
)

type Node struct {
	Title            string           `yaml:"title,omitempty" json:"title,omitempty"`
	Description      string           `yaml:"description,omitempty" json:"description,omitempty"`
	Type             string           `yaml:"type,omitempty" json:"type,omitempty"` // string, number, integer, boolean, object, array, null
	Default          interface{}      `yaml:"default,omitempty" json:"default,omitempty"`
	Examples         []interface{}    `yaml:"examples,omitempty" json:"examples,omitempty"`
	Enum             []interface{}    `yaml:"enum,omitempty" json:"enum,omitempty"`
	MinLength        int32            `yaml:"minLength,omitempty" json:"minLength,omitempty"`
	MaxLength        int32            `yaml:"maxLength,omitempty" json:"maxLength,omitempty"`
	Pattern          string           `yaml:"pattern,omitempty" json:"pattern,omitempty"`
	Format           string           `yaml:"format,omitempty" json:"format,omitempty"`
	Minimum          int32            `yaml:"minimum,omitempty" json:"minimum,omitempty"`
	Maximum          int32            `yaml:"maximum,omitempty" json:"maximum,omitempty"`
	ExclusiveMinimum int32            `yaml:"exclusiveMinimum,omitempty" json:"exclusiveMinimum,omitempty"`
	ExclusiveMaximum int32            `yaml:"exclusiveMaximum,omitempty" json:"exclusiveMaximum,omitempty"`
	MultipleOf       int32            `yaml:"multipleOf,omitempty" json:"multipleOf,omitempty"`
	Items            *Node            `yaml:"items,omitempty" json:"items,omitempty"`
	AdditionalItems  bool             `yaml:"additionalItems,omitempty" json:"additionalItems,omitempty"`
	MinItems         int32            `yaml:"minItems,omitempty" json:"minItems,omitempty"`
	MaxItems         int32            `yaml:"maxItems,omitempty" json:"maxItems,omitempty"`
	UniqueItems      bool             `yaml:"uniqueItems,omitempty" json:"uniqueItems,omitempty"`
	Properties       map[string]*Node `yaml:"properties,omitempty" json:"properties,omitempty"`
	Required         []string         `yaml:"required,omitempty" json:"required,omitempty"`

	compiled *gojsonschema.Schema
}

func (n *Node) Compile() (err error) {
	n.compiled, err = gojsonschema.NewSchema(gojsonschema.NewGoLoader(n))
	if err != nil {
		return eris.Wrapf(err, "failed to process json schema")
	}
	return
}

func (n *Node) Validate(val interface{}) error {
	err := n.compileIfNot()
	if err != nil {
		return eris.Wrapf(err, "failed to compile validation schema")
	}

	result, err := n.compiled.Validate(gojsonschema.NewGoLoader(val))
	if err != nil {
		return eris.Wrapf(err, "failed to process the input object")
	}

	if !result.Valid() {
		return eris.Errorf("validation failed with following errors: \n%s", formatErrors(result.Errors()))
	}

	return nil
}

func (n *Node) ApplyDefaultsToMSI(inp map[string]interface{}) {
	if inp == nil {
		return
	}

	for k, v := range n.Properties {
		if _, isSet := inp[k]; isSet {
			continue
		}

		inp[k] = v.Default
	}
}

func (n *Node) ApplyDefaultsToArr(inp []interface{}) {
	if n.Items == nil {
		return
	}

	for i, v := range inp {
		if v != nil {
			continue
		}

		inp[i] = n.Items.Default
	}
}

func (n *Node) compileIfNot() error {
	if n.compiled == nil {
		return n.Compile()
	}
	return nil
}

func formatErrors(errs []gojsonschema.ResultError) string {
	s := []string{}
	for _, e := range errs {
		s = append(s, e.String())
	}

	return fmt.Sprintf("\t%s", strings.Join(s, "\n\t"))
}

// Implement the sql.Scanner interface to take care of unmarshaling
// the serialized form (stored in the database) into the Go Node structure
func (n *Node) Scan(value interface{}) error {
	var data []byte

	switch v := value.(type) {
	case string:
		data = []byte(v)
	case []byte:
		data = v
	default:
		return errors.New("type assertion to string or []byte failed")
	}

	return json.Unmarshal(data, n)
}

// Implement the driver.Valuer interface to serialize the Node struct
//
//	into a format suitable for storing in the database.
func (n Node) Value() (driver.Value, error) {
	return json.Marshal(n)
}

func (n *Node) ToProto() *terrariumpb.JSONSchema {
	if n == nil {
		return nil
	}

	// Create the base proto representation
	protoSchema := &terrariumpb.JSONSchema{
		Title:       n.Title,
		Description: n.Description,
		Type:        n.Type,
	}

	// If properties exist OR the type is an "object",
	// then we can convert each property
	if n.Properties != nil || n.Type == "object" {
		protoSchema.Properties = make(map[string]*terrariumpb.JSONSchema)

		for key, prop := range n.Properties {
			protoSchema.Properties[key] = prop.ToProto()
		}
	}

	return protoSchema
}
