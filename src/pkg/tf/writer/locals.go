// Copyright (c) Ollion
// SPDX-License-Identifier: Apache-2.0

package writer

import (
	"io"
	"sort"

	"github.com/cldcvr/terrarium/src/pkg/utils"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/rotisserie/eris"
)

func WriteLocals(val map[string]interface{}, out io.Writer) error {
	// Convert the map to cty.Value
	ctyData, err := utils.ToCtyValue(val)
	if err != nil {
		return eris.Wrapf(err, "error converting given value to hcl: %v", val)
	}
	data := ctyData.AsValueMap()

	// Create a new HCL file
	f := hclwrite.NewEmptyFile()
	rootBody := f.Body()

	// Create a new block for "locals"
	localsBlock := rootBody.AppendNewBlock("locals", nil)

	// Sort the keys for deterministic output
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Set the local variables in the HCL block in sorted order
	for _, k := range keys {
		localsBlock.Body().SetAttributeValue(k, data[k])
	}

	// Write the HCL content to the provided io.Writer
	_, err = f.WriteTo(out)
	if err != nil {
		return eris.Wrap(err, "error writing hcl to the output buffer")
	}

	return nil
}
