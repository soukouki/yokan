package object

func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s, parent: nil}
}

func NewInferitEnvironment(parent *Environment) *Environment {
	env := NewEnvironment()
	env.parent = parent
	return env
}


type Environment struct {
	store map[string]Object
	parent *Environment
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.parent != nil {
		obj, ok = e.parent.Get(name)
	}
	return obj, ok
}

func (e *Environment) Set(name string, val Object) {
	e.store[name] = val
}
