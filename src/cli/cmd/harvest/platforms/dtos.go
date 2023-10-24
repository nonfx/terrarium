// Copyright (c) Ollion
// SPDX-License-Identifier: Apache-2.0

package platforms

import (
	"github.com/cldcvr/terrarium/src/pkg/jsonschema"
	"github.com/xeipuuv/gojsonschema"
)

type Terrarium struct {
	// DependencyInterfaces holds the list of all interfaces defined in the file.
	Terrarium Components `yaml:"components,omitempty"`
}

type Components []Component

type Component struct {
	ID          string           `yaml:"id"`
	Title       string           `yaml:"title"`
	Description string           `yaml:"description"`
	Inputs      *jsonschema.Node `yaml:"inputs"`
	Outputs     *jsonschema.Node `yaml:"outputs"`
}

func (i *Component) Init() {
	if i.Inputs != nil {
		i.Inputs.Type = gojsonschema.TYPE_OBJECT
	}
	if i.Outputs != nil {
		i.Outputs.Type = gojsonschema.TYPE_OBJECT
	}
}

func (iArr Components) Init() {
	for i := range iArr {
		iArr[i].Init()
	}
}

func (f Terrarium) Init() {
	if f.Terrarium != nil {
		f.Terrarium.Init()
	}
}

type Platform struct {
	Title       string     `yaml:"title"`
	Description string     `yaml:"description"`
	RepoURL     string     `yaml:"repo_url"`
	RepoDir     string     `yaml:"repo_directory"`
	Revisions   []Revision `yaml:"revisions"`
}

type Revision struct {
	Label string `yaml:"label"`
	Type  string `yaml:"type"`
}
