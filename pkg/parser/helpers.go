package parser

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/go-logr/logr"
	"github.com/go-swagger/pkg/types"

	"github.com/go-swagger/pkg/utils"
)

// 正则表达式的名字不够好，不应该根据用途。而应该根据模式特点
var (
	reStructPrefix = regexp.MustCompile(`type\s+[a-zA-Z0-9]+\s+struct\s?{`)
	reSvcDefPrefix = regexp.MustCompile(`service\s+[a-zA-Z0-9]+\s{`)

	reToken          = regexp.MustCompile(`[a-zA-Z0-9]+`)
	reIsComment      = regexp.MustCompile(`^\s*//.*`)
	reExtractComment = regexp.MustCompile(`[^/]+`)
	reTag            = regexp.MustCompile("`.*`")
	//By default . does not match newlines. To change that, use the s flag.  https://go.dev/src/regexp/syntax/doc.go
	reHandlerDoc             = regexp.MustCompile(`(?sU)@doc\(.*\)`)
	reHanderHander           = regexp.MustCompile(`@handler.*`)
	regUnicodeStr            = regexp.MustCompile(`[^ \n\t]+`)
	reEmptyLineWithSpace     = regexp.MustCompile(`^\s*$`)
	reHandlerToken           = regexp.MustCompile(`[^\s\t]+`)
	reBracesContent          = regexp.MustCompile(`{.*}`)
	reMultiLineBracesContent = regexp.MustCompile(`(?sm){.*}`)
	reFoundModule            = regexp.MustCompile(`module\s+.*`)

	ErrMultiHandlerFound     = errors.New("multi handler def found")
	ErrInvalidHttpHandlerDef = errors.New("invalid http handler def")
)

// dir 相对路径
func extractStruct(dir string) ([]types.StructRecord, error) {

	ret := make([]types.StructRecord, 0)
	files, err := ioutil.ReadDir(".")
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		if strings.HasSuffix(f.Name(), ".go") {
			// ret = append(ret, tret...)
		}
	}
	return ret, nil
}

//从结构体抽取 struct 定义块
//必须保证注释中不包含字符 "{" 或 "}"
func extractBracesBlock(log logr.Logger, r io.Reader, pattern *regexp.Regexp) ([]string, error) {
	content, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	ret := make([]string, 0)
	current := make([]string, 0)

	var extractStructDef bool
	var leftQ int

	for _, line := range strings.Split(string(content), "\n") {
		log.V(8).Info("extractBracesBlock", "line", line, "leftQ", leftQ)
		if extractStructDef {
			current = append(current, line)
			leftQ += strings.Count(line, "{")
			leftQ -= strings.Count(line, "}")
			if leftQ == 0 {
				extractStructDef = false
				ret = append(ret, strings.Join(current, "\n"))
				current = make([]string, 0)
			}
		} else {
			if pattern.FindAllString(line, 1) != nil {
				extractStructDef = true
				current = append(current, line)
				leftQ = 1
			}

		}
	}
	return ret, nil
}

func extractHanderDefBlocks(s string) []string {
	ret := make([]string, 0)
	current := make([]string, 0)
	for _, line := range strings.Split(strings.Trim(s, " {}\n\t"), "\n") {
		if strings.Trim(line, " \t") == "" {
			if len(current) != 0 {
				ret = append(ret, strings.Join(current, "\n"))
				current = make([]string, 0)
			}
		} else {
			current = append(current, line)
		}
	}

	if len(current) != 0 {
		block := strings.Join(current, "\n")
		if strings.Contains(block, "@handler") {
			ret = append(ret, block)
		}
	}
	return ret
}

func isCommentLine(s string) bool {
	return len(reIsComment.FindAllString(s, -1)) > 0
}

func trimOmitempty(s string) string {
	return strings.Split(strings.Trim(s, `"`), ",")[0]
}

func usingWhiteSpace(s string) string {
	return strings.Replace(s, "\t", " ", -1)
}

func trimSpace(s string) string {
	return strings.Trim(s, " ")
}

func mayAddPathToStructKind(modulePrefixName, s string) string {

	proc := func(name string) string {
		if strings.Contains(name, ".") {
			return "pkg" + "." + name
		} else {
			return modulePrefixName + "." + name
		}
	}

	if utils.IsGoBuiltinTypes(s) {
		return s
	} else {
		if strings.HasPrefix(s, "[]") {
			ns := strings.Trim(s, "[]")
			if utils.IsGoBuiltinTypes(ns) {
				return s
			} else {
				return "[]" + proc(ns)
			}
		}
		if strings.HasPrefix(s, "map[") {
			fmt.Printf("not supported kind when add path to struct:(%s)", s)
			return s
		}

		return proc(s)
	}

}
