package types

import "strings"

type SourceInfo struct {
	FileName string
}

type Tag struct {
	Json     string
	Bson     string
	Position string
	Binding  bool
	Form     string
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

func (o *TypeD) GetKind() string {
	return strings.Trim(o.Kind, "[]")
}

func (o *TypeD) GetMapKind() []string {
	return strings.Split(strings.Trim(o.Kind, "map["), "]")
}

type Field struct {
	Name     string `json:"name,omitempty"`
	Tag      *Tag   `json:"tag,omitempty"`
	Kind     TypeD  `json:"kind,omitempty"`
	Comments string `json:"comments,omitempty"`
}

type StructRecord struct {
	Name   string
	Fields []Field
	SInfo  SourceInfo
}
