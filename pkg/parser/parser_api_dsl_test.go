package parser

import (
	"testing"

	"github.com/go-swagger/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestParserHandler(t *testing.T) {
	t.Run("parser_doc", func(t *testing.T) {
		t1 := `
		@doc(
			summary: 获取用户信息
		)`
		res := parserDoc(testLogger, t1)
		assert.Equal(t, *res, types.HandlerDoc{Summary: "获取用户信息"})
	})

	t.Run("parser_handler_body_method", func(t *testing.T) {
		t1 := `(pkg.config.UserSearchReq)`
		res := parseHandlerBodyMethod(t1)
		assert.Equal(t, *res, types.HandlerBodyParams{Name: "pkg.config.UserSearchReq"})

		t2 := `(pkg.config.UserSearchReq{list=[]string, records=[]int})`
		res = parseHandlerBodyMethod(t2)
		m := make(map[string]types.TypeD)
		m["list"] = types.TypeD{Kind: "[]string"}
		m["records"] = types.TypeD{Kind: "[]int"}
		assert.Equal(t, *res, types.HandlerBodyParams{Name: "pkg.config.UserSearchReq", EmbedReplaces: m})
	})

	t.Run("parser_handler_info", func(t *testing.T) {
		t1 := ` get /api/user/search (pkg.config.UserSearchReq) return (pkg.config.UserInfoReply)  `
		var handler types.HttpHandler
		err := parseHanlderInfo(testLogger, t1, &handler)
		assert.NoError(t, err)
		assert.Equal(t, *handler.Req, types.HandlerBodyParams{Name: "pkg.config.UserSearchReq"})
		assert.Equal(t, *handler.Res, types.HandlerBodyParams{Name: "pkg.config.UserInfoReply"})
		handler.Req = nil
		handler.Res = nil
		assert.Equal(t, handler, types.HttpHandler{Endpoint: "/api/user/search", Method: "get"})

		t2 := ` get /api/user/search  return  (pkg.config.UserInfoReply)  `
		var handler2 types.HttpHandler
		err = parseHanlderInfo(testLogger, t2, &handler2)
		assert.NoError(t, err)
		assert.Equal(t, *handler2.Res, types.HandlerBodyParams{Name: "pkg.config.UserInfoReply"})
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
		res, err := parserHandler(testLogger, t1)
		assert.NoError(t, err)
		assert.Equal(t, *res.Doc, types.HandlerDoc{Summary: "获取用户信息"})
		assert.Equal(t, res.Name, "getUserInfo")
		assert.Equal(t, res.Method, "get")
		assert.Equal(t, res.Endpoint, "/api/user/:id")

		res, err = parserHandler(testLogger, t2)
		assert.NoError(t, err)
		assert.Nil(t, res.Doc)
		assert.Equal(t, res.Name, "getUserInfo")
		assert.Equal(t, res.Method, "get")
		assert.Equal(t, res.Endpoint, "/api/user/:id")
	})

}
