package dependency

import (
	"strings"

	"github.com/cldcvr/terrarium/src/pkg/jsonschema"
)

type Taxonomy string

type Dependency struct {
	Taxonomy    Taxonomy         `yaml:"taxonomy"`
	Title       string           `yaml:"title"`
	Description string           `yaml:"description"`
	Inputs      *jsonschema.Node `yaml:"inputs"`
	Outputs     *jsonschema.Node `yaml:"outputs"`
}

func (t Taxonomy) Parse() (taxons []string) {
	return strings.Split(string(t), "/")
}

func NewTaxonomy(taxons ...string) Taxonomy {
	return Taxonomy(strings.Join(taxons, "/"))
}
