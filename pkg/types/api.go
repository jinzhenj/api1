package types

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
	Kind         *TypeD
	Value        string
	RelatedNames map[string]bool
}

type HttpHandler struct {
	Resource string
	Name     string
	Method   string
	Endpoint string
	Doc      *HandlerDoc
	Req      *HandlerBodyParams
	Res      *HandlerBodyParams
}
