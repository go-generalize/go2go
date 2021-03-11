package go2go

import (
	"go/types"
	"strings"
)

func getBasicTypeName(k types.BasicKind) string {
	switch k {
	case types.Bool:
		return "bool"
	case types.Int:
		return "int"
	case types.Int8:
		return "int8"
	case types.Int16:
		return "int16"
	case types.Int32:
		return "int32"
	case types.Int64:
		return "int64"
	case types.Uint:
		return "uint"
	case types.Uint8:
		return "uint8"
	case types.Uint16:
		return "uint16"
	case types.Uint32:
		return "uint32"
	case types.Uint64:
		return "uint64"
	case types.Uintptr:
		return "uintptr"
	case types.Float32:
		return "float32"
	case types.Float64:
		return "float64"
	case types.Complex64:
		return "complex64"
	case types.Complex128:
		return "complex128"
	case types.String:
		return "string"
	default:
		return "interface{}" // Unsupported type
	}
}

// SplitPackegeStruct - package.structを分割
func SplitPackegeStruct(s string) (string, string) {
	idx := strings.LastIndex(s, ".")

	return s[:idx], s[idx+1:]
}
