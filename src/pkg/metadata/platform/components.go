package platform

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/cldcvr/terraform-config-inspect/tfconfig"
	"github.com/cldcvr/terrarium/src/pkg/jsonschema"
	"github.com/icza/backscanner"
	"github.com/xeipuuv/gojsonschema"
	"github.com/zclconf/go-cty/cty"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	docCommentMatcher    = regexp.MustCompile("^\\s*#.+$")
	docCommentArgMatcher = regexp.MustCompile("#\\s*(.+?)(\\s*\\[(.+)\\])?:\\s*(.+)")
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
		docs, _, err := parseBlockDocComment(mc.Module, mc.Pos)
		if err != nil {
			return
		}

		c.Title = tfValueToTitle(id, nil) // default to component name in title format
		if docs, ok := docs["component"]; ok {
			docTitle := docs[0]
			if docTitle != "" {
				c.Title = docTitle
			}
			c.Description = docs[1]
		}

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

	fieldDoc, _, err := parseBlockDocComment(m, v.Pos)
	if err != nil {
		return
	}

	cv, _ := v.Expression.Value(nil)
	valMap := cv.AsValueMap()
	if mv, ok := valMap["default"]; ok {
		extractSchema(c.Inputs, mv, "", fieldDoc)
	}
}

func parseBlockDocComment(m *tfconfig.Module, blockPos tfconfig.SourcePos) (args map[string][2]string, lines []string, err error) {
	if blockPos.Filename == "" {
		return
	}
	lines, err = readCommentLinesAbove(blockPos.Filename, blockPos.StartByte)
	if err != nil {
		return
	}

	args = make(map[string][2]string, len(lines))
	for _, line := range lines {
		if groups := docCommentArgMatcher.FindStringSubmatch(line); groups != nil {
			args[groups[1]] = [2]string{groups[3], groups[4]}
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

func extractSchema(existingSchema *jsonschema.Node, value cty.Value, fieldName string, fieldDoc map[string][2]string) {
	existingSchema.Title = tfValueToTitle(fieldName, nil) // default to field name in title format
	if docString, ok := fieldDoc[fieldName]; ok {
		docTitle := docString[0]
		if docTitle != "" {
			existingSchema.Title = docTitle
		}
		existingSchema.Description = docString[1]
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
			valueFieldName := strings.TrimPrefix(fmt.Sprintf("%s.%s", fieldName, key), ".")
			extractSchema(existingSchema.Properties[key], val, valueFieldName, fieldDoc)
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
			extractSchema(existingSchema.Items, listVal[0], fieldName, fieldDoc)
		}
	}
}
