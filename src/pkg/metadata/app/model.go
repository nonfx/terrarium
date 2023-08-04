package app

import "github.com/cldcvr/terrarium/src/pkg/metadata/taxonomy"

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

	// Service denotes a specific dependency that best classifies the app itself,
	// it can only be of the type `service.*`.
	// id of this dependency is automatically set to app id.
	// it is used to setup deployment pipeline in Code Pipes and allow other
	// apps to use this app as dependency.
	Service Dependency `yaml:"service"`

	// Dependencies lists the required services, databases, and other components that the application relies on.
	Dependencies Dependencies `yaml:"dependencies"`
}

type Dependencies []Dependency

// Dependency represents a single dependency of the application,
// which could be a database, another service, cache, etc.
type Dependency struct {
	// ID is a required identifier for the dependency in the project, which must start with
	// an alphabet character, can only contain alphanumeric characters,
	// and must not be longer than 20 characters.
	ID string `yaml:"id"`

	// Use indicates the specific taxon in the taxonomy hierarchy.
	Use taxonomy.Taxon `yaml:"type"`

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
