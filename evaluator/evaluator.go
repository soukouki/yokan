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
	case *ast.InfixExpression:
		return evalInfixExpression(Eval(node.Left), node.Operator, Eval(node.Right))

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
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
		return &object.TypeMisMatchError{Name: "PlusPrefixOperator", Wants: object.INTEGER_OBJ, Got: right}
	}
	return right
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return &object.TypeMisMatchError{Name: "MinusPrefixOperator", Wants: object.INTEGER_OBJ, Got: right}
	}
	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalInfixExpression(left object.Object, operator string, right object.Object) object.Object {
	switch operator {
	case "+":
		return evalPlusInfixOperatorExpression(left, right)
	case "-":
		return evalMinusInfixOperatorExpression(left, right)
	case "*":
		return evalStarInfixOperatorExpression(left, right)
	case "/":
		return evalSlashInfixOperatorExpression(left, right)
	default:
		return &object.OtherError{Msg: fmt.Sprintf("Invalid operator '%s' in InfixExpression.", operator)}
	}
}

func evalPlusInfixOperatorExpression(left object.Object, right object.Object) object.Object {
	if left.Type() != object.INTEGER_OBJ {
		return &object.TypeMisMatchError{Name: "PlusInfixOperator", Wants: object.INTEGER_OBJ, Got: left}
	}
	if right.Type() != object.INTEGER_OBJ {
		return &object.TypeMisMatchError{Name: "PlusInfixOperator", Wants: object.INTEGER_OBJ, Got: right}
	}
	l := left.(*object.Integer).Value
	r := right.(*object.Integer).Value
	return &object.Integer{Value: l+r}
}

func evalMinusInfixOperatorExpression(left object.Object, right object.Object) object.Object {
	if left.Type() != object.INTEGER_OBJ {
		return &object.TypeMisMatchError{Name: "MinusInfixOperator", Wants: object.INTEGER_OBJ, Got: left}
	}
	if right.Type() != object.INTEGER_OBJ {
		return &object.TypeMisMatchError{Name: "MinusInfixOperator", Wants: object.INTEGER_OBJ, Got: right}
	}
	l := left.(*object.Integer).Value
	r := right.(*object.Integer).Value
	return &object.Integer{Value: l-r}
}

func evalStarInfixOperatorExpression(left object.Object, right object.Object) object.Object {
	if left.Type() != object.INTEGER_OBJ {
		return &object.TypeMisMatchError{Name: "StarInfixOperator", Wants: object.INTEGER_OBJ, Got: left}
	}
	if right.Type() != object.INTEGER_OBJ {
		return &object.TypeMisMatchError{Name: "StarInfixOperator", Wants: object.INTEGER_OBJ, Got: right}
	}
	l := left.(*object.Integer).Value
	r := right.(*object.Integer).Value
	return &object.Integer{Value: l*r}
}

func evalSlashInfixOperatorExpression(left object.Object, right object.Object) object.Object {
	if left.Type() != object.INTEGER_OBJ {
		return &object.TypeMisMatchError{Name: "SlashInfixOperator", Wants: object.INTEGER_OBJ, Got: left}
	}
	if right.Type() != object.INTEGER_OBJ {
		return &object.TypeMisMatchError{Name: "SlashInfixOperator", Wants: object.INTEGER_OBJ, Got: right}
	}
	l := left.(*object.Integer).Value
	r := right.(*object.Integer).Value
	if r==0 {
		return &object.OtherError{Msg: "Zero division Error"}
	}
	return &object.Integer{Value: l/r}
}