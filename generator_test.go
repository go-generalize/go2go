package go2go

import (
	"fmt"
	"strings"
	"testing"

	"github.com/go-generalize/go2ts/pkg/parser"
	tstypes "github.com/go-generalize/go2ts/pkg/types"
	"github.com/k0kubun/pp"
)

func TestGenerate_Standard(t *testing.T) {
	filter := func(opt *parser.FilterOpt) bool {
		if opt.Dependency {
			return true
		}
		if !opt.BasePackage {
			return false
		}
		if !opt.Exported {
			return false
		}

		return strings.HasSuffix(opt.Name, "Request") || strings.HasSuffix(opt.Name, "Response")
	}

	psr, err := parser.NewParser("./testdata/standard", filter)

	if err != nil {
		t.Fatal(err)
	}

	res, err := psr.Parse()

	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)

	gen := NewGenerator(res, []string{})
	fmt.Println(gen.Generate())

	pp.Println(gen.generatorParam.Objects)
	pp.Println(gen.generatorParam.Consts)
}

func TestGenerate_External(t *testing.T) {
	filter := func(opt *parser.FilterOpt) bool {
		if opt.Dependency {
			return true
		}
		if !opt.BasePackage {
			return false
		}
		if !opt.Exported {
			return false
		}

		return strings.HasSuffix(opt.Name, "Request") || strings.HasSuffix(opt.Name, "Response")
	}

	psr, err := parser.NewParser("./testdata/external", filter)

	if err != nil {
		t.Fatal(err)
	}

	res, err := psr.Parse()

	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)

	gen := NewGenerator(res, []string{})

	gen.ExternalGenerator = func(t tstypes.Type) (*GeneratedType, bool) {
		_, ok := t.(*tstypes.String)

		if !ok {
			return nil, false
		}

		return &GeneratedType{
			Path: "github.com/go-generalize/go2ts/testdata/external/external",
			Name: "MyString",
		}, true
	}

	fmt.Println(gen.Generate())

	pp.Println(gen.generatorParam.Objects)
	pp.Println(gen.generatorParam.Consts)
}
