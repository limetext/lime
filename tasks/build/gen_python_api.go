// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

// This file generates python bindings for the backend api.
// Similar code would need to be written to generate bindings for other
// scripting languages, though much code is likely shared.

package main

import (
	"fmt"
	"github.com/limetext/lime/backend"
	"github.com/limetext/text"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var re = regexp.MustCompile(`\p{Lu}`)

func pyname(in string) string {
	switch in {
	case "RowCol":
		return "_rowcol"
	case "String":
		return "Str"
	case "Len":
		return "SeqLen"
	case "Get":
		return "SeqGet"
	}
	return re.ReplaceAllStringFunc(in, func(a string) string { return "_" + strings.ToLower(a) })
}

func pyretvar(name string, ot reflect.Type) (string, error) {
	return fmt.Sprintf("\npy%s, err = toPython(%s)", name, name), nil
}

func pyret(ot reflect.Type) (string, error) {
	v, err := pyretvar("ret", ot)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(`
				var pyret py.Object
				var err error
				%s
				return pyret, err
				`, v), nil
}

func typename(t reflect.Type) string {
	switch t.Kind() {
	case reflect.Slice:
		return "[]" + typename(t.Elem())
	case reflect.Ptr:
		return "*" + typename(t.Elem())
	default:
		return t.String()
	}
}
func pytogoconv(in, set, name string, returnsValue bool, t reflect.Type) (string, error) {
	r := ""
	if returnsValue {
		r = "nil, "
	}
	ty := typename(t)
	return fmt.Sprintf(`
		if v3, err2 := fromPython(%s);  err2 != nil {
			return %serr2
		} else {
	 		if v2, ok := v3.(%s); !ok {
	 			return %sfmt.Errorf("Expected type %s for %s, not %%s", %s.Type())
	 		} else {
	 			%s = v2
	 		}
	 	}`, in, r, ty, r, ty, name, in, set), nil
}

func generatemethod(m reflect.Method, t2 reflect.Type, callobject, name string) (ret string, err error) {
	var (
		args string
		rv   string
		in   = m.Type.NumIn() - 1
		out  = m.Type.NumOut()
		call string
	)

	if in > 0 {
		args = "tu *py.Tuple"
	}
	if m.Name == "Get" && in == 1 && m.Type.In(1).Kind() == reflect.Int && out == 1 {
		args = "arg0 int64"
	}
	if m.Name == "String" {
		rv = "string"
	} else if m.Name == "Len" {
		rv = "int64"
	} else {
		rv = "(py.Object, error)"
	}

	ret += fmt.Sprintf("\nfunc %s (%s) %s {", name, args, rv)

	if m.Name == "Get" && in == 1 && m.Type.In(1).Kind() == reflect.Int && out == 1 {
		ret += `var (
				pyret0 py.Object
				err error
			)
			if l := ` + callobject + `Len(); int(arg0) >= l || arg0 < 0 {
				return nil, py.NewError(py.IndexError, "%d >= %d || %d < 0", arg0, l, arg0)
			}
			`
		ret += "\nret0 := " + callobject + m.Name + "(int(arg0))"
		if r, err := pyretvar("ret0", m.Type.Out(0)); err != nil {
			return "", err
		} else {
			ret += r
			ret += `
						if err != nil {
							// TODO: do the py objs need to be freed?
							return nil, err
						}
						`
		}
		ret += "return pyret0, err\n}\n"
		return ret, nil
	} else if m.Name == "Len" {
		return ret + "return int64(" + callobject + m.Name + "())\n}\n", nil
	}
	if in > 0 {
		ret += "\n\tvar ("
		for j := 1; j <= in; j++ {
			ret += fmt.Sprintf("\n\t\targ%d %s", j, m.Type.In(j))
		}
		r := ""
		if m.Name != "String" {
			r = "nil, "
		}
		ret += "\n\t)"

		for j := 1; j <= in; j++ {
			t := m.Type.In(j)
			name := fmt.Sprintf("arg%d", j)
			msg := fmt.Sprintf("%s.%s() %s", t2, m.Name, name)
			pygo, err := pytogoconv("v", name, msg, m.Name != "String", t)
			if err != nil {
				return "", err
			}
			if t.Kind() == reflect.Map && t.Key().Kind() == reflect.String {
				ret += fmt.Sprintf(`
						%s = make(%s)
						if v, err := tu.GetItem(%d); err == nil {%s}`, name, t, j-1, pygo)
			} else {
				ret += fmt.Sprintf(`
						if v, err := tu.GetItem(%d); err != nil {
							return %serr
						} else {%s}`, j-1, r, pygo)
			}
		}
	}

	if in > 0 {
		call = callobject + m.Name + "("
		for j := 1; j <= in; j++ {
			if j > 1 {
				call += ", "
			}
			call += fmt.Sprintf("arg%d", j)
		}
		call += ")"
	} else {
		call = callobject + m.Name + "()"
	}
	if m.Name == "String" {
		ret += "\n\treturn " + call
	} else if out == 1 && m.Type.Out(0).Name() == "error" {
		ret += fmt.Sprintf(`
			if err := %s; err != nil {
				return nil, err
			} else {
				return toPython(nil)
			}`, call)
	} else if out > 0 {
		ret += "\n\t"
		for j := 0; j < out; j++ {
			if j > 0 {
				ret += ", "
			}
			ret += fmt.Sprintf("ret%d", j)
		}
		ret += " := " + call
		ret += "\nvar err error"
		for j := 0; j < out; j++ {
			ret += fmt.Sprintf("\nvar pyret%d py.Object\n", j)
			if r, err := pyretvar(fmt.Sprintf("ret%d", j), m.Type.Out(j)); err != nil {
				return "", err
			} else {
				ret += r
				ret += "\nif err != nil {"
				for k := 0; k < j; k++ {
					ret += fmt.Sprintf("\npyret%d.Decref()", k)
				}
				ret += "\nreturn nil, err\n}"
			}
		}
		if out == 1 {
			ret += "\n\treturn pyret0, err"
		} else {
			// TODO: does PackTuple take ownership of the objects?
			ret += "\n\treturn py.PackTuple("
			for j := 0; j < out; j++ {
				if j > 0 {
					ret += ", "
				}
				ret += fmt.Sprintf("pyret%d", j)
			}
			ret += ")"
		}
	} else {
		ret += "\n\t" + call + "\n\treturn toPython(nil)"
	}
	ret += "\n}\n"
	return
}

func generatemethodsEx(t reflect.Type, ignorefunc func(name string) bool, callobject string, name func(t reflect.Type, m reflect.Method) string) (methods string) {
	t2 := t
	if t.Kind() == reflect.Ptr {
		t2 = t.Elem()
	}

	for i := 0; i < t.NumMethod(); i++ {
		var (
			m      = t.Method(i)
			reason string
		)

		if m.Name[0] != strings.ToUpper(m.Name[:1])[0] {
			reason = "unexported"
			goto skip
		}
		if ignorefunc != nil && ignorefunc(m.Name) {
			reason = "in skip list"
			goto skip
		}

		if m, err := generatemethod(m, t2, callobject, name(t2, m)); err != nil {
			reason = err.Error()
			goto skip
		} else {
			methods += m
		}

		continue
	skip:
		fmt.Printf("Skipping method %s.%s: %s\n", t2, m.Name, reason)
	}
	return

}

func generatemethods(t reflect.Type, ignorefunc func(name string) bool) (methods string) {
	return generatemethodsEx(t, ignorefunc, "o.data.", func(t2 reflect.Type, m reflect.Method) string {
		return fmt.Sprintf("\n(o *%s) Py%s", t2.Name(), pyname(m.Name))
	})
}

func generateWrapper(ptr reflect.Type, canCreate bool, ignorefunc func(name string) bool) (ret string) {
	t := ptr
	if t.Kind() == reflect.Ptr {
		t = ptr.Elem()
	}
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

	cons := ""
	if canCreate {
		cons = fmt.Sprintf(`
			func (o *%s) PyInit(args *py.Tuple, kwds *py.Dict) error {
				if args.Size() > %d {
					return fmt.Errorf("Expected at most %d arguments")
				}
			`, t.Name(), t.NumField(), t.NumField())
		ok := true
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			pygo, err := pytogoconv("v", "o.data."+f.Name, t.Name()+"."+f.Name, false, f.Type)
			if err != nil {
				ok = false
				break
			}
			cons += fmt.Sprintf(`
					if args.Size() > %d {
						if v, err := args.GetItem(%d); err != nil {
							return err
						} else {%s
						}
					}
				`, i, i, pygo)
		}
		if !ok {
			cons = ""
		} else {
			cons += "\n\treturn nil\n}"
		}
	}
	if cons == "" {
		ret += fmt.Sprintf(`
			func (o *%s) PyInit(args *py.Tuple, kwds *py.Dict) error {
				return fmt.Errorf("Can't initialize type %s")
			}`, t.Name(), t.Name())
	}
	ret += cons
	ret += generatemethods(ptr, ignorefunc)
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if !f.Anonymous && f.Name[0] == strings.ToUpper(f.Name[:1])[0] {
			if ignorefunc != nil && ignorefunc(f.Name) {
				goto skip
			}

			if r, err := pyret(f.Type); err != nil {
				fmt.Printf("Skipping field %s.%s: %s\n", t.Name(), f.Name, err)
			} else if pygo, err := pytogoconv("v", "o.data."+f.Name, t.Name()+"."+f.Name, false, f.Type); err != nil {
				fmt.Printf("Skipping field %s.%s: %s\n", t.Name(), f.Name, err)
			} else {
				ret += fmt.Sprintf(`
					func (o *%s) PyGet%s() (py.Object, error) {
						ret := o.data.%s%s
					}

					func (o *%s) PySet%s(v py.Object) error {%s
						return  nil
					}

					`, t.Name(), pyname(f.Name), f.Name, r,
					t.Name(), pyname(f.Name), pygo,
				)
			}
		skip:
		}
	}

	return
}

