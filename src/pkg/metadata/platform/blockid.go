package platform

import (
	"fmt"
	"sort"
	"strings"

	"github.com/cldcvr/terraform-config-inspect/tfconfig"
	"golang.org/x/exp/constraints"
)

func NewBlockID(t BlockType, name string) BlockID {
	return BlockID(fmt.Sprintf("%s.%s", t, name))
}

func NewBlockType(t string) BlockType {
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

func (bid BlockID) Parse() (t BlockType, name string) {
	spl := strings.SplitN(string(bid), ".", 2)
	if len(spl) < 2 {
		return
	}

	t = NewBlockType(spl[0])
	name = spl[1]
	return
}

func (b BlockID) IsComponent() bool {
	bt, bn := b.Parse()
	return bt == BlockType_ModuleCall && strings.HasPrefix(bn, ComponentPrefix)
}

func (b BlockID) GetComponentName() string {
	if !b.IsComponent() {
		return ""
	}

	_, bn := b.Parse()
	return strings.TrimPrefix(bn, ComponentPrefix)
}

func (b BlockID) FindRequirements(m *tfconfig.Module) []BlockID {
	t, _ := b.Parse()
	switch t {
	case BlockType_ModuleCall:
		return b.findModuleCallReq(m)
	case BlockType_Resource:
		return b.findResourceReq(m)
	case BlockType_Data:
		return b.findResourceDataReq(m)
	case BlockType_Output:
		return b.findOutputReq(m)
	}

	return []BlockID{}
}

func (b BlockID) findModuleCallReq(m *tfconfig.Module) []BlockID {
	_, name := b.Parse()
	mc := m.ModuleCalls[name]
	if mc == nil {
		return []BlockID{}
	}

	return b.findReqBlockIDs(mc.Dependencies, m)
}

func (b BlockID) findResourceReq(m *tfconfig.Module) []BlockID {
	_, name := b.Parse()
	mr := m.ManagedResources[name]
	if mr == nil {
		return []BlockID{}
	}

	return b.findReqBlockIDs(mr.Dependencies, m)
}

func (b BlockID) findResourceDataReq(m *tfconfig.Module) []BlockID {
	_, name := b.Parse()
	dr := m.DataResources[name]
	if dr == nil {
		return []BlockID{}
	}

	return b.findReqBlockIDs(dr.Dependencies, m)
}

func (b BlockID) findOutputReq(m *tfconfig.Module) []BlockID {
	_, name := b.Parse()
	op := m.Outputs[name]
	if op == nil {
		return []BlockID{}
	}

	return b.findReqBlockIDs(op.Dependencies, m)
}

func (b BlockID) findReqBlockIDs(dependencies map[string]tfconfig.AttributeReference, m *tfconfig.Module) []BlockID {
	requirements := []BlockID{}

	for _, v := range dependencies {
		bId, found := getBlockIDFromTFAttribute(v, m)
		if found {
			requirements = appendSortedUnique(requirements, bId)
		}
	}

	return requirements
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
