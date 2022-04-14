package parser

import (
	"testing"

	"github.com/go-swagger/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestParserHandler(t *testing.T) {
	pfObj := ParserApiFile{log: testLogger, fileName: ""}

	t.Run("parser_doc", func(t *testing.T) {
		t1 := `
		@doc(
			summary: 获取用户信息
		)`
		res := pfObj.parserDoc(t1)
		assert.Equal(t, *res, types.HandlerDoc{Summary: "获取用户信息"})
	})

	t.Run("parser_handler_body_method", func(t *testing.T) {

		t1 := `(pkg.config.UserSearchReq)`
		res := pfObj.parseHandlerBodyMethod(t1)
		assert.Equal(t, *res.Kind, types.TypeD{Val: "pkg.config.UserSearchReq"})
		res.Kind = nil
		assert.Equal(t, *res, types.HandlerBodyParams{
			Value:        t1,
			RelatedNames: map[string]bool{"pkg.config.UserSearchReq": true}})

		t2 := `(pkg.config.UserSearchReq{list=[]string, records=[]int})`
		res = pfObj.parseHandlerBodyMethod(t2)
		m := make(map[string]types.TypeD)
		m["list"] = types.TypeD{Val: "[]string"}
		m["records"] = types.TypeD{Val: "[]int"}
		assert.Equal(t, *res.Kind, types.TypeD{Val: "pkg.config.UserSearchReq"})
		res.Kind = nil
		assert.Equal(t, *res, types.HandlerBodyParams{
			Value:        t2,
			RelatedNames: map[string]bool{"pkg.config.UserSearchReq": true, "[]string": true, "[]int": true}})
	})

	t.Run("parser_handler_info", func(t *testing.T) {
		t1 := ` get /api/user/search (pkg.config.UserSearchReq) return (pkg.config.UserInfoReply)  `
		var handler types.HttpHandler
		err := pfObj.parseHanlderInfo(t1, &handler)
		assert.NoError(t, err)
		handler.Req.Kind = nil
		handler.Res.Kind = nil
		assert.Equal(t, *handler.Req, types.HandlerBodyParams{
			Value:        "(pkg.config.UserSearchReq)",
			RelatedNames: map[string]bool{"pkg.config.UserSearchReq": true}})
		assert.Equal(t, *handler.Res, types.HandlerBodyParams{
			Value:        "(pkg.config.UserInfoReply)",
			RelatedNames: map[string]bool{"pkg.config.UserInfoReply": true}})
		handler.Req = nil
		handler.Res = nil
		assert.Equal(t, handler, types.HttpHandler{Endpoint: "/api/user/search", Method: "get"})

		t2 := ` get /api/user/search  return  (pkg.config.UserInfoReply)  `
		var handler2 types.HttpHandler
		err = pfObj.parseHanlderInfo(t2, &handler2)
		assert.NoError(t, err)
		handler2.Res.Kind = nil
		assert.Equal(t, *handler2.Res, types.HandlerBodyParams{
			Value:        "(pkg.config.UserInfoReply)",
			RelatedNames: map[string]bool{"pkg.config.UserInfoReply": true}})
		handler2.Res = nil
		assert.Equal(t, handler2, types.HttpHandler{Endpoint: "/api/user/search", Method: "get"})
	})

	t.Run("parser_handler", func(t *testing.T) {

		t1 := `

		@doc(
			summary: 获取用户信息
		)
		@handler getUserInfo
		get /api/user/:id (pkg.config.UserInfoReq) return (pkg.config.UserInfoReply)
		
		`
		t2 := `

		@handler getUserInfo
		get /api/user/:id (pkg.config.UserInfoReq) return (pkg.config.UserInfoReply)
		
		`
		res, err := pfObj.parserHandler(t1)
		assert.NoError(t, err)
		assert.Equal(t, *res.Doc, types.HandlerDoc{Summary: "获取用户信息"})
		assert.Equal(t, res.Name, "getUserInfo")
		assert.Equal(t, res.Method, "get")
		assert.Equal(t, res.Endpoint, "/api/user/:id")

		res, err = pfObj.parserHandler(t2)
		assert.NoError(t, err)
		assert.Nil(t, res.Doc)
		assert.Equal(t, res.Name, "getUserInfo")
		assert.Equal(t, res.Method, "get")
		assert.Equal(t, res.Endpoint, "/api/user/:id")
	})

	t.Run("parser_api_def", func(t *testing.T) {
		t1 := `
		service user {

			@doc(
				summary: 注册
			)
			@handler register
			post /api/user/register (pkg.types.RegisterReq)
	  
		} sss
		`
		res, err := pfObj.ParserApiDef(t1)
		assert.NoError(t, err)
		assert.Equal(t, len(res), 1)
		assert.Equal(t, res[0].Name, "register")
		assert.Equal(t, res[0].Method, "post")
		assert.Equal(t, res[0].Resource, "user")
		assert.Equal(t, res[0].Endpoint, "/api/user/register")
	})
}

func TestExtractNestedReplacedStruct(t *testing.T) {
	t1 := `pkg.config.UserSearchReq`
	t2 := `pkg.config.UserSearchReq{list=[]string , records=[]int}`
	t3 := `pkg.config.UserSearchReq{list=[]string , records=[]int , pager=Paged{count=int}}`

	name, m := ExtractNestedReplacedStruct(t1)
	assert.Equal(t, name, "pkg.config.UserSearchReq")
	assert.Nil(t, m)

	name, m = ExtractNestedReplacedStruct(t2)
	assert.Equal(t, name, "pkg.config.UserSearchReq")
	assert.Equal(t, m, map[string]string{"list": "[]string", "records": "[]int"})

	name, m = ExtractNestedReplacedStruct(t3)
	assert.Equal(t, name, "pkg.config.UserSearchReq")
	assert.Equal(t, m, map[string]string{"list": "[]string", "records": "[]int", "pager": "Paged{count=int"})

}
