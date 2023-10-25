// Copyright (c) Ollion
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"github.com/rotisserie/eris"
	"github.com/zclconf/go-cty/cty"
)

func ToCtyValue(v interface{}) (out cty.Value, err error) {
	switch v := v.(type) {
	case string:
		return cty.StringVal(v), err
	case int:
		return cty.NumberIntVal(int64(v)), err
	case float64:
		return cty.NumberFloatVal(v), err
	case bool:
		return cty.BoolVal(v), err
	case map[string]interface{}:
		data := make(map[string]cty.Value)
		for k, val := range v {
			data[k], err = ToCtyValue(val)
			if err != nil {
				return cty.NilVal, err
			}
		}
		return cty.ObjectVal(data), err
	case []interface{}:
		var values []cty.Value
		for _, val := range v {
			var value cty.Value
			value, err = ToCtyValue(val)
			if err != nil {
				return cty.NilVal, err
			}
			values = append(values, value)
		}
		return cty.TupleVal(values), err
	default:
		return cty.NilVal, eris.Errorf("unsupported type: %T", v)
	}
}
