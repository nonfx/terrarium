package platform

type Component struct {
	ID           string
	ParentTaxons []string
	Title        string
	Description  string
	Inputs       map[string]interface{}
	Outputs      map[string]interface{}
	Extends      ComponentExtends
}

type ComponentExtends struct {
	ID     string
	Inputs map[string]interface{}
}

type GraphBlock struct {
	ID          string
	Connections []string
}

type Components []Component

type Graph []GraphBlock
