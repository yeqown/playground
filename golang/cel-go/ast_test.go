package main

import (
	"testing"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/operators"
	"github.com/google/cel-go/test"

	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

func Test_ast(t *testing.T) {
	exp := test.ExprIdent(1, "doctor")
	exp2 := test.ExprIdent(8, "doctor.id")
	exp3 := test.ExprIdent(8, "doctor.age")

	var (
		// doctor.id + 12 > 3000 && doctor.age > 30
		parsed = &exprpb.ParsedExpr{
			Expr: test.ExprCall(6, operators.LogicalAnd,
				test.ExprCall(7, operators.Add, test.ExprSelect(8, exp, "id"), test.ExprLiteral(3, int64(12))),
				test.ExprCall(9, operators.Add, test.ExprSelect(10, exp, "age"), test.ExprLiteral(6, int64(30))),
			),
		}
		// 也可以
		// doctor.id + 12 > 3000 && doctor.age > 30
		parsed2 = &exprpb.ParsedExpr{
			Expr: test.ExprCall(6, operators.LogicalAnd,
				test.ExprCall(7, operators.Add, exp2, test.ExprLiteral(3, int64(12))),
				test.ExprCall(9, operators.Add, exp3, test.ExprLiteral(6, int64(30))),
			),
		}
	)

	_ = parsed
	ast := cel.ParsedExprToAst(parsed2)
	s, err := cel.AstToString(ast)
	if err != nil {
		panic(err)
	}

	println(s)
}

// 1. 字面值
// 2. 子函数操作
type exprNode struct {
	// 1 + 2 - 3 * 4 / 5 .
	op uint

	expr exprValue

	left  exprValue
	right exprValue

	operands []*exprValue
}

type exprValue interface {
}

type exprTree struct {
	root *exprNode

	id int64
}

var _ = `
{
	op: 0,
	expr: "doctor.id"
	left: null,
	right: null,
	operands: [],
}
`

// doctor.id + 32
var _ = `
{
	op: 1,
	expr: "",
	left: {
		op: 0,
		expr: "doctor.id",
		left: null,
		right: null,
		operands: []
	},
	right: {
		op: 0,
		expr: "32",
		left: null,
		right: null,
		operands: []
	},
	operands: [],
}
`
