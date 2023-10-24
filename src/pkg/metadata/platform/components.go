// Copyright (c) Ollion
// SPDX-License-Identifier: Apache-2.0

package platform

import (
	"regexp"
	"sort"
	"strings"

	"github.com/cldcvr/terraform-config-inspect/tfconfig"
	"github.com/cldcvr/terrarium/src/pkg/jsonschema"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/xeipuuv/gojsonschema"
	"github.com/zclconf/go-cty/cty"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	docCommentMatcher    = regexp.MustCompile("^\\s*#.+$")
	docCommentArgMatcher = regexp.MustCompile("^#\\s*@(.+?)\\s*:\\s*(.+)$")
)

func NewComponents(platformModule *tfconfig.Module) Components {
	cArr := Components{}
	cArr.Parse(platformModule)

	return cArr
}

func (cArr Components) GetByID(id string) *Component {
	for i, v := range cArr {
		if v.ID == id {
			return &cArr[i]
		}
	}

	return nil
}

func (cArr *Components) Append(c Component) *Component {
	(*cArr) = append((*cArr), c)
	return &(*cArr)[len(*cArr)-1]
}

func (cArr *Components) Parse(platformModule *tfconfig.Module) {
	for k, mc := range platformModule.ModuleCalls {
		if !strings.HasPrefix(k, ComponentPrefix) {
			continue
		}

		id := strings.TrimPrefix(k, ComponentPrefix)
		c := cArr.GetByID(id)
		if c == nil {
			c = cArr.Append(Component{ID: id})
		}
		docs := getBlockDoc(mc.Pos)

		c.Title = tfValueToTitle(id, nil) // default to component name in title format
		SetValueFromDocIfFound(&c.Title, docCommentTitleArgTag, docs)
		SetValueFromDocIfFound(&c.Description, docCommentDescArgTag, docs)

		c.fetchInputs(platformModule)
		c.fetchOutputs(platformModule)
	}

	sort.Slice(*cArr, func(i, j int) bool {
		return (*cArr)[i].ID < (*cArr)[j].ID
	})
}

func (c *Component) fetchInputs(m *tfconfig.Module) {
	varName := ComponentPrefix + c.ID
	v := m.Locals[varName]
	if v == nil {
		return
	}

	if c.Inputs == nil {
		c.Inputs = &jsonschema.Node{}
	}

	fieldDoc := getLocalInputBlockDocs(v)

	cv, _ := v.Expression.Value(nil)
	valMap := cv.AsValueMap()
	if mv, ok := valMap["default"]; ok {
		extractSchema(c.Inputs, mv, "", fieldDoc)
	}
}

func getLocalInputBlockDocs(block *tfconfig.Local) (docsByKey map[string]map[string]string) {
	docsByKey = make(map[string]map[string]string)
	switch localExpr := block.Expression.(type) {
	case *hclsyntax.ObjectConsExpr:
		for _, kv := range localExpr.ExprMap() { // default input block
			switch inputExpr := kv.Value.(type) {
			case *hclsyntax.ObjectConsExpr:
				for _, input := range inputExpr.ExprMap() { // individual input values
					inputName, _ := input.Key.Value(nil) // input name value
					inputRange := input.Key.Range()      // input name start position (i.e. read doc comment from here upwards)
					docsByKey[inputName.AsString()] = getBlockDoc(tfconfig.SourcePos{
						Filename:  block.ParentPos.Filename,
						Line:      inputRange.Start.Line,
						StartByte: inputRange.Start.Byte,
						EndLine:   inputRange.End.Line,
						EndByte:   inputRange.End.Byte,
					})
				}
			}
		}
	}
	return
}

func getBlockDoc(blockPos tfconfig.SourcePos) (args map[string]string) {
	args, _ = GetDoc(blockPos.Filename, blockPos.StartByte, true)
	return
}

func (c *Component) fetchOutputs(m *tfconfig.Module) {
	prefix := ComponentPrefix + c.ID + "_"
	if c.Outputs == nil {
		c.Outputs = &jsonschema.Node{
			Type: gojsonschema.TYPE_OBJECT,
		}
	}

	if c.Outputs.Properties == nil {
		c.Outputs.Properties = map[string]*jsonschema.Node{}
	}

	for k, v := range m.Outputs {
		if !strings.HasPrefix(k, prefix) {
			continue
		}

		outputKey := strings.TrimPrefix(k, prefix)
		node, ok := c.Outputs.Properties[outputKey]
		if !ok {
			node = &jsonschema.Node{}
			c.Outputs.Properties[outputKey] = node
		}
		node.Title = tfValueToTitle(outputKey, &prefix)
		node.Description = v.Description
	}
}

// tfValueToTitle removes the common component prefix and converts snake-case "default_db_name" to "Default Db Name"
func tfValueToTitle(value string, prefix *string) string {
	if prefix == nil {
		c := ComponentPrefix
		prefix = &c
	}
	return cases.Title(language.Und, cases.NoLower).String(strings.ReplaceAll(strings.TrimPrefix(value, *prefix), "_", " "))
}

func extractSchema(existingSchema *jsonschema.Node, value cty.Value, fieldName string, fieldDocs map[string]map[string]string) {
	existingSchema.Title = tfValueToTitle(fieldName, nil) // default to field name in title format
	if fieldDoc, ok := fieldDocs[fieldName]; ok {
		SetValueFromDocIfFound(&existingSchema.Title, docCommentTitleArgTag, fieldDoc)
		SetValueFromDocIfFound(&existingSchema.Description, docCommentDescArgTag, fieldDoc)
		SetListFromDocIfFound(&existingSchema.Enum, docCommentEnumArgTag, fieldDoc)
	}

	switch value.Type().FriendlyName() {
	case "object":
		mapValue := value.AsValueMap()

		existingSchema.Type = gojsonschema.TYPE_OBJECT
		if existingSchema.Properties == nil {
			existingSchema.Properties = map[string]*jsonschema.Node{}
		}

		for key, val := range mapValue {
			existingSchema.Properties[key] = &jsonschema.Node{}
			extractSchema(existingSchema.Properties[key], val, key, fieldDocs)
		}
	case "string":
		existingSchema.Type = gojsonschema.TYPE_STRING
		existingSchema.Default = value.AsString()
		return
	case "number":
		existingSchema.Type = gojsonschema.TYPE_NUMBER
		existingSchema.Default, _ = value.AsBigFloat().Float64()
		return
	case "bool":
		existingSchema.Type = gojsonschema.TYPE_BOOLEAN
		existingSchema.Default = value.True()
		return
	case "tuple":
		listVal := value.AsValueSlice()
		existingSchema.Type = gojsonschema.TYPE_ARRAY
		if existingSchema.Items == nil {
			existingSchema.Items = &jsonschema.Node{}
		}

		if len(listVal) > 0 {
			extractSchema(existingSchema.Items, listVal[0], fieldName, fieldDocs)
		}
	}
}
