package object

import (
	"fmt"
	"yokan/ast"
	"yokan/utility"
)

type ObjectType string

const (
	FUNCTION_OBJ = "FUNCTION"

	INTEGER_OBJ = "INTEGER"
	STRING_OBJ = "STRING"
	BOOLEAN_OBJ = "BOOLEAN"
	NULL_OBJ = "NULL"
	
	ERROR_OBJ = "ERROR"
	SHOULD_NOT_VIEWABLE_OBJ = "SHOULD_NOT_VIEWABLE"
)

type Object interface {
	Type() ObjectType
	String() string
}


type Function struct {
	Parameters []ast.Identifier
	Body []ast.Statement
	// 実行時じゃなくて定義時の環境を持たないといけないので、Functionが環境を保つ必要がある
	Env *Environment
}
func (f *Function) String() string {
	var param []string
	for _, p := range f.Parameters {
		param = append(param, p.Name)
	}
	var body []string
	for _, b := range f.Body {
		body = append(body, b.String())
	}
	return utility.FunctionString(param, body)
}
func (f *Function) Type() ObjectType {
	return FUNCTION_OBJ
}


// 値

type Integer struct {
	Value int64
}
func (i *Integer) String() string {
	return fmt.Sprintf("%d", i.Value)
}
func (i *Integer) Type() ObjectType {
	return INTEGER_OBJ
}

type String struct {
	Value string
}
func (s *String) String() string {
	return utility.Quote(s.Value)
}
func (s *String) Type() ObjectType {
	return STRING_OBJ
}

type Boolean struct {
	Value bool
}
func (b *Boolean) String() string {
	if b.Value {
		return "true"
	} else {
		return "false"
	}
}
func (b *Boolean) Type() ObjectType {
	return BOOLEAN_OBJ
}

type Null struct { }
func (n *Null) String() string {
	return "null"
}
func (n *Null) Type() ObjectType {
	return NULL_OBJ
}

// エラー

type Error interface {
	Object
	ErrorObject()
}

type TypeMisMatchError struct {
	Name string
	Expected string
	Got Object
}
func (e *TypeMisMatchError) ErrorObject() { }
func (e *TypeMisMatchError) String() string {
	return fmt.Sprintf("%s Expected %s but got '%s'", e.Name, e.Expected, e.Got.Type())
}
func (e *TypeMisMatchError) Type() ObjectType {
	return ERROR_OBJ
}

type OtherError struct {
	Msg string
}
func (e *OtherError) ErrorObject() { }
func (e *OtherError) String() string {
	return e.Msg
}
func (e *OtherError) Type() ObjectType {
	return ERROR_OBJ
}


// 文の戻り値

type ReturnValueOsStatement struct { }
func (r *ReturnValueOsStatement) String() string {
	return "THIS VALUE SHOULD NOT VIEWABLE"
}
func (e *ReturnValueOsStatement) Type() ObjectType {
	return SHOULD_NOT_VIEWABLE_OBJ
}