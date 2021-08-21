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
}