package main

import (
	"fmt"
	"io/ioutil"
	"lime/backend"
	"reflect"
	"regexp"
	"strings"
)

var re = regexp.MustCompile(`\p{Lu}`)

func pyname(in string) string {
	if in == "String" {
		return "Str"
	}
	return re.ReplaceAllStringFunc(in, func(a string) string { return "_" + strings.ToLower(a) })
}

func pytype(t reflect.Type) string {
	switch t.Kind() {
	case reflect.Struct:
		return "*" + t.Name()
	case reflect.Int:
		return "*py.Int"
	default:
		panic(t.Kind())
	}
}

func pyret(ot reflect.Type) (ret string) {
	switch ot.Kind() {
	case reflect.Struct:
		return fmt.Sprintf(`
	pyret, err := %sClass.Alloc(1)
	if err != nil {
		return nil, err
	} else if v2, ok := pyret.(*%s); !ok {
		return nil, fmt.Errorf("Unable to convert return value to the right type?!: %%s", pyret)
	} else {
		v2.data = ret
		return v2, nil
	}`, pyname(ot.Name()), ot.Name())
	case reflect.Bool:
		return `
	if ret {
		return py.True, nil
	} else {
		return py.False, nil
	}`
	case reflect.Int:
		return "\n\treturn py.NewInt(ret), nil"
	default:
		panic(ot.Kind())
	}
	return
}

func pyacc(ot reflect.Type) string {
	switch ot.Kind() {
	case reflect.Ptr, reflect.Struct:
		return ".data"
	case reflect.Int:
		return ".Int()"
	default:
		panic(ot.Kind())
	}
}

func generateWrapper(t reflect.Type, canCreate bool) (ret string) {
	if t.Kind() != reflect.Struct {
		panic(t.Kind())
	}
	it := t.String()
	if !canCreate {
		it = "*" + it
	}
	ret += fmt.Sprintf(`
var %sClass = py.Class{
	Name:    "sublime.%s",
	Pointer: (*%s)(nil),
}

type %s struct {
	py.BaseObject
	data %s
}
`, pyname(t.Name()), t.Name(), t.Name(), t.Name(), it)

	for i := 0; i < t.NumMethod(); i++ {
		var (
			m    = t.Method(i)
			args string
			rv   string
			in   = m.Type.NumIn() - 1
			out  = m.Type.NumOut()
		)
		if out > 1 {
			panic(fmt.Errorf("Can't handle out > 1: %d", out))
		}

		if in > 0 {
			args = "tu *py.Tuple, kw *py.Dict"
		}
		if m.Name == "String" {
			rv = "string"
		} else if m.Type.NumOut() == 0 {
			rv = "error"
		} else {
			rv = "(py.Object, error)"
		}

		ret += fmt.Sprintf(`
func (o *%s) Py%s(%s) %s {`, t.Name(), pyname(m.Name), args, rv)

		if in > 0 {
			ret += "\n\tvar ("
			for j := 1; j <= in; j++ {
				ty := pytype(m.Type.In(j))
				ret += fmt.Sprintf("\n\t\targ%d %s", j, ty)
			}
			ret += "\n\t)"
			ret += fmt.Sprintf(`
	if tu.Size() != %d {
		return nil, fmt.Errorf("Expected %d arguments but got %%d", tu.Size())
	}`, in, in)

			for j := 1; j <= in; j++ {
				ty := pytype(m.Type.In(j))
				ret += fmt.Sprintf(`
	if v, err := tu.GetItem(%d); err != nil {
		return nil, err
	} else if v2, ok := v.(%s); !ok {
		return nil, fmt.Errorf("Expected type %s for argument %d, not %%s", v.Type())
	} else {
		arg%d = v2
	}`, j-1, ty, ty, j, j)
			}
		}

		call := ""
		if in > 0 {
			call = "o.data." + m.Name + "("
			for j := 1; j <= in; j++ {
				if j > 1 {
					call += ", "
				}
				call += fmt.Sprintf("arg%d", j)
				call += pyacc(m.Type.In(j))
				call += ")"
			}
		} else {
			call = "o.data." + m.Name + "()"
		}
		if m.Name == "String" {
			ret += "\n\treturn " + call
		} else if out > 0 {
			ret += "\n\tret := " + call
			ret += pyret(m.Type.Out(0))
		} else {
			ret += "\n\t" + call + "\n\treturn nil"
		}
		ret += "\n}\n"
	}

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if !f.Anonymous {
			ret += fmt.Sprintf(`
func (o *%s) PyGet%s() (py.Object, error) {
	ret := o.data.%s%s
}

func (o *%s) PySet%s(v py.Object) error {
	if v2, ok := v.(%s); !ok {
		return fmt.Errorf("Expected type %s not %%s", v.Type())
	} else {
		o.data.%s = v2%s
		return nil
	}
}

`, t.Name(), pyname(f.Name), f.Name, pyret(f.Type),
				t.Name(), pyname(f.Name), pytype(f.Type), f.Name, f.Name, pyacc(f.Type),
			)
		}
	}

	return
}

func main() {
	data := [][]string{
		{"../backend/sublime/region.go", generateWrapper(reflect.TypeOf(backend.Region{}), true)},
	}
	for _, gen := range data {
		wr := `// This file was generated as part of a build step and shouldn't be manually modified
package sublime

import (
	"fmt"
	"github.com/qur/gopy/lib"
	"lime/backend"
)
` + gen[1]
		if err := ioutil.WriteFile(gen[0], []byte(wr), 0644); err != nil {
			panic(err)
		}
	}
}
