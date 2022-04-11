package parser

import (
	"bytes"
	"strings"

	"github.com/go-logr/logr"

	"github.com/go-swagger/pkg/types"
)

const (
	retDef              = "return"
	maxEmbedStructDepth = 5
)

func ParserApiDef(log logr.Logger, content string) ([]types.HttpHandler, error) {
	ret := make([]types.HttpHandler, 0)
	mayHandlerStrContent, err := extractBracesBlock(log, bytes.NewReader([]byte(content)), reSvcDefPrefix)
	if err != nil {
		log.Error(err, "failed to parse api def when extract block", "content", content)
		return nil, err
	}
	log.V(6).Info("parser_api_def", "found blocks", mayHandlerStrContent)
	for _, line := range mayHandlerStrContent {
		noBraceContents := reMultiLineBracesContent.FindAllString(line, -1)
		log.V(6).Info("parser_api_def", "no remove braces", noBraceContents, "line", line)
		if len(noBraceContents) == 0 {
			continue
		}
		handler, err := parserHandler(log, strings.Trim(noBraceContents[0], " \n\t{}"))
		if err != nil {
			log.Error(err, "failed to parse hanlder", "content", line)
			return nil, err
		}
		ret = append(ret, *handler)
	}
	return ret, nil
}

func parserHandler(log logr.Logger, s string) (*types.HttpHandler, error) {
	var ret types.HttpHandler
	ns := strings.Trim(usingWhiteSpace(s), " \n\t")
	log.V(8).Info("parser Handler", "input", s, "normalized input", ns)

	// may get doc

	docCotents := reHandlerDoc.FindAllString(ns, -1)
	if len(docCotents) != 0 {
		ret.Doc = parserDoc(log, docCotents[0])
		log.V(8).Info("parser handler", "doc", docCotents)
	}

	ns = reHandlerDoc.ReplaceAllString(ns, "")

	// get handler
	handlerContents := reHanderHander.FindAllString(ns, 1)
	if len(handlerContents) != 1 {
		log.Error(ErrMultiHandlerFound, "extract handler def", "handlers", handlerContents, "params", ns)
	}
	handlerTokens := reToken.FindAllString(handlerContents[0], -1)
	ret.Name = handlerTokens[1]

	ns = reHanderHander.ReplaceAllString(ns, "")

	// get the rest info of handler, parser `get /api/user/search (pkg.config.UserSearchReq) returns (pkg.config.UserInfoReply)`
	for _, line := range strings.Split(ns, "\n") {
		if len(reEmptyLineWithSpace.FindAllString(line, -1)) > 0 {
			continue
		}
		parseHanlderInfo(log, line, &ret)
		break
	}
	return &ret, nil
}

// underlayer functions
func parserDoc(log logr.Logger, s string) *types.HandlerDoc {
	log.V(6).Info("parse doc", "source string", s)
	var ret types.HandlerDoc
	ns := strings.Trim(usingWhiteSpace(s), " \n\t")

	strs := strings.Split(ns, "\n")
	if !strings.Contains(strs[0], "@doc") {
		log.Info("no doc found when trying parse", "content", strs[0], "source", s)
		return nil
	}

	for i, line := range strs[1:] {
		log.V(6).Info("parser doc", "line_no", i, "content", line)
		if i == len(strs)-2 {
			break
		}
		matched := regUnicodeStr.FindAllString(line, -1)
		if len(matched) == 0 {
			log.Info("no token found when trying parse doc", "content", line)
			continue
		}
		log.V(6).Info("matched doc", "records", matched)
		switch matched[0] {
		case string(types.HandlerDocAttrSummary):
			ret.Summary = strings.Join(matched[1:], " ")
		case string(types.HandlerDocAttrDescription):
			ret.Description = strings.Join(matched[1:], " ")
		default:
			log.Info("unkonwn content of doc", "content", line)
		}
	}

	return &ret
}

func parseHanlderInfo(log logr.Logger, s string, handler *types.HttpHandler) error {
	tokens := reHandlerToken.FindAllString(strings.Trim(s, " \t\n"), -1)
	log.V(6).Info("parser handler info", "tokens", tokens)
	if 2 > len(tokens) {
		return ErrInvalidHttpHandlerDef
	}
	handler.Method = tokens[0]
	handler.Endpoint = tokens[1]

	if len(tokens) > 2 {
		if tokens[2] == retDef {
			handler.Res = parseHandlerBodyMethod(tokens[3])
		} else {
			handler.Req = parseHandlerBodyMethod(tokens[2])
			if len(tokens) == 5 { // why is 5? len([method, url, req, return, res]) == 5
				handler.Res = parseHandlerBodyMethod(tokens[4])
			}
		}
	}
	return nil
}

// (pkg.config.UserSearchReq{list=[]string, records=[]int})
func parseHandlerBodyMethod(s string) *types.HandlerBodyParams {
	var ret types.HandlerBodyParams
	str := strings.Trim(s, "()")
	namesMap := make(map[string]bool)
	ret.Name = reBracesContent.ReplaceAllString(str, "")

	params := []string{str}

	for depth := 0; depth < maxEmbedStructDepth; depth++ {
		if len(params) == 0 {
			break
		}
		tmp := make([]string, 0)

		for _, line := range params {
			newName, m := extractNestedReplacedStruct(line)
			namesMap[newName] = true

			for _, val := range m {
				tmp = append(tmp, val)
			}
		}
		params = tmp
	}
	ret.RelatedNames = namesMap
	return &ret
}

func extractNestedReplacedStruct(s string) (string, map[string]string) {
	idx := strings.Index(s, "{")
	if idx == -1 {
		return s, nil
	}
	m := make(map[string]string)
	for _, line := range strings.Split(strings.Trim(s[idx+1:], "}"), ",") {
		idxE := strings.Index(line, "=")
		if idxE >= 0 {
			key := trimSpace(line[:idxE])
			val := trimSpace(line[idxE+1:])
			m[key] = val
		}
	}

	return trimSpace(s[:idx]), m
}