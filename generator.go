package go2go

import (
	"bytes"
	"crypto/sha1"
	_ "embed"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"text/template"

	tstypes "github.com/go-generalize/go2ts/pkg/types"
	"golang.org/x/xerrors"
)

type Generator struct {
	types map[string]tstypes.Type
	generatorParam

	converted   map[string]string
	prereserved map[string]string
	reserved    map[string]struct{}
}

type objectEntry struct {
	Field string
	Type  string
	Tag   string
}

type object struct {
	Name   string
	Fields []objectEntry
}

type constantEnum struct {
	Name  string
	Value string
}

type constant struct {
	Name  string
	Base  string
	Enums []constantEnum
}

type generatorParam struct {
	Consts  []constant
	Objects []object

	UseTimePackage bool
}

type metadata struct {
	upperStructName string
	inlineIndex     int
}

func NewGenerator(types map[string]tstypes.Type) *Generator {
	return &Generator{
		types: types,

		converted:   map[string]string{},
		reserved:    map[string]struct{}{},
		prereserved: map[string]string{},
	}
}

func (g *Generator) convert(v tstypes.Type, meta *metadata) string {
	switch v := v.(type) {
	case *tstypes.Array:
		return "[]" + g.convert(v.Inner, meta)
	case *tstypes.Object:
		return g.convertObject(v, meta)
	case *tstypes.String:
		return g.convertString(v, meta)
	case *tstypes.Number:
		return g.convertNumber(v, meta)
	case *tstypes.Boolean:
		return "bool"
	case *tstypes.Date:
		g.UseTimePackage = true
		return "time.Time"
	case *tstypes.Nullable:
		return "*" + g.convert(v.Inner, meta)
	case *tstypes.Any:
		return "interface{}"
	case *tstypes.Map:
		return fmt.Sprintf("map[%s]%s", g.convert(v.Key, meta), g.convert(v.Value, meta))
	default:
		panic("unsupported")
	}
}

func (g *Generator) convertString(str *tstypes.String, upper *metadata) string {
	if len(str.Enum) == 0 {
		return "string"
	}

	if name, ok := g.converted[str.Name]; ok {
		return name
	}

	name := g.getConvertedType(str.Name, upper)
	consts := make([]constantEnum, 0, len(str.RawEnum))

	_, orig := SplitPackegeStruct(str.Name)

	for _, e := range str.RawEnum {
		key := name + e.Key
		if strings.HasPrefix(e.Key, orig) {
			key = name + strings.TrimPrefix(e.Key, orig)
		}

		consts = append(consts, constantEnum{
			Name:  key,
			Value: strconv.Quote(e.Value),
		})
	}

	g.Consts = append(g.Consts, constant{
		Name:  name,
		Base:  "string",
		Enums: consts,
	})

	return name
}

func (g *Generator) convertNumber(num *tstypes.Number, upper *metadata) string {
	if len(num.Enum) == 0 {
		return getBasicTypeName(num.RawType)
	}

	if name, ok := g.converted[num.Name]; ok {
		return name
	}

	name := g.getConvertedType(num.Name, upper)

	enums := make([]constantEnum, 0, len(num.RawEnum))
	_, orig := SplitPackegeStruct(num.Name)

	for _, e := range num.RawEnum {
		key := name + e.Key
		if strings.HasPrefix(e.Key, orig) {
			key = name + strings.TrimPrefix(e.Key, orig)
		}

		enums = append(enums, constantEnum{
			Name:  key,
			Value: fmt.Sprint(e.Value),
		}) // Support multiple types
	}

	g.Consts = append(g.Consts, constant{
		Name:  name,
		Base:  getBasicTypeName(num.RawType),
		Enums: enums,
	})

	return name
}

func (g *Generator) convertObject(obj *tstypes.Object, upper *metadata) string {
	var converted object

	if name, ok := g.converted[obj.Name]; ok {
		return name
	}

	name := g.getConvertedType(obj.Name, upper)

	entries := make([]tstypes.ObjectEntry, 0, len(obj.Entries))
	for _, v := range obj.Entries {
		entries = append(entries, v)
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].RawName < entries[j].RawName
	})

	for i, e := range entries {
		converted.Fields = append(converted.Fields, objectEntry{
			Field: e.RawName,
			Type:  g.convert(e.Type, &metadata{upperStructName: name, inlineIndex: i}),
		})
	}

	converted.Name = name
	g.Objects = append(g.Objects, converted)

	return name
}

func (g *Generator) getConvertedType(fullName string, meta *metadata) string {
	var name string

	if fullName == "" {
		name = meta.upperStructName + "Inline" + fmt.Sprintf("%03d", meta.inlineIndex)
	} else {
		_, name = SplitPackegeStruct(fullName)

		prev, prereserved := g.prereserved[fullName]
		_, reserved := g.reserved[name]
		if (prereserved && prev != fullName) || reserved {
			hash := fmt.Sprintf("%x", sha1.Sum([]byte(fullName)))

			name = name + "_" + hash[:4]
		}

		g.reserved[name] = struct{}{}
	}
	g.converted[fullName] = name

	return name
}

//go:embed templates/types.go.tmpl
var templateBase string

func (g *Generator) Generate() (string, error) {
	for _, v := range g.types {
		g.convert(v, nil)
	}

	sort.Slice(g.Objects, func(i, j int) bool {
		return g.Objects[i].Name < g.Objects[j].Name
	})
	sort.Slice(g.Consts, func(i, j int) bool {
		return g.Consts[i].Name < g.Consts[j].Name
	})

	tmpl := template.Must(template.New("").Parse(templateBase))

	buf := bytes.NewBuffer(nil)
	if err := tmpl.Execute(buf, g.generatorParam); err != nil {
		return "", xerrors.Errorf("failed to generate template: %w", err)
	}

	return buf.String(), nil
}
