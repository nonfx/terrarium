package jsonschema

import "github.com/xeipuuv/gojsonschema"

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
