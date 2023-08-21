package dependency

import (
	"github.com/cldcvr/terrarium/src/pkg/jsonschema"
	"github.com/rotisserie/eris"
	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/yaml.v3"
)

// File represents the structure of the Dependency Interface file.
type File struct {
	// DependencyInterfaces holds the list of all interfaces defined in the file.
	DependencyInterfaces Interfaces `yaml:"dependency-interfaces,omitempty"`
}

// Interfaces is a slice of Interface, representing multiple Dependency Interfaces.
type Interfaces []Interface

// Interface represents a single Dependency Interface.
type Interface struct {
	ID string `yaml:"id,omitempty"`
	// Taxonomy is the identifier for the dependency represented by a Taxon.
	Taxonomy string `yaml:",omitempty"`
	// Title is the display title of the dependency.
	Title string `yaml:",omitempty"`
	// Description provides detailed information about the dependency.
	Description string `yaml:",omitempty"`
	// Inputs is a JSON Schema spec defining the structure of input attributes.
	Inputs *jsonschema.Node `yaml:",omitempty"`
	// Outputs is a JSON Schema spec defining the structure of output attributes.
	Outputs *jsonschema.Node `yaml:",omitempty"`
}

func (i *Interface) Init() {
	if i.Inputs != nil {
		i.Inputs.Type = gojsonschema.TYPE_OBJECT
	}
	if i.Outputs != nil {
		i.Outputs.Type = gojsonschema.TYPE_OBJECT
	}
}

func (iArr Interfaces) Init() {
	for i := range iArr {
		iArr[i].Init()
	}
}

func (f File) Init() {
	if f.DependencyInterfaces != nil {
		f.DependencyInterfaces.Init()
	}
}

// NewFile is a function that takes a byte slice of YAML data,
// unmarshals it into a File struct, and returns the struct.
// If there is an error during unmarshalling, it returns the error.
func NewFile(data []byte) (*File, error) {
	// Create a new File struct.
	f := &File{}
	// Unmarshal the YAML data into the File struct.
	err := yaml.Unmarshal(data, f)
	if err != nil {
		return nil, eris.Wrapf(err, "failed to unmarshal data")
	}

	f.Init()

	// If there is no error, return the File struct.
	return f, nil
}
