package main

import "fmt"

const OPENAI_MODEL = "text-embedding-ada-002" // "text-embedding-ada-002" //
var OPENAI_PROMPT = 0

type EmbeddingsIndex map[string]map[string][]float64 // ResourceType -> ResourceAttribute -> Embeddings

func (e EmbeddingsIndex) GetByTFToken(t TFToken) []float64 {
	if resType, exists := e[t.ResourceType]; exists {
		if resAttrVec, exists := resType[t.ResourceAttribute]; exists {
			return resAttrVec
		}
	}

	return nil
}

func (e *EmbeddingsIndex) SetByTFToken(t TFToken, vec []float64) EmbeddingsIndex {
	if e == nil {
		e = &EmbeddingsIndex{}
	}

	if _, exists := (*e)[t.ResourceType]; !exists {
		(*e)[t.ResourceType] = map[string][]float64{}
	}

	(*e)[t.ResourceType][t.ResourceAttribute] = vec

	return *e
}

type TFToken struct {
	ResourceType      string
	ResourceAttribute string
}

func (t TFToken) Prompt() string {
	switch OPENAI_PROMPT {
	case 0:
		return t.String()
	case 1:
		return fmt.Sprintf("terraform resource attribute: %s.%s", t.ResourceType, t.ResourceAttribute)
	case 2:
		return fmt.Sprintf("terraform resource type: %s; terraform resource attribute: %s;", t.ResourceType, t.ResourceAttribute)
	case 3:
		return fmt.Sprintf("%s of %s", t.ResourceAttribute, t.ResourceType)
	case 4:
		return fmt.Sprintf("%s in %s", t.ResourceAttribute, t.ResourceType)
	case 5:
		return fmt.Sprintf("terraform resource attribute `%s` in resource type `%s`", t.ResourceAttribute, t.ResourceType)
	case 6:
		return fmt.Sprintf("`%s` is a terraform resource attribute in resource type `%s.%s`", t.ResourceAttribute, t.ResourceType, t.ResourceAttribute)
	case 7:
		return fmt.Sprintf("`%s.%s` provides terraform cloud resource attribute `%s`", t.ResourceType, t.ResourceAttribute, t.ResourceAttribute)
	}
	panic(fmt.Errorf("invalid prompt number"))
}

func (t TFToken) String() string {
	return fmt.Sprintf("%s.%s", t.ResourceType, t.ResourceAttribute)
}

type OpenAIEmbeddingRespObj struct {
	Object    string
	Index     int
	Embedding []float64
}

type OpenAIEmbeddingResp struct {
	Object string
	Data   []OpenAIEmbeddingRespObj
}
