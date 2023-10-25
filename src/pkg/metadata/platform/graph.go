// Copyright (c) Ollion
// SPDX-License-Identifier: Apache-2.0

package platform

import (
	"sort"
	"strings"

	"github.com/cldcvr/terraform-config-inspect/tfconfig"
	"github.com/rotisserie/eris"
	"golang.org/x/exp/slices"
)

func NewGraph(platformModule *tfconfig.Module) Graph {
	g := Graph{}
	g.Parse(platformModule)
	return g
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

type GraphWalkerCB func(blockId BlockID) error

// Walk a function to traverse all requirements starting from the given nodes
// and call cb exactly once for each node that is connected to the given set of nodes.
// Terraform Outputs are traversed differently in the end, such that, each
// output that is able to resolve with the blocks been traversed, are selected.
func (g *Graph) Walk(roots []BlockID, cb GraphWalkerCB) error {
	roots = slices.Compact(roots)
	traverser := make([]BlockID, len(roots)) // nodes before `i` are visited and after `i` are queued
	copy(traverser, roots)

	err := g.traverseRootBlocks(&traverser, cb)
	if err != nil {
		return eris.Wrap(err, "error traversing hcl blocks")
	}

	err = g.traverseOutputBlocks(&traverser, cb)
	if err != nil {
		return eris.Wrap(err, "error traversing hcl output blocks")
	}

	return nil
}

// traverse all requirements starting from the given nodes
func (g *Graph) traverseRootBlocks(traverser *[]BlockID, cb GraphWalkerCB) error {
	for i := 0; i < len(*traverser); i++ {
		node := g.GetByID((*traverser)[i])
		if node == nil {
			continue
		}

		err := cb(node.ID)
		if err != nil {
			return err
		}

		g.appendRequirements(traverser, node.Requirements)
	}

	return nil
}

func (g *Graph) appendRequirements(traverser *[]BlockID, requirements []BlockID) {
	for _, bID := range requirements {
		if !slices.Contains(*traverser, bID) {
			*traverser = append(*traverser, bID)
		}
	}
}

// traverse outputs whose requirements are already traversed.
func (g *Graph) traverseOutputBlocks(traverser *[]BlockID, cb GraphWalkerCB) error {
	for _, node := range *g {
		bt, _ := node.ID.Parse()
		if bt != BlockType_Output {
			continue
		}

		if !g.allDependenciesTraversed(traverser, node.Requirements) {
			continue
		}

		err := cb(node.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *Graph) allDependenciesTraversed(traverser *[]BlockID, requirements []BlockID) bool {
	for _, bId := range requirements {
		if !slices.Contains(*traverser, bId) {
			return false
		}
	}

	return true
}
