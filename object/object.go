package object

import (
	"fmt"
	"yokan/utility"
)

type ObjectType string

const (
	INTEGER_OBJ = "INTEGER"
	STRING_OBJ = "STRING"
	BOOLEAN_OBJ = "BOOLEAN"
	
	ERROR_OBJ = "ERROR"
	SHOULD_NOT_VIEWABLE_OBJ = "SHOULD_NOT_VIEWABLE"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}


// 値

type Integer struct {
	Value int64
}
func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}
func (i *Integer) Type() ObjectType {
	return INTEGER_OBJ
}

type String struct {
	Value string
}
func (s *String) Inspect() string {
	return utility.Quote(s.Value)
}
func (s *String) Type() ObjectType {
	return STRING_OBJ
}

type Boolean struct {
	Value bool
}
func (b *Boolean) Inspect() string {
	if b.Value {
		return "true"
	} else {
		return "false"
	}
}
func (b *Boolean) Type() ObjectType {
	return BOOLEAN_OBJ
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
func (e *TypeMisMatchError) Inspect() string {
	return fmt.Sprintf("%s Expected %s but got '%s'", e.Name, e.Expected, e.Got.Type())
}
func (e *TypeMisMatchError) Type() ObjectType {
	return ERROR_OBJ
}

type UnboundedVariableError struct {
	Name string
}
func (e *UnboundedVariableError) ErrorObject() { }
func (e *UnboundedVariableError) Inspect() string {
	return fmt.Sprintf("%s is unbouded variable", e.Name)
}
func (e *UnboundedVariableError) Type() ObjectType {
	return ERROR_OBJ
}

type OtherError struct {
	Msg string
}
func (e *OtherError) ErrorObject() { }
func (e *OtherError) Inspect() string {
	return e.Msg
}
func (e *OtherError) Type() ObjectType {
	return ERROR_OBJ
}


// 文の戻り値

type ReturnValueOsStatement struct { }
func (r *ReturnValueOsStatement) Inspect() string {
	return "THIS VALUE SHOULD NOT VIEWABLE"
}
func (e *ReturnValueOsStatement) Type() ObjectType {
	return SHOULD_NOT_VIEWABLE_OBJ
}