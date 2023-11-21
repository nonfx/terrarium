// Copyright (c) Ollion
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/cldcvr/terrarium/src/pkg/pb/terrariumpb"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/structpb"
	"gopkg.in/yaml.v3"
)

// App multiple apps configuration
type Apps []App

// App represents the main application configuration.
type App struct {
	// ID is a required identifier for the app in the project, which must start with
	// an alphabet character, can only contain alphanumeric characters,
	// and must not be longer than 20 characters.
	ID string `yaml:"id"`

	// Name describes the human-friendly name for the application.
	Name string `yaml:"name"`

	// EnvPrefix is the prefix used for the environment variables in this app.
	// If not set, defaults to an empty string.
	EnvPrefix string `yaml:"env_prefix"`

	// Compute denotes a specific dependency that best classifies the app itself,
	// it can only be of the type `compute/*`.
	// id of this dependency is automatically set to app id.
	// it is used to setup deployment pipeline in Code Pipes and allow other
	// apps to use this app as dependency.
	Compute Dependency `yaml:"compute"`

	// Dependencies lists the required services, databases, and other components that the application relies on.
	Dependencies Dependencies `yaml:"dependencies"`
}

func (e App) WrapProtoMessage() (message *anypb.Any, err error) {
	platformArtifact, err := e.ProtoValue()
	if err != nil {
		return nil, fmt.Errorf("invalid platform-artifact data: %v", err)
	}
	message, err = anypb.New(platformArtifact)
	if err != nil {
		return nil, fmt.Errorf("could not create platform-artifact message: %v", err)
	}
	return
}

func (e App) ProtoValue() (*terrariumpb.App, error) {
	compute, err := e.Compute.ProtoValue()
	if err != nil {
		return nil, fmt.Errorf("invalid compute data: %v", err)
	}
	dependencies, err := e.Dependencies.ProtoValue()
	if err != nil {
		return nil, fmt.Errorf("invalid dependencies data: %v", err)
	}

	return &terrariumpb.App{
		Id:           e.ID,
		Name:         e.Name,
		EnvPrefix:    e.EnvPrefix,
		Compute:      compute,
		Dependencies: dependencies,
	}, nil
}

func (e *App) ScanProto(m *terrariumpb.App) {
	if m != nil {
		e.ID = m.Id
		e.Name = m.Name
		e.EnvPrefix = m.EnvPrefix
		e.Compute.ScanProto(m.Compute)
		e.Dependencies.ScanProto(m.Dependencies)
	}
}

func (e App) Value() (driver.Value, error) {
	return json.Marshal(e)
}

func (e *App) Scan(value interface{}) error {
	data, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("invalid data")
	}
	return json.Unmarshal(data, e)
}

func (e App) ToFileBytes() ([]byte, error) {
	return yaml.Marshal(e)
}

// IsEquivalent returns if this and the other app generate the same infrastructure.
func (e App) IsEquivalent(other App) (isEquivalent bool) {
	if isEquivalent = e.Compute.IsEquivalent(other.Compute); !isEquivalent {
		return
	}

	if isEquivalent = len(e.Dependencies) == len(other.Dependencies); !isEquivalent {
		return
	}

	for _, thisItem := range e.Dependencies {
		for _, otherItem := range other.Dependencies {
			if isEquivalent = thisItem.IsEquivalent(otherItem); isEquivalent {
				break
			}
		}

		if !isEquivalent {
			return
		}
	}

	return
}

type Dependencies []Dependency

func (e Dependencies) ProtoValue() (values []*terrariumpb.AppDependency, err error) {
	values = make([]*terrariumpb.AppDependency, len(e))
	for i := range e {
		values[i], err = e[i].ProtoValue()
		if err != nil {
			return
		}
	}
	return
}

func (e *Dependencies) ScanProto(m []*terrariumpb.AppDependency) {
	*e = make(Dependencies, len(m))
	for i := range m {
		(*e)[i].ScanProto(m[i])
	}
}

// Dependency represents a single dependency of the application,
// which could be a database, another service, cache, etc.
type Dependency struct {
	// ID is a required identifier for the dependency in the project, which must start with
	// an alphabet character, can only contain alphanumeric characters,
	// and must not be longer than 20 characters.
	ID string `yaml:"id"`

	// Use indicates the specific dependency interface ID that is used to provision an app dependency.
	Use string `yaml:"use"`

	// EnvPrefix is used to prefix the output env vars in order to avoid collision
	// Defaults to dependency id upper case.
	EnvPrefix string `yaml:"env_prefix"`

	// Inputs represents customization options for the selected dependency interface.
	Inputs map[string]interface{} `yaml:"inputs"`

	// Outputs maps environment variables to dependency outputs. Keys are app env name (without prefix) and
	// values are Mustache templates using dependency outputs.
	// The default env var name format is `<app_env_prefix>_<dependency_env_prefix>_<dependency_output_name_to_upper>`.
	Outputs map[string]string `yaml:"outputs"`

	// NoProvision indicates whether the dependency is provisioned in another app.
	// If true, this dependency is shared and its inputs are set in another app
	// and its outputs are made available here.
	NoProvision bool `yaml:"no_provision"`
}

func (e Dependency) ProtoValue() (*terrariumpb.AppDependency, error) {
	inputs, err := structpb.NewStruct(e.Inputs)
	if err != nil {
		return nil, fmt.Errorf("invalid inputs data: %v", err)
	}
	return &terrariumpb.AppDependency{
		Id:          e.ID,
		Use:         e.Use,
		EnvPrefix:   e.EnvPrefix,
		Inputs:      inputs,
		Outputs:     e.Outputs,
		NoProvision: e.NoProvision,
	}, nil
}

func (e *Dependency) ScanProto(m *terrariumpb.AppDependency) {
	if m != nil {
		e.ID = m.Id
		e.Use = m.Use
		e.EnvPrefix = m.EnvPrefix
		e.Inputs = m.Inputs.AsMap()
		e.Outputs = m.Outputs
		e.NoProvision = m.NoProvision
	}
}

// IsEquivalent returns if this and the other dependency generate the same infrastructure.
func (e Dependency) IsEquivalent(other Dependency) bool {
	return e.ID == other.ID &&
		e.Use == other.Use &&
		e.NoProvision == other.NoProvision &&
		reflect.DeepEqual(e.Inputs, other.Inputs)
}

func NewApp(content []byte) (*App, error) {
	out := &App{}
	err := yaml.Unmarshal(content, out)
	if err != nil {
		return nil, err
	}

	return out, nil
}
