package parser

import (
	"strings"

	"github.com/go-logr/logr"

	"github.com/go-swagger/pkg/types"
)

const (
	retDef = "return"
)

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

func parseHandlerBodyMethod(s string) *types.HandlerBodyParams {
	var ret types.HandlerBodyParams
	str := strings.Trim(s, "()")
	replacedS := reBracesContent.FindAllString(str, 1)
	if len(replacedS) == 1 {
		trimS := strings.Trim(replacedS[0], "{}")
		replacedM := make(map[string]types.TypeD)
		for _, line := range strings.Split(trimS, ",") {
			kv := strings.Split(line, "=")
			replacedM[trimSpace(kv[0])] = types.TypeD{Kind: trimSpace(kv[1])}
		}
		ret.EmbedReplaces = replacedM
	}
	ret.Name = reBracesContent.ReplaceAllString(str, "")
	return &ret
}