const path = "./backend/sublime"

var keep = regexp.MustCompile(`^(.+(_test|_manual)\.go|.+\.py)|(doc\.go)$`)

func cleanup() {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if fi, err := f.Readdir(0); err != nil {
		panic(err)
	} else {
		for _, f := range fi {
			if !f.IsDir() && !keep.MatchString(f.Name()) {
				fn := filepath.Join(path, f.Name())
				fmt.Println("removing", fn)
				os.Remove(fn)
			}
		}
	}
}

func main() {
	cleanup()
	var sublime_methods = ""
	sn := func(t reflect.Type, m reflect.Method) string {
		sn := "sublime_" + m.Name
		sublime_methods += fmt.Sprintf("{Name: \"%s\", Func: %s},\n", pyname(m.Name)[1:], sn)
		return sn
	}
	data := [][]string{
		{"./backend/sublime/region_generated.go", generateWrapper(reflect.TypeOf(text.Region{}), true, regexp.MustCompile("Cut").MatchString)},
		{"./backend/sublime/regionset_generated.go", generateWrapper(reflect.TypeOf(&text.RegionSet{}), false, regexp.MustCompile("Less|Swap|Adjust|Has|Cut").MatchString)},
		{"./backend/sublime/edit_generated.go", generateWrapper(reflect.TypeOf(&backend.Edit{}), false, regexp.MustCompile("Apply|Undo").MatchString)},
		{"./backend/sublime/view_generated.go", generateWrapper(reflect.TypeOf(&backend.View{}), false, regexp.MustCompile("Buffer|Syntax|CommandHistory|Show|AddRegions|UndoStack|Transform|Reload|Save|Close|ExpandByClass|Erased|FileChanged|Inserted|Find$|^Status").MatchString)},
		{"./backend/sublime/window_generated.go", generateWrapper(reflect.TypeOf(&backend.Window{}), false, regexp.MustCompile("OpenFile|SetActiveView|Close").MatchString)},
		{"./backend/sublime/settings_generated.go", generateWrapper(reflect.TypeOf(&text.Settings{}), false, regexp.MustCompile("Parent|Set|Get|UnmarshalJSON|MarshalJSON").MatchString)},
		{"./backend/sublime/view_buffer_generated.go", generatemethodsEx(
			reflect.TypeOf(backend.GetEditor().Console().Buffer()),
			regexp.MustCompile("Erase|Insert|Substr|SetFile|AddCallback|AddObserver|RemoveObserver|Data|Runes|Settings|Index|Close|Unlock|Lock").MatchString,
			"o.data.Buffer().",
			func(t reflect.Type, m reflect.Method) string {
				mn := ""
				switch m.Name {
				case "Line", "LineR", "FullLine", "FullLineR", "WordR", "Word":
					mn = strings.ToLower(m.Name)
				case "Id":
					mn = "Py_buffer_id"
				default:
					mn = "Py" + pyname(m.Name)
				}
				return "(o *View) " + mn
			})},
		{"./backend/sublime/commands_generated.go", generatemethodsEx(reflect.TypeOf(backend.GetEditor().CommandHandler()),
			regexp.MustCompile("RunWindowCommand|RunTextCommand|RunApplicationCommand|RegisterWithDefault").MatchString,
			"backend.GetEditor().CommandHandler().",
			sn),
		},
		{"./backend/sublime/frontend_generated.go", generatemethodsEx(reflect.TypeOf(backend.GetEditor().Frontend()),
			regexp.MustCompile("Show|VisibleRegion|DefaultAction").MatchString,
			"backend.GetEditor().Frontend().",
			sn),
		},
		{"./backend/sublime/sublime_api_generated.go", generatemethodsEx(reflect.TypeOf(backend.GetEditor()),
			regexp.MustCompile("Info|HandleInput|CommandHandler|Windows|Frontend|Console|SetActiveWindow|Init|Watch|Observe|SetClipboardFuncs|KeyBindings").MatchString,
			"backend.GetEditor().",
			sn),
		},
	}
	data[len(data)-1][1] += fmt.Sprintf(`var sublime_methods = []py.Method{
		%s
	}`, sublime_methods)
	var year = strconv.FormatInt(int64(time.Now().Year()), 10)

	for _, gen := range data {
		if gen[0] == "" {
			continue
		}
		wr := `// Copyright ` + year + ` The lime Authors.
			// Use of this source code is governed by a 2-clause
			// BSD-style license that can be found in the LICENSE file.

			// This file was generated by tasks/build/gen_python_api.go and shouldn't be manually modified

			package sublime

			import (
				"fmt"
				"github.com/limetext/gopy/lib"
				"github.com/limetext/lime/backend"
				"github.com/limetext/text"
			)
			var (
				_ = backend.View{}
				_ = text.Region{}
				_ = fmt.Errorf
			)
			` + gen[1]
		if err := ioutil.WriteFile(gen[0], []byte(wr), 0644); err != nil {
			panic(err)
		} else {
			c := exec.Command("go", "fmt", gen[0])
			if o, err := c.CombinedOutput(); err != nil {
				panic(fmt.Errorf("%s, %s", o, err))
			} else {
				fmt.Printf("%s", string(o))
			}
		}
	}
}
