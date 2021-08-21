package object

import (
	"fmt"
)

var Buildins = map[string]Object{
	"true": &Boolean{Value: true},
	"false": &Boolean{Value: false},
	"null": &Null{ },
	"puts": &Buildin{
		Fn: func(args ...Object) Object {
			for _, arg := range args {
				fmt.Println(arg.String())
			}
			return &Null{ }
		},
	},
	"if": &Buildin{
		Fn: func(args ...Object) Object {
			if len(args) != 3 {
				return &OtherError{Msg: fmt.Sprintf("if need 3 arguments. but got %d", len(args))}
			}
			cond := args[0]
			t := args[1]
			f := args[2]
			if cond.Type() != BOOLEAN_OBJ {
				return &TypeMisMatchError{Name: "if", Expected: BOOLEAN_OBJ, Got: cond}
			}
			if cond.(*Boolean).Value {
				return t
			} else {
				return f
			}
		},
	},
}