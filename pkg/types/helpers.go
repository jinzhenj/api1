package types

import (
	"strings"
)

var (
	goBuiltinTypes = map[string]bool{"int": true, "int32": true, "int64": true,
		"float": true, "float32": true, "float64": true,
		"uint": true, "uint32": true, "uint64": true,
		"string": true, "bool": true,
		"interface{}": true,
	}
)

func IsGoBuiltinTypes(s string) bool {
	return goBuiltinTypes[s]
}

func IsInterface(s string) bool {
	return s == "interface{}"
}

func IsComposedByBuiltin(s string) bool {
	if goBuiltinTypes[s] {
		return true
	}
	if strings.HasPrefix(s, "[]") && goBuiltinTypes[strings.Trim(s, "[]")] {
		return true
	}
	if strings.HasPrefix(s, "map[") {
		strs := strings.Split(strings.Trim(s, "map["), "]")
		return len(strs) == 2 && goBuiltinTypes[strs[0]] && goBuiltinTypes[strs[1]]
	}
	return false
}
