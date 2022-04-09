package types

import "strings"

var (
	goBuiltinTypes = map[string]bool{"int": true, "int32": true, "int64": true,
		"float": true, "float32": true, "float64": true,
		"uint": true, "uint32": true, "uint64": true,
		"string": true, "bool": true}
)

type SourceInfo struct {
	FileName string
}

type Tag struct {
	Json    string
	Bson    string
	Binding bool
	Form    string
}

type TypeD struct {
	Kind string
}

func (o *TypeD) IsArray() bool {
	return strings.HasPrefix(o.Kind, "[]")
}

func (o *TypeD) IsMap() bool {
	return strings.HasPrefix(o.Kind, "map")
}

func (o *TypeD) IsComposedByBuiltin() bool {
	if goBuiltinTypes[o.Kind] {
		return true
	}
	if o.IsArray() && goBuiltinTypes[strings.Trim(o.Kind, "[]")] {
		return true
	}
	if o.IsMap() {
		strs := strings.Split(strings.Trim(o.Kind, "map["), "]")
		return len(strs) == 2 && goBuiltinTypes[strs[0]] && goBuiltinTypes[strs[1]]
	}
	return false
}

func (o *TypeD) GetKind() string {
	return strings.Trim(o.Kind, "[]")
}

func (o *TypeD) GetMapKind() []string {
	return strings.Split(strings.Trim(o.Kind, "map["), "]")
}

type Field struct {
	Name     string
	Tag      *Tag
	Kind     TypeD
	Comments string
}

type StructRecord struct {
	Name   string
	Fields []Field
	SInfo  SourceInfo
}
