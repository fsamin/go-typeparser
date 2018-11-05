package typeparser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"strings"
)

type Type struct {
	file    *ast.File
	docs    []string
	spec    *ast.TypeSpec
	strct   *ast.StructType
	iface   *ast.InterfaceType
	fields  []Field
	methods []Method
}

func (t Type) IsConcrete() bool {
	return t.strct != nil
}

func (t Type) IsInterface() bool {
	return t.iface != nil
}

func (t Type) Name() string {
	return fmt.Sprintf("%v", t.spec.Name)
}

func (t Type) Docs() List {
	return List(t.docs)
}

func (t Type) Fields() []Field {
	return t.fields
}

func (t Type) FieldNames() List {
	names := []string{}
	for _, f := range t.fields {
		names = append(names, f.field.Names[0].Name)
	}
	return List(names)
}

func (t Type) Field(name string) *Field {
	for _, f := range t.fields {
		if f.Name() == name {
			return &f
		}
	}
	return nil
}

func (t Type) Methods() []Method {
	return t.methods
}

func (t Type) MethodNames() List {
	names := []string{}
	for _, m := range t.methods {
		names = append(names, m.Name())
	}
	return List(names)
}

func (t Type) Method(s string) *Method {
	for _, m := range t.methods {
		if m.Name() == s {
			return &m
		}
	}
	return nil
}

type Method struct {
	method *ast.Field
}

func (m Method) Name() string {
	return m.method.Names[0].Name
}

func (m Method) ParamNames() List {
	names := []string{}
	for _, p := range m.Params() {
		names = append(names, p.Name())
	}
	return List(names)
}

func (m Method) ParamTypes() List {
	types := []string{}
	for _, p := range m.Params() {
		types = append(types, p.Type())
	}
	return List(types)
}

func (m Method) Params() []Param {
	f, ok := m.method.Type.(*ast.FuncType)
	if !ok {
		return nil
	}

	params := []Param{}
	for _, p := range f.Params.List {
		params = append(params, Param{
			param: p,
		})
	}

	return params
}

func (m Method) ResultTypes() List {
	types := []string{}
	for _, p := range m.Results() {
		types = append(types, p.Type())
	}
	return List(types)
}

func (m Method) Results() []Param {
	f, ok := m.method.Type.(*ast.FuncType)
	if !ok {
		return nil
	}

	params := []Param{}
	for _, p := range f.Results.List {
		params = append(params, Param{
			param: p,
		})
	}

	return params
}

type Param struct {
	param *ast.Field
}

func (p Param) Name() string {
	if len(p.param.Names) == 0 {
		return ""
	}
	return p.param.Names[0].Name
}

func (p Param) Type() string {
	switch t := p.param.Type.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.Ellipsis:
		return "..." + fmt.Sprintf("%v", t.Elt)
	default:
		return fmt.Sprintf("%v", p.param.Type)
	}
}

type Field struct {
	field  *ast.Field
	goType string
}

func (f Field) Tags() List {
	if f.field.Tag == nil {
		return List{}
	}
	tags := strings.Replace(f.field.Tag.Value, "`", "", 2)
	return List(strings.Split(tags, " "))
}

func (f Field) TagValue(s string) List {
	return f.Tags().
		Map(func(t string) string {
			if !strings.HasPrefix(t, s+":") {
				return ""
			}
			val := strings.Replace(t, s+":", "", 1)
			val = strings.Replace(val, "\"", "", 2)
			return val
		}).
		Filter(func(s string) bool {
			if s == "" {
				return false
			}
			return true
		}).
		Explode(func(s string) []string {
			return strings.Split(s, ",")
		})
}

func (f Field) Name() string {
	return f.field.Names[0].Name
}

func (f Field) Type() string {
	return f.goType
}

func Parse(filename string) ([]Type, error) {
	src, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", src, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	types := []Type{}

	for _, d := range f.Decls {
		t := Type{
			file: f,
		}

		gen, ok := d.(*ast.GenDecl)
		if !ok {
			continue
		}

		if len(gen.Specs) != 1 {
			continue
		}

		typeSpec, ok := gen.Specs[0].(*ast.TypeSpec)
		if !ok {
			continue
		}

		for _, c := range gen.Doc.List {
			t.docs = append(t.docs, c.Text)
		}

		t.spec = typeSpec

		t.strct, ok = typeSpec.Type.(*ast.StructType)
		if ok {
			for _, f := range t.strct.Fields.List {
				start, end := f.Type.Pos()-1, f.Type.End()-1
				t.fields = append(t.fields, Field{
					field:  f,
					goType: string(src[start:end]),
				})
			}
		} else {
			t.iface, ok = typeSpec.Type.(*ast.InterfaceType)
			if !ok {
				continue
			}
			for _, m := range t.iface.Methods.List {
				t.methods = append(t.methods, Method{
					method: m,
				})
			}
		}

		types = append(types, t)
	}
	return types, nil
}
