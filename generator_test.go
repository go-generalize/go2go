package go2go

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/go-generalize/go2ts/pkg/parser"
	tstypes "github.com/go-generalize/go2ts/pkg/types"
	"github.com/google/go-cmp/cmp"
	"github.com/k0kubun/pp"
)

func TestGenerate(t *testing.T) {
	type testCase struct {
		dir               string
		expected          string
		externalGenerator func(tstypes.Type) (*GeneratedType, bool)
	}

	cases := []testCase{
		{
			dir:      "testdata/standard",
			expected: "testdata/standard/expected/expected.go",
		},
		{
			dir:      "testdata/external",
			expected: "testdata/external/expected/expected.go",
			externalGenerator: func(t tstypes.Type) (*GeneratedType, bool) {
				_, ok := t.(*tstypes.String)

				if !ok {
					return nil, false
				}

				return &GeneratedType{
					Path: "github.com/go-generalize/go2ts/testdata/external/external",
					Name: "MyString",
				}, true
			},
		},
	}

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

	for _, c := range cases {
		c := c
		t.Run(c.dir, func(t *testing.T) {
			psr, err := parser.NewParser(c.dir, filter)

			if err != nil {
				t.Fatal(err)
			}

			res, err := psr.Parse()

			if err != nil {
				t.Fatal(err)
			}
			fmt.Println(res)

			gen := NewGenerator(res, []string{})
			gen.ExternalGenerator = c.externalGenerator

			generated, err := gen.Generate()

			if err != nil {
				t.Fatal(err)
			}

			fmt.Println(generated)
			pp.Println(gen.generatorParam.Objects)
			pp.Println(gen.generatorParam.Consts)

			b, err := os.ReadFile(c.expected)

			if err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(string(b), generated); diff != "" {
				t.Error(diff)
			}

		})
	}
}
