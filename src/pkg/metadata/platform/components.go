package platform

import (
	"sort"
	"strings"

	"github.com/cldcvr/terraform-config-inspect/tfconfig"
	"github.com/cldcvr/terrarium/src/pkg/jsonschema"
	"github.com/xeipuuv/gojsonschema"
	"github.com/zclconf/go-cty/cty"
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
	for k := range platformModule.ModuleCalls {
		if !strings.HasPrefix(k, ComponentPrefix) {
			continue
		}

		id := strings.TrimPrefix(k, ComponentPrefix)
		c := cArr.GetByID(id)
		if c == nil {
			c = cArr.Append(Component{ID: id})
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

	cv, _ := v.Value(nil)
	valMap := cv.AsValueMap()
	if v, ok := valMap["default"]; ok {
		extractSchema(c.Inputs, v)
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
		if c.Outputs.Properties[outputKey] == nil {
			c.Outputs.Properties[outputKey] = &jsonschema.Node{
				Description: v.Description,
			}
		}
	}
}

func extractSchema(existingSchema *jsonschema.Node, value cty.Value) {
	switch value.Type().FriendlyName() {
	case "object":
		mapValue := value.AsValueMap()

		existingSchema.Type = gojsonschema.TYPE_OBJECT
		if existingSchema.Properties == nil {
			existingSchema.Properties = map[string]*jsonschema.Node{}
		}

		for key, val := range mapValue {
			existingSchema.Properties[key] = &jsonschema.Node{}
			extractSchema(existingSchema.Properties[key], val)
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
			extractSchema(existingSchema.Items, listVal[0])
		}
	}
}
