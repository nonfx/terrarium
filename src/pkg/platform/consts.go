package platform

import (
	"fmt"
	"regexp"
)

const (
	COMPONENT_PREFIX    = "tr_component_"
	TAXON_PREFIX        = "tr_taxon_"
	ENABLED_BOOL_SUFFIX = "_enabled"
)

var (
	componentNamePattern = regexp.MustCompile(fmt.Sprintf("^%s(.*?)(%s)?$", COMPONENT_PREFIX, ENABLED_BOOL_SUFFIX))
)

type TFBlockType int

const (
	TFBlock_Undefined = iota
	TFBlock_Module
	TFBlock_Resource
	TFBlock_Data
	TFBlock_Local
	TFBlock_Variable
	TFBlock_Output
)

type TFBlock struct {
	ID          string
	Type        TFBlockType
	Connections Graph
}

type Graph []TFBlock

func ParseComponentName(blockName string) (componentName string) {
	if !componentNamePattern.MatchString(blockName) {
		return ""
	}

	return componentNamePattern.ReplaceAllString(blockName, "$1")
}
