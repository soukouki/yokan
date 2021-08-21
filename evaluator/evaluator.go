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
	err, ok := checkTypeIsInteger("PlusPrefixOperator", right)
	if !ok { return err }
	return right
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	err, ok := checkTypeIsInteger("MinuPrefixOperator", right)
	if !ok { return err }
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
	case "==":
		return evalEqInfixOperatorExpression(left, right)
	case "!=":
		return evalNotEqInfixOperatorExpression(left, right)
	default:
		return &object.OtherError{Msg: fmt.Sprintf("Invalid operator '%s' in InfixExpression.", operator)}
	}
}

func evalPlusInfixOperatorExpression(left object.Object, right object.Object) object.Object {
	{
		err, ok := checkTypeIsInteger("PlusInfixOperator", left)
		if !ok { return err }
	}
	{
		err, ok := checkTypeIsInteger("PlusInfixOperator", right)
		if !ok { return err }
	}
	l := left.(*object.Integer).Value
	r := right.(*object.Integer).Value
	return &object.Integer{Value: l+r}
}

func evalMinusInfixOperatorExpression(left object.Object, right object.Object) object.Object {
	{
		err, ok := checkTypeIsInteger("MinusInfixOperator", left)
		if !ok { return err }
	}
	{
		err, ok := checkTypeIsInteger("MinusInfixOperator", right)
		if !ok { return err }
	}
	l := left.(*object.Integer).Value
	r := right.(*object.Integer).Value
	return &object.Integer{Value: l-r}
}

func evalStarInfixOperatorExpression(left object.Object, right object.Object) object.Object {
	{
		err, ok := checkTypeIsInteger("StarInfixOperator", left)
		if !ok { return err }
	}
	{
		err, ok := checkTypeIsInteger("StarInfixOperator", right)
		if !ok { return err }
	}
	l := left.(*object.Integer).Value
	r := right.(*object.Integer).Value
	return &object.Integer{Value: l*r}
}

func evalSlashInfixOperatorExpression(left object.Object, right object.Object) object.Object {
	{
		err, ok := checkTypeIsInteger("SlashInfixOperator", left)
		if !ok { return err }
	}
	{
		err, ok := checkTypeIsInteger("SlashInfixOperator", right)
		if !ok { return err }
	}
	l := left.(*object.Integer).Value
	r := right.(*object.Integer).Value
	if r==0 {
		return &object.OtherError{Msg: "Zero division Error"}
	}
	return &object.Integer{Value: l/r}
}

func checkTypeIsInteger(name string, val object.Object) (object.Object, bool) {
	if val.Type() != object.INTEGER_OBJ {
		return &object.TypeMisMatchError{Name: name, Expected: object.INTEGER_OBJ, Got: val}, false
	}
	return nil, true
}

var comparableInEqInfixOperatorTypes = []object.ObjectType {
	object.INTEGER_OBJ,
	object.STRING_OBJ,
	object.BOOLEAN_OBJ,
}
var comparableInEqInfixOperatorTypesName = object.INTEGER_OBJ+", "+object.STRING_OBJ+", "+object.BOOLEAN_OBJ

func evalEqInfixOperatorExpression(left object.Object, right object.Object) object.Object {
	if !contains(left.Type(), comparableInEqInfixOperatorTypes) {
		return &object.TypeMisMatchError{
			Name: "EqInfixOperator",
			Expected: comparableInEqInfixOperatorTypesName,
			Got: left,
		}
	}
	if !contains(right.Type(), comparableInEqInfixOperatorTypes) {
		return &object.TypeMisMatchError{
			Name: "EqInfixOperator",
			Expected: comparableInEqInfixOperatorTypesName,
			Got: right,
		}
	}
	if (
		left.Type() == object.INTEGER_OBJ &&
		right.Type() == object.INTEGER_OBJ &&
		left.(*object.Integer).Value == right.(*object.Integer).Value ||
		left.Type() == object.STRING_OBJ &&
		right.Type() == object.STRING_OBJ &&
		left.(*object.String).Value == right.(*object.String).Value ||
		left.Type() == object.BOOLEAN_OBJ &&
		right.Type() == object.BOOLEAN_OBJ &&
		left.(*object.Boolean).Value == right.(*object.Boolean).Value) {
		return &object.Boolean{Value: true}
	}
	return &object.Boolean{Value: false}
}

func evalNotEqInfixOperatorExpression(left object.Object, right object.Object) object.Object {
	obj := evalEqInfixOperatorExpression(left, right)
	if obj.Type() == object.BOOLEAN_OBJ {
		return &object.Boolean{Value: !obj.(*object.Boolean).Value}
	} else {
		return obj
	}
}

func contains(target object.ObjectType, types []object.ObjectType) bool {
	ret := false
	for _, t := range types {
		ret = ret || target == t
	}
	return ret
}
