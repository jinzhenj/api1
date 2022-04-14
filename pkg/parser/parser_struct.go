package parser

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/go-logr/logr"

	"github.com/go-swagger/pkg/types"
	"github.com/go-swagger/pkg/utils"
)

type ParserStructFile struct {
	log              logr.Logger
	FileName         string
	ModulePrefixName string
}

func NewParserStructFile(log logr.Logger, FileName string) *ParserStructFile {
	return &ParserStructFile{log: log, FileName: FileName, ModulePrefixName: strings.ReplaceAll(filepath.Dir(FileName), "/", ".")}
}

func ParserStructFromDirs(log logr.Logger, dirs []string) ([]types.StructRecord, error) {
	fliter := func(fn string) bool {
		return strings.HasSuffix(fn, ".go")
	}

	files := make([]string, 0)

	for _, dir := range dirs {
		cfiles, err := utils.ListFiles(dir, fliter)
		if err != nil {
			return nil, err
		}
		files = append(files, cfiles...)
	}

	ret := make([]types.StructRecord, 0)
	for _, fileName := range files {
		pFileObj := ParserStructFile{log: log, FileName: fileName, ModulePrefixName: strings.ReplaceAll(filepath.Dir(fileName), "/", ".")}
		records, err := pFileObj.ParserStructFromFile()
		if err != nil {
			return nil, err
		}
		ret = append(ret, records...)
	}

	return ret, nil
}

func (o ParserStructFile) ParserStructFromFile() ([]types.StructRecord, error) {
	byteData, err := ioutil.ReadFile(o.FileName)
	if err != nil {
		return nil, err
	}

	mayHandlerStrContent, err := extractBracesBlock(o.log, bytes.NewReader([]byte(byteData)), reStructPrefix)
	if err != nil {
		o.log.Error(err, "failed to parse go struct when extract block", "content", byteData)
		return nil, err
	}
	o.log.V(6).Info("parser_struct", "found blocks", mayHandlerStrContent)

	ret := make([]types.StructRecord, 0)
	for _, content := range mayHandlerStrContent {
		record := o.ParseStruct(content)
		if record != nil {
			record.SInfo = types.SourceInfo{FileName: o.FileName}
			ret = append(ret, *record)
		}
	}
	return ret, nil
}

// real work function
func (o ParserStructFile) ParseStruct(s string) *types.StructRecord {
	var ret types.StructRecord
	fields := make([]types.Field, 0)
	contents := strings.Split(s, "\n")

	ret.Name = reToken.FindAllString(contents[0], -1)[1]
	ret.RelativePathName = o.ModulePrefixName + "." + ret.Name

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
			o.parserField(line, &field)
			fields = append(fields, field)
			comments = make([]string, 0)
		}
	}
	ret.Fields = fields
	return &ret
}

// underlayer functions
func (o ParserStructFile) parserField(s string, field *types.Field) {
	ns := strings.TrimLeft(usingWhiteSpace(s), " ")
	strs := reHandlerToken.FindAllString(ns, -1)
	field.Name = strs[0]
	field.Kind = types.TypeD{Kind: mayAddPathToStructKind(o.ModulePrefixName, strs[1])}

	tagContent := reTag.FindString(s)
	if len(tagContent) > 0 {
		field.Tag = o.parserTag(tagContent)
		if field.Tag.Json == "" {
			field.Tag.Json = field.Name
		}
	} else {
		field.Tag = &types.Tag{Json: field.Name}
	}
}

//TODO: use regexp
func (o ParserStructFile) parserTag(s string) *types.Tag {
	var ret types.Tag
	for _, str := range reHandlerToken.FindAllString(strings.Trim(s, "`"), -1) {
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
		case "position":
			ret.Position = trimOmitempty(keys[1])
		case "form":
			ret.Form = trimOmitempty(keys[1])
		case "bson":
			ret.Bson = trimOmitempty(keys[1])
		default:

		}
	}
	return &ret
}
