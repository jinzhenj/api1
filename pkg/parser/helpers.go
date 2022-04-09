package parser

import (
	"errors"
	"io"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/go-swagger/pkg/types"
)

var (
	reStructPrefix   = regexp.MustCompile(`type\s+[a-zA-Z0-9]+\s+struct\s?{`)
	reToken          = regexp.MustCompile(`[a-zA-Z0-9]+`)
	reIsComment      = regexp.MustCompile(`^\s*//.*`)
	reExtractComment = regexp.MustCompile(`[^/]+`)
	reTag            = regexp.MustCompile("`.*`")
	//By default . does not match newlines. To change that, use the s flag.  https://go.dev/src/regexp/syntax/doc.go
	reHandlerDoc         = regexp.MustCompile(`(?sU)@doc\(.*\)`)
	reHanderHander       = regexp.MustCompile(`@handler.*`)
	regUnicodeStr        = regexp.MustCompile(`[^ \n\t]+`)
	reEmptyLineWithSpace = regexp.MustCompile(`^\s*$`)
	reHandlerToken       = regexp.MustCompile(`[^\s\t]+`)
	reBracesContent      = regexp.MustCompile(`{.*}`)

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
func extractBracesBlock(r io.Reader) ([]string, error) {
	content, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	ret := make([]string, 0)
	current := make([]string, 0)

	var extractStructDef bool
	var leftQ int

	for _, line := range strings.Split(string(content), "\n") {
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
			if reStructPrefix.FindAllString(line, 1) != nil {
				extractStructDef = true
				current = append(current, line)
				leftQ = 1
			}

		}
	}
	return ret, nil
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
