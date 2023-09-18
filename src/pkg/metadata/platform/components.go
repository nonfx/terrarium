// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package platform

import (
	"io"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/cldcvr/terraform-config-inspect/tfconfig"
	"github.com/cldcvr/terrarium/src/pkg/jsonschema"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/icza/backscanner"
	"github.com/xeipuuv/gojsonschema"
	"github.com/zclconf/go-cty/cty"
	"golang.org/x/exp/slices"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const (
	docCommentTitleArgTag = "title"
	docCommentDescArgTag  = "description"
	docCommentEnumArgTag  = "enum"
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
		setValueFromDocIfFound(&c.Title, docCommentTitleArgTag, docs)
		setValueFromDocIfFound(&c.Description, docCommentDescArgTag, docs)

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
	head, args, _, _ := parseBlockDocComment(blockPos)

	if _, ok := args[docCommentDescArgTag]; !ok && head != nil { // if there is no description tag use the head lines instead
		args[docCommentDescArgTag] = strings.Join(head, "\n")
	}

	return
}

// parseBlockDocComment parses docummentation comment to sections:
//
// # Quo in officia nobis autem pariatur sit tenetur ut dolores.		<--- head line
// # Deleniti asperiores quaerat.										<--- head line
// # @title: Incidunt aperiam sit facilis.								<--- argument tag
// # Voluptatem officiis aperiam.										<--- reminder
//
// It returns head (lines before the first argument tag with comment symbol removed), args (argument tags grouped by tag name), lines (list of all lines).
func parseBlockDocComment(blockPos tfconfig.SourcePos) (head []string, args map[string]string, lines []string, err error) {
	if blockPos.Filename == "" {
		return
	}
	lines, err = readCommentLinesAbove(blockPos.Filename, blockPos.StartByte)
	if err != nil {
		return
	}
	slices.Reverse(lines) // lines are read bottom-up

	head = make([]string, 0, len(lines))
	args = make(map[string]string, len(lines))
	for _, line := range lines {
		if groups := docCommentArgMatcher.FindStringSubmatch(line); groups != nil {
			args[groups[1]] = groups[2]
		}
		if len(args) < 1 { // all lines before the first argument tag form the comment head
			head = append(head, strings.TrimSpace(strings.TrimPrefix(line, "#")))
		}
	}

	return
}

// read all comment lines (ignoring empty lines) above a given end byte until the first non-comment or the end of file
func readCommentLinesAbove(filename string, endPos int) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := backscanner.New(file, endPos)
	for {
		line, _, err := scanner.Line()
		if err == io.EOF {
			return lines, nil
		} else if err != nil {
			return lines, err
		}

		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == "" {
			continue
		} else if docCommentMatcher.MatchString(trimmedLine) {
			lines = append(lines, trimmedLine)
			continue
		}

		return lines, nil
	}
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

func setListFromDocIfFound(values *[]interface{}, valueTagName string, fieldDoc map[string]string) {
	var enumCSV string
	if ok := setValueFromDocIfFound(&enumCSV, valueTagName, fieldDoc); ok {
		enumValues := strings.Split(enumCSV, ",")
		*values = make([]interface{}, 0, len(enumValues))
		for _, value := range enumValues {
			*values = append(*values, strings.TrimSpace(value))
		}
	}
}

func setValueFromDocIfFound(value *string, valueTagName string, fieldDoc map[string]string) (ok bool) {
	if newValue, exists := fieldDoc[valueTagName]; exists && newValue != "" {
		*value = newValue
		ok = true
	}
	return
}

func extractSchema(existingSchema *jsonschema.Node, value cty.Value, fieldName string, fieldDocs map[string]map[string]string) {
	existingSchema.Title = tfValueToTitle(fieldName, nil) // default to field name in title format
	if fieldDoc, ok := fieldDocs[fieldName]; ok {
		setValueFromDocIfFound(&existingSchema.Title, docCommentTitleArgTag, fieldDoc)
		setValueFromDocIfFound(&existingSchema.Description, docCommentDescArgTag, fieldDoc)
		setListFromDocIfFound(&existingSchema.Enum, docCommentEnumArgTag, fieldDoc)
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
