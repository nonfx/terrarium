package platform

import (
	"fmt"
	"sort"
	"strings"

	"github.com/cldcvr/terraform-config-inspect/tfconfig"
	"golang.org/x/exp/constraints"
)

func NewBlockID(t BlockType, bKey string) BlockID {
	switch t {
	case BlockType_Data:
		return BlockID(bKey)
	default:
		return BlockID(fmt.Sprintf("%s.%s", t, bKey))
	}
}

// GetBlockType returns a BlockType from pre-defined set of constants.
// This would be similar to typecast except, it changes unrecognized values to BlockType_Undefined
func GetBlockType(t string) BlockType {
	bt := BlockType(t)
	switch bt {
	case BlockType_ModuleCall,
		BlockType_Resource,
		BlockType_Data,
		BlockType_Local,
		BlockType_Variable,
		BlockType_Output,
		BlockType_Provider:

		return bt
	default:
		return BlockType_Undefined
	}
}

func (bID BlockID) Parse() (t BlockType, bKey string) {
	spl := strings.SplitN(string(bID), ".", 2)
	if len(spl) < 2 {
		return
	}

	t = GetBlockType(spl[0])

	switch t {
	case BlockType_Data:
		bKey = string(bID)
	default:
		bKey = spl[1]
	}
	return
}

func (b BlockID) ParseComponent() (bt BlockType, componentName string) {
	bt, bn := b.Parse()
	if !strings.HasPrefix(bn, ComponentPrefix) {
		return
	}

	switch bt {
	case BlockType_ModuleCall, BlockType_Local:
		return bt, strings.TrimPrefix(bn, ComponentPrefix)
	}

	return bt, ""
}

func (bID BlockID) GetBlock(m *tfconfig.Module) (b ParsedBlock, found bool) {
	bt, bn := bID.Parse()
	switch bt {
	case BlockType_ModuleCall:
		b, found = m.ModuleCalls[bn]
		return

	case BlockType_Resource:
		b, found = m.ManagedResources[bn]
		return

	case BlockType_Data:
		b, found = m.DataResources[bn]
		return

	case BlockType_Local:
		b, found = m.Locals[bn]
		return

	case BlockType_Variable:
		b, found = m.Variables[bn]
		return

	case BlockType_Output:
		b, found = m.Outputs[bn]
		return

	case BlockType_Provider:
		b, found = m.RequiredProviders[bn]
		return
	}

	return
}

func (bID BlockID) FindRequirements(m *tfconfig.Module) (requirements []BlockID) {
	requirements = []BlockID{}

	b, found := bID.GetBlock(m)
	if !found || b == nil {
		return
	}

	dg, ok := b.(BlockDependencyGetter)
	if !ok || dg == nil {
		return
	}

	for _, v := range dg.GetDependencies() {
		bId, found := getBlockIDFromTFAttribute(v, m)
		if found {
			requirements = appendSortedUnique(requirements, bId)
		}
	}

	pg, ok := b.(BlockProviderGetter)
	if !ok || pg.GeProviderName() == "" {
		return
	}

	pbID := NewBlockID(BlockType_Provider, pg.GeProviderName())
	requirements = appendSortedUnique(requirements, pbID)

	return
}

func getBlockIDFromTFAttribute(v tfconfig.AttributeReference, m *tfconfig.Module) (BlockID, bool) {
	switch v.Type() {
	case "":
		return "", false
	case "module":
		_, found := m.ModuleCalls[v.Name()]
		return NewBlockID(BlockType_ModuleCall, v.Name()), found

	case "local":
		_, found := m.Locals[v.Name()]
		return NewBlockID(BlockType_Local, v.Name()), found

	case "var":
		_, found := m.Variables[v.Name()]
		return NewBlockID(BlockType_Variable, v.Name()), found
	}

	resName := fmt.Sprintf("%s.%s", v.Type(), v.Name())
	if _, found := m.ManagedResources[resName]; found {
		return NewBlockID(BlockType_Resource, resName), found
	}

	dataName := "data." + resName
	if _, found := m.DataResources[dataName]; found {
		return NewBlockID(BlockType_Data, dataName), found
	}

	return "", false
}

func appendSortedUnique[T constraints.Ordered](arr []T, val T) []T {
	idx := sort.Search(len(arr), func(i int) bool {
		return arr[i] >= val
	})

	// return original array if value already exists
	if idx < len(arr) && arr[idx] == val {
		return arr
	}

	// Insert the value at the correct position in the sorted array
	arr = append(arr[:idx], append([]T{val}, arr[idx:]...)...)

	return arr
}
