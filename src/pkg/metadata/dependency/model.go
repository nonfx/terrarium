package dependency

import (
	"strings"

	"github.com/cldcvr/terrarium/src/pkg/jsonschema"
)

type Taxonomy string

type Dependency struct {
	Taxonomy    Taxonomy
	Title       string
	Description string
	Inputs      *jsonschema.Node
	Outputs     *jsonschema.Node
}

func (t Taxonomy) Parse() (taxons []string) {
	return strings.Split(string(t), "/")
}

func NewTaxonomy(taxons ...string) Taxonomy {
	return Taxonomy(strings.Join(taxons, "/"))
}
