package object

import (
	"fmt"
	"yokan/utility"
)

type ObjectType string

const (
	INTEGER_OBJ = "INTEGER"
	STRING_OBJ = "STRING"
	ERROR_OBJ = "ERROR"
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


// エラー

type Error interface {
	Object
	ErrorObject()
}

type TypeMisMatchError struct {
	Name string
	Wants string
	Got Object
}
func (e *TypeMisMatchError) ErrorObject() { }
func (e *TypeMisMatchError) Inspect() string {
	return fmt.Sprintf("%s wants %s but got '%s'", e.Name, e.Wants, e.Got.Type())
}
func (e *TypeMisMatchError) Type() ObjectType {
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
