package go2go

import (
	"fmt"
	"strings"
	"testing"

	"github.com/go-generalize/go2ts/pkg/parser"
	"github.com/k0kubun/pp"
)

func TestGenerate(t *testing.T) {
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

	psr, err := parser.NewParser("./test", filter)

	if err != nil {
		t.Fatal(err)
	}

	res, err := psr.Parse()

	if err != nil {
		t.Fatal(err)
	}

	gen := NewGenerator(res)
	fmt.Println(gen.Generate())

	pp.Println(gen.generatorParam.Objects)
	pp.Println(gen.generatorParam.Consts)
}
