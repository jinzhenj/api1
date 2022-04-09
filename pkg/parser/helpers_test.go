package parser

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegs(t *testing.T) {
	t.Run("found_struct_def", func(t *testing.T) {
		t1 := "type s struct {"
		t2 := "type s struct{"

		s := reStructPrefix.FindAllString(t1, 1)
		assert.Equal(t, s, []string{t1})

		s = reStructPrefix.FindAllString(t2, 1)
		assert.Equal(t, s, []string{t2})

		t3 := "types struct{"
		t4 := "type sstruct{"

		s = reStructPrefix.FindAllString(t3, 1)
		assert.Nil(t, s)

		s = reStructPrefix.FindAllString(t4, 1)
		assert.Nil(t, s)
	})

	t.Run("found_struct_name", func(t *testing.T) {
		t1 := "type User struct {"
		s := reToken.FindAllString(t1, -1)
		assert.True(t, len(s) > 0)
		assert.Equal(t, s[1], "User")

	})

	t.Run("field_comment", func(t *testing.T) {
		t1 := " // abc"
		s := reIsComment.FindAllString(t1, 1)
		assert.True(t, len(s) > 0)
		assert.Equal(t, s[0], " // abc")

		t2 := "// abc"
		s = reIsComment.FindAllString(t2, 1)
		assert.True(t, len(s) > 0)
		assert.Equal(t, s[0], "// abc")

		s = reExtractComment.FindAllString(t1, -1)
		assert.Equal(t, s[1], " abc")

		s = reExtractComment.FindAllString(t2, -1)
		assert.Equal(t, s[0], " abc")

		t3 := " //"
		s = reExtractComment.FindAllString(t3, -1)
		assert.True(t, len(s) == 1)
	})

	t.Run("found_tag", func(t *testing.T) {
		t1 := "`" + `json:"fake_json,omitempty" binding:"required,omitempty" form:"fake_form,omitempty"` + "`"
		s := reTag.FindAllString(t1, -1)
		assert.True(t, len(s) > 0)
		assert.Equal(t, s[0], t1)

		t2 := "  " + t1 + "  "
		s = reTag.FindAllString(t2, -1)
		assert.True(t, len(s) > 0)
		assert.Equal(t, s[0], t1)

	})

	t.Run("found_handler_doc", func(t *testing.T) {
		t1 := `
		@xxxx 
		@doc(
			summary: hello
		)
		@handler login
		`

		t2 := `

		@doc(
			summary: hello
		)
		@handler getUserInfo
		get /api/user/:id (pkg.config.UserInfoReq) return (pkg.config.UserInfoReply)
		
		`

		expected := `
		@doc(
			summary: hello
		)
		`
		s := reHandlerDoc.FindAllString(t1, -1)
		assert.True(t, len(s) > 0)
		assert.Equal(t, s[0], strings.Trim(expected, " \n\t"))

		s = reHandlerDoc.FindAllString(t2, -1)
		assert.True(t, len(s) == 1)
		assert.Equal(t, s[0], strings.Trim(expected, " \n\t"))
	})

	t.Run("found_handler_handler", func(t *testing.T) {
		t1 := `
		@handler login
	`
		s := reHanderHander.FindAllString(t1, -1)
		assert.True(t, len(s) > 0)
		assert.Equal(t, s[0], strings.Trim(t1, " \n\t"))

	})

	t.Run("found_empty_line_spaces", func(t *testing.T) {

		t1 := `   `
		t2 := ` abc `

		s := reEmptyLineWithSpace.FindAllString(t1, -1)
		assert.True(t, len(s) == 1)
		assert.True(t, t1 == s[0])

		s = reEmptyLineWithSpace.FindAllString(t2, -1)
		assert.True(t, len(s) == 0)
	})

	t.Run("found_http_handler_token", func(t *testing.T) {
		t1 := "  	get /api/user/search (pkg.config.UserSearchReq) returns (pkg.config.UserInfoReply) "

		s := reHandlerToken.FindAllString(t1, -1)
		res := []string{"get", "/api/user/search", "(pkg.config.UserSearchReq)", "returns", "(pkg.config.UserInfoReply)"}
		assert.True(t, len(s) == 5)
		assert.Equal(t, s[0], res[0])
		assert.Equal(t, s[1], res[1])
		assert.Equal(t, s[2], res[2])
		assert.Equal(t, s[3], res[3])
		assert.Equal(t, s[4], res[4])
	})

	t.Run("found_braces_content", func(t *testing.T) {
		t1 := "abc{test1}"
		s := reBracesContent.FindAllString(t1, 1)
		assert.True(t, len(s) == 1)
		assert.Equal(t, s[0], "{test1}")
	})
}

func TestIsComment(t *testing.T) {
	t1 := " // abc"
	assert.True(t, isCommentLine(t1))
}

func TestExtractStructContent(t *testing.T) {
	t.Run("extract_struct", func(t *testing.T) {
		// one struct
		t1 := `
		xxx ni 
	type User struct {
ID int 
 // hello world
Name string` + "`json:\"name\"` \n RoleIDs []int64 \n} ssos"

		res, err := extractBracesBlock(bytes.NewReader([]byte(t1)))
		assert.NoError(t, err)

		assert.Equal(t, 1, len(res))
		assert.Equal(t, 1, strings.Count(res[0], "{"))
		assert.Equal(t, 1, strings.Count(res[0], "}"))

		// two struct
		t2 := t1 + t1
		res, err = extractBracesBlock(bytes.NewReader([]byte(t2)))
		assert.NoError(t, err)

		assert.Equal(t, 2, len(res))
		assert.Equal(t, 1, strings.Count(res[0], "{"))
		assert.Equal(t, 1, strings.Count(res[0], "}"))
	})

}
