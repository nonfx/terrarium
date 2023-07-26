package platform

import (
	"fmt"
	"sort"
	"strings"

	"github.com/cldcvr/terraform-config-inspect/tfconfig"
)

func (bid BlockID) Parse() (t BlockType, name string) {
	spl := strings.SplitN(string(bid), ".", 2)
	if len(spl) < 2 {
		return
	}

	t = NewBlockType(spl[0])
	name = spl[1]
	return
}

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

func NewGraph(platformModule *tfconfig.Module) Graph {
	g := Graph{}
	g.Parse(platformModule)
	return g
}

func (g Graph) GetByID(id BlockID) *GraphNode {
	for i, v := range g {
		if v.ID == id {
			return &g[i]
		}
	}

	return nil
}

func (g *Graph) Append(bID BlockID, requirements []BlockID) *GraphNode {
	(*g) = append((*g), GraphNode{ID: bID, Requirements: requirements})
	return &(*g)[len(*g)-1]
}

func (g *Graph) Parse(srcModule *tfconfig.Module) {
	toTraverse := map[BlockID]struct{}{}

	for k := range srcModule.ModuleCalls {
		if !strings.HasPrefix(k, ComponentPrefix) {
			continue
		}
		bID := NewBlockID(BlockType_ModuleCall, k)
		toTraverse[bID] = struct{}{}
	}

	for k := range srcModule.Outputs {
		bID := NewBlockID(BlockType_Output, k)
		toTraverse[bID] = struct{}{}
	}

	for len(toTraverse) > 0 {
		for bID := range toTraverse {
			if g.GetByID(bID) != nil {
				delete(toTraverse, bID)
				continue
			}

			blockRequirements := bID.FindRequirements(srcModule)
			g.Append(bID, blockRequirements)
			delete(toTraverse, bID)

			for _, reqBId := range blockRequirements {
				toTraverse[reqBId] = struct{}{}
			}
		}
	}

	sort.Slice(*g, func(i, j int) bool {
		return (*g)[i].ID < (*g)[j].ID
	})
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

	return b.findInputReq(mc.Inputs, m)
}

func (b BlockID) findResourceReq(m *tfconfig.Module) []BlockID {
	_, name := b.Parse()
	mr := m.ManagedResources[name]
	if mr == nil {
		return []BlockID{}
	}

	return b.findInputReq(mr.Inputs, m)

}

func (b BlockID) findResourceDataReq(m *tfconfig.Module) []BlockID {
	_, name := b.Parse()
	dr := m.DataResources[name]
	if dr == nil {
		return []BlockID{}
	}

	return b.findInputReq(dr.Inputs, m)
}

func (b BlockID) findOutputReq(m *tfconfig.Module) []BlockID {
	_, name := b.Parse()
	op := m.Outputs[name]
	if op == nil {
		return []BlockID{}
	}

	return b.findInputReq(map[string]tfconfig.AttributeReference{"": op.Value}, m)
}

func (b BlockID) findInputReq(inputs map[string]tfconfig.AttributeReference, m *tfconfig.Module) []BlockID {
	requirements := []BlockID{}

	for _, v := range inputs {
		switch v.Type() {
		case "":
		case "module":
			if _, ok := m.ModuleCalls[v.Name()]; ok {
				requirements = append(requirements, NewBlockID(BlockType_ModuleCall, v.Name()))
			}
		case "local":
			if _, ok := m.Locals[v.Name()]; ok {
				requirements = append(requirements, NewBlockID(BlockType_Local, v.Name()))
			}
		case "var":
			if _, ok := m.Variables[v.Name()]; ok {
				requirements = append(requirements, NewBlockID(BlockType_Variable, v.Name()))
			}
		default:
			resName := fmt.Sprintf("%s.%s", v.Type(), v.Name())
			if _, ok := m.ManagedResources[resName]; ok {
				requirements = append(requirements, NewBlockID(BlockType_Resource, resName))
			}

			dataName := "data." + resName
			if _, ok := m.DataResources[dataName]; ok {
				requirements = append(requirements, NewBlockID(BlockType_Data, dataName))
			}
		}
	}

	sort.Slice(requirements, func(i, j int) bool {
		return requirements[i] < requirements[j]
	})

	return requirements
}
