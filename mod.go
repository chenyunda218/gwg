package gwg

import (
	"fmt"
	"os"
	"strings"
)

type Code interface {
	String() string
}

type Mod struct {
	Module    string
	Go        string
	Gitignore string
	Package   Package
}

type Package struct {
	Name    string
	Imports []Import
	Codes   []Code
}

type Interface struct {
	Name    string
	Methods []Method
}

func (i Interface) String() (c string) {
	c = fmt.Sprintf("type %s interface {\n", i.Name)
	for _, m := range i.Methods {
		c += fmt.Sprintf("%s\n", m.Prototype())
	}
	c += "}\n"
	return c
}

func (i *Interface) AddMethod(ms ...Method) Interface {
	i.Methods = append(i.Methods, ms...)
	return *i
}

func (p Package) Wirte(path string) {
	os.MkdirAll(path, os.ModePerm)
	f, _ := os.Create(fmt.Sprintf("%s/%s.go", path, p.Name))
	defer f.Close()
	f.Write([]byte(p.String()))
}

func (p *Package) AddImport(is ...Import) Package {
	p.Imports = append(p.Imports, is...)
	return *p
}

func (p *Package) AddCode(c ...Code) Package {
	p.Codes = append(p.Codes, c...)
	return *p
}

func (p Package) String() (c string) {
	c = fmt.Sprintf("package %s\n", p.Name)
	for _, i := range p.Imports {
		c += i.String()
	}
	for _, code := range p.Codes {
		c += code.String()
	}
	return c
}

type Import struct {
	Packages []string
}

func (i Import) String() (c string) {
	for _, p := range i.Packages {
		c += fmt.Sprintf("import \"%s\"\n", p)
	}
	return c
}

func (i *Import) Add(is ...string) Import {
	i.Packages = append(i.Packages, is...)
	return *i
}

type Pair struct {
	Left  string
	Right string
}

func (p Pair) String() string {
	if p.Left == "" {
		return p.Right
	}
	return fmt.Sprint(p.Left, " ", p.Right)
}

func (p *Pair) SLeft(left string) Pair {
	p.Left = left
	return *p
}
func (p *Pair) SRight(right string) Pair {
	p.Right = right
	return *p
}

type Parameters struct {
	Pairs []Pair
}

func (p Parameters) String() string {
	var ps []string
	for _, p := range p.Pairs {
		ps = append(ps, p.String())
	}
	return strings.Join(ps, ", ")
}

func (p *Parameters) Add(ps ...Pair) Parameters {
	p.Pairs = append(p.Pairs, ps...)
	return *p
}

type Outputs struct {
	Pairs []Pair
}

func (o Outputs) String() (c string) {
	if len(o.Pairs) == 0 {
		return ""
	} else if len(o.Pairs) == 1 {
		if o.Pairs[0].Left == "" {
			return o.Pairs[0].Right
		}
		return fmt.Sprintf("(%s)", o.Pairs[0].String())
	}
	var s []string
	for _, p := range o.Pairs {
		s = append(s, p.String())
	}
	return fmt.Sprintf("(%s)", strings.Join(s, ", "))
}

type Tag struct {
	Label   string
	Content string
}

func (t Tag) String() string {
	return fmt.Sprintf("%s:\"%s\"", t.Label, t.Content)
}

func NewProperty(label, t string, tags ...Tag) Property {
	return Property{
		Label: label,
		Type:  t,
		Tags:  tags,
	}
}

type Property struct {
	Label string
	Type  string
	Tags  []Tag
}

func (p Property) String() string {
	var tags []string
	if len(p.Tags) != 0 {
		for _, t := range p.Tags {
			tags = append(tags, t.String())
		}
	}
	var tag string
	if len(tags) > 0 {
		tag = "`" + strings.Join(tags, " ") + "`"
	}
	return fmt.Sprintf("%s %s %s", p.Label, p.Type, tag)
}

type Func struct {
	Name       string
	Parameters Parameters
	Outputs    Outputs
	Lines      []Line
}

func (f Func) Call(args ...interface{}) Code {
	var argsString []string
	for _, arg := range args {
		switch arg.(type) {
		case float32:
		case float64:
		case int64:
		case int32:
		case int:
			argsString = append(argsString, fmt.Sprint(arg))
		case string:
			argsString = append(argsString, fmt.Sprintf("\"%s\"", arg))
		}
	}
	return Line{Content: fmt.Sprintf("%s(%s)\n", f.Name, strings.Join(argsString, ", "))}
}

func (f Func) String() (c string) {
	c = fmt.Sprintf("func %s(%s) %s {\n", f.Name, f.Parameters.String(), f.Outputs.String())
	for _, line := range f.Lines {
		c += fmt.Sprintf("%s\n", line.String())
	}
	c += fmt.Sprintln("}")
	return c
}

func (f *Func) AddLine(l ...Line) Func {
	f.Lines = append(f.Lines, l...)
	return *f
}

func (f *Func) AddParameter(ps ...Pair) Func {
	f.Parameters.Add(ps...)
	return *f
}

type Method struct {
	Ref        bool
	Name       string
	Struct     *Struct
	Parameters Parameters
	Outputs    Outputs
}

func (m *Method) Prototype() string {
	var outputs string
	if len(m.Outputs.Pairs) != 0 {
		outputs = " (" + m.Outputs.String() + ")"
	}
	return fmt.Sprintf("%s(%s)%s", m.Name, m.Parameters.String(), outputs)
}

type Struct struct {
	Name         string
	Combinations []string
	Properties   []Property
	Methods      map[string]Method
}

func (s Struct) String() (c string) {
	c = fmt.Sprintln("type", s.Name, "struct {")
	for _, com := range s.Combinations {
		c += fmt.Sprintln(com)
	}
	for _, p := range s.Properties {
		c += fmt.Sprintln(p)
	}
	c += "}\n"
	return c
}

func (s *Struct) AddMethod(methods ...Method) Struct {
	if s.Methods == nil {
		s.Methods = make(map[string]Method)
	}
	for i, m := range methods {
		method := methods[i]
		method.Struct = s
		s.Methods[m.Name] = method
	}
	return *s
}

func (s *Struct) AddProperty(ps ...Property) Struct {
	s.Properties = append(s.Properties, ps...)
	return *s
}

func (s *Struct) AddCombination(c ...string) Struct {
	s.Combinations = append(s.Combinations, c...)
	return *s
}

type Line struct {
	Content string
}

func (l Line) String() string {
	return l.Content
}

type Enums struct {
	Title  string
	Values []string
}

func (e *Enums) Add(values ...string) Enums {
	e.Values = append(e.Values, values...)
	return *e
}

func (e Enums) String() (c string) {
	c = fmt.Sprintf("type %s string\n", e.Title)
	for _, v := range e.Values {
		c += fmt.Sprintf("const %s %s = \"%s\"\n", v, e.Title, v)
	}
	return c
}

type Good string

const HELLO Good = "HELLO"
