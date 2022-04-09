package types

import "strings"

/*
@doc(
	summary: "用户搜索"
)
@handler searchUser
get /api/user/search (pkg.config.UserSearchReq) returns (pkg.config.UserInfoReply)

*/

type HandlerDoc struct {
	Summary     string
	Description string
}

type HandlerBodyParams struct {
	Name          string
	EmbedReplaces map[string]TypeD // existed embed field that need replaced
}

func (o *HandlerBodyParams) IsThisStruct(relativeFilePath, structName string) bool {
	if relativeFilePath == "./" {
		return o.Name == structName
	}
	return o.Name == strings.ReplaceAll(relativeFilePath, "/", ".")+structName

}

type HttpHandler struct {
	Name     string
	Method   string
	Endpoint string
	Doc      *HandlerDoc
	Req      *HandlerBodyParams
	Res      *HandlerBodyParams
}

func (o *HttpHandler) CapitalName() string {
	return strings.Title(o.Name)
}
