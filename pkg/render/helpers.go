package render

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/go-swagger/pkg/types"
)

func foundStructDef(name string, arr []types.StructRecord) *types.StructRecord {
	for _, record := range arr {
		if name == GoStructName2SvcDefName(&record) {
			return &record
		}
	}
	return nil
}

//TODO: 和 utils 里面判断 goBuiltin 抽取出来集合一下
func mapGoTypesToSwagger(kind string) string {
	switch kind {
	case "int", "int32", "int64":
		return "integer"
	case "float", "float32", "float64":
		return "number"
	case "bool":
		return "boolean"
	default:
		fmt.Printf("not supported go types:(%s) when convert to swagger type\n", kind)
		return kind
	}
}

func GoStructName2SvcDefName(o *types.StructRecord) string {
	if o == nil {
		return ""
	}
	return strings.Replace(filepath.Join(filepath.Dir(o.SInfo.FileName), o.Name), "/", ".", -1)
}
