package evaluator

import (
	"fmt"

	"yokan/ast"
	"yokan/object"
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	
	case *ast.PrefixExpression:
		return evalPrefixExpression(node.Operator, Eval(node.Right))

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	}
	return &object.OtherError{Msg: "Not yet implemented"}
}

func evalStatements(stmts []ast.Statement) object.Object {
	var result object.Object
	for _, stmt := range stmts {
		result = Eval(stmt)
	}
	return result
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "+":
		return evalPlusPrefixOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return &object.OtherError{Msg: fmt.Sprintf("Invalid operator '%s' in PrefixExpression.", operator)}
	}
}

func evalPlusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return &object.TypeMisMatchError{Name: "PlusPrefixOperator", Wants: "INTEGER", Got: right}
	}
	return right
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return &object.TypeMisMatchError{Name: "MinusPrefixOperator", Wants: "INTEGER", Got: right}
	}
	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}