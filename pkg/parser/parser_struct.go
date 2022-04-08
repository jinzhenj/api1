package parser

import (
	"strings"

	"github.com/go-swagger/pkg/types"
)

func ParseStruct(relativePath, s string) *types.StructRecord {
	ret := types.StructRecord{SInfo: types.SourceInfo{FileName: relativePath}}
	fields := make([]types.Field, 0)
	contents := strings.Split(s, "\n")

	ret.Name = reStructName.FindAllString(contents[0], -1)[1]

	comments := make([]string, 0)
	for _, line := range contents[1:] {
		if strings.Contains(line, "}") {
			break
		}
		if isCommentLine(line) {
			comment := reExtractComment.FindAllString(line, -1)
			if len(comment) > 0 {
				comments = append(comments, comment[len(comment)-1])
			}
		} else {
			field := types.Field{Comments: strings.Join(comments, "\n")}
			parserField(line, &field)
			fields = append(fields, field)
			comments = make([]string, 0)
		}
	}
	ret.Fields = fields
	return &ret
}

// underlayer functions
func parserField(s string, field *types.Field) {
	ns := strings.TrimLeft(usingWhiteSpace(s), " ")
	strs := strings.Split(ns, " ")
	field.Name = strings.Trim(strs[0], " ")
	if strings.Contains(strs[1], "[]") {
		field.IsArray = true
	}
	field.Kind = strings.Trim(strs[1], " []")

	tagContent := reTag.FindString(s)
	if len(tagContent) > 0 {
		field.Tag = parserTag(tagContent)
		if field.Tag.Json == "" {
			field.Tag.Json = field.Name
		}
	} else {
		field.Tag = &types.Tag{Json: field.Name}
	}
}

func parserTag(s string) *types.Tag {
	var ret types.Tag
	for _, str := range strings.Split(strings.Trim(s, "`"), " ") {
		str := strings.Trim(str, " ")
		if len(str) == 0 {
			continue
		}
		keys := strings.Split(str, ":")

		if strings.Contains(keys[1], "-") {
			continue
		}
		switch keys[0] {
		case "json":
			ret.Json = trimOmitempty(keys[1])
		case "binding":
			if trimOmitempty(keys[1]) == "required" {
				ret.Binding = true
			}
		case "form":
			ret.Form = trimOmitempty(keys[1])
		case "bson":
			ret.Bson = trimOmitempty(keys[1])
		default:

		}
	}
	return &ret
}

func trimOmitempty(s string) string {
	return strings.Split(strings.Trim(s, `"`), ",")[0]
}

func usingWhiteSpace(s string) string {
	return strings.Replace(s, "\t", " ", -1)
}
