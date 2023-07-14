package parser

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"
)

func IsBool(expr hcl.Expression) bool {
	switch t := expr.(type) {
	case *hclsyntax.BinaryOpExpr:
		return true
	case *hclsyntax.LiteralValueExpr:
		return t.Val.Type() == cty.Bool
	case *hclsyntax.FunctionCallExpr:
		return t.Name == "anytrue" || t.Name == "alltrue"
	}
	return false
}

func IsObject(expr hcl.Expression) (ok bool) {
	_, ok = expr.(*hclsyntax.ObjectConsExpr)
	return
}

func IsCollection(expr hcl.Expression) (ok bool) {
	_, ok = expr.(*hclsyntax.ForExpr)
	return
}
