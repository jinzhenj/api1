package types

import (
	"regexp"
	"strings"

	"github.com/fatih/color"
)

var (
	reBrace = regexp.MustCompile(`{.*}`)
)

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
	Val string
}

func (o *TypeD) IsSupported() bool {
	return !o.IsMap()
}

func (o *TypeD) IsArray() bool {
	return strings.HasPrefix(o.Val, "[]")
}

func (o *TypeD) IsMap() bool {
	return strings.HasPrefix(o.Val, "map[")
}

func (o *TypeD) GetKind() string {
	if o.IsMap() {
		color.Red("not supported kind:(%s)\n", o.Val)
		return o.Val
	}
	if o.IsBuiltin() {
		return o.Tidy()
	}
	strs := strings.Split(o.Tidy(), ".")
	return strs[len(strs)-1]
}

func (o *TypeD) IsBuiltin() bool {
	tidyS := o.Tidy()
	return IsGoBuiltinTypes(tidyS)
}

func (o *TypeD) GetModule() string {
	if o.IsBuiltin() {
		return ""
	}
	strs := strings.Split(o.Tidy(), ".")
	return strings.Join(strs[:len(strs)-1], "/")
}

// []int ==> [int]
// map[string]bool ==> map[string]bool
// []pkg.types.ListRecords{data=xxxx} ==> pkg.types.ListRecords
func (o *TypeD) Tidy() string {
	ns := o.Val
	if o.IsArray() {
		ns = strings.Trim(o.Val, "[]")
	}
	if strings.Contains(ns, "{") {
		ns = reBrace.ReplaceAllString(ns, "")
	}
	return ns
}

func (o *TypeD) GetMapKind() []string {
	return strings.Split(strings.Trim(o.Val, "map["), "]")
}

type Field struct {
	Name     string `json:"name,omitempty"`
	Tag      *Tag   `json:"tag,omitempty"`
	Kind     TypeD  `json:"kind,omitempty"`
	Comments string `json:"comments,omitempty"`
}

type StructRecord struct {
	Name             string
	RelativePathName string
	Fields           []Field
	SInfo            SourceInfo
}
