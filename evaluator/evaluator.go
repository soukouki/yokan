package evaluator

import (
	"fmt"

	"yokan/ast"
	"yokan/object"
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalStatements(node.Statements, env)
	
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.Assign:
		return evalAssign(*node, env)
	
	case *ast.FunctionLiteral:
		return &object.Function{Parameters: node.Arguments, Body: node.Body, Env: env}
	case *ast.FunctionCalling:
		function := Eval(node.Function, env)
		if isError(function) { return function }
		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) { return args[0] }
		return applyFunction(function, args)
	
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) { return right }
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) { return left }
		right := Eval(node.Right, env)
		if isError(right) { return right }
		return evalInfixExpression(left, node.Operator, right)
	case *ast.Identifier:
		name := node.Name
		val, ok := env.Get(name)
		if !ok { return &object.OtherError{Msg: fmt.Sprintf("%s is unbouded variable", name)} }
		return val

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	}
	return &object.OtherError{Msg: fmt.Sprintf("%T is not yet implemented", node)}
}

func evalStatements(stmts []ast.Statement, env *object.Environment) object.Object {
	var result object.Object = &object.Null{ } 
	for _, stmt := range stmts {
		result = Eval(stmt, env)
		if isError(result) {
			return result
		}
	}
	return result
}

func evalAssign(assign ast.Assign, env *object.Environment) object.Object {
	val := Eval(assign.Value, env)
	if isError(val) { return val }
	env.Set(assign.Name.Name, val)
	return &object.ReturnValueOsStatement{ }
}

func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object
	for _, e := range exps {
		evaled := Eval(e, env)
		if isError(evaled) { return []object.Object{evaled} }
		result = append(result, evaled)
	}
	return result
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	function, ok := fn.(*object.Function)
	if !ok {
		return &object.OtherError {
			Msg: fmt.Sprintf("%s(%s) is not a function", fn.Type(), fn.Inspect()),
		}
	}
	if len(function.Parameters) != len(args) {
		return &object.OtherError {
			Msg: fmt.Sprintf("Function need %d params, but got %d params", len(function.Parameters), len(args)),
		}
	}
	inheritEnv := inheritFunctionEnv(function, args)
	return evalStatements(function.Body, inheritEnv)
}

func inheritFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewInferitEnvironment(fn.Env)
	for paramIdx, param := range fn.Parameters {
		env.Set(param.Name, args[paramIdx])
	}
	return env
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
	case "<":
		return evalLTInfixOperatorExpression(left, right)
	case "<=":
		return evalLTEQInfixOperatorExpression(left, right)
	case ">":
		return evalGTInfixOperatorExpression(left, right)
	case ">=":
		return evalGTEQInfixOperatorExpression(left, right)
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
	return not(evalEqInfixOperatorExpression(left, right))
}

func evalLTInfixOperatorExpression(left object.Object, right object.Object) object.Object {
	{
		err, ok := checkTypeIsInteger("LTInfixOperator", left)
		if !ok { return err }
	}
	{
		err, ok := checkTypeIsInteger("LTInfixOperator", right)
		if !ok { return err }
	}
	l := left.(*object.Integer).Value
	r := right.(*object.Integer).Value
	return &object.Boolean{Value: l<r} 
}

func evalLTEQInfixOperatorExpression(left object.Object, right object.Object) object.Object {
	{
		err, ok := checkTypeIsInteger("LTInfixOperator", left)
		if !ok { return err }
	}
	{
		err, ok := checkTypeIsInteger("LTInfixOperator", right)
		if !ok { return err }
	}
	l := left.(*object.Integer).Value
	r := right.(*object.Integer).Value
	return &object.Boolean{Value: l<=r} 
}

func evalGTInfixOperatorExpression(left object.Object, right object.Object) object.Object {
	return not(evalLTEQInfixOperatorExpression(left, right))
}

func evalGTEQInfixOperatorExpression(left object.Object, right object.Object) object.Object {
	return not(evalLTInfixOperatorExpression(left, right))
}

func not(obj object.Object) object.Object {
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

func isError(obj object.Object) bool {
	return obj.Type() == object.ERROR_OBJ
}