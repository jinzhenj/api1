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
		s := reStructName.FindAllString(t1, -1)
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

		res, err := extractStructFromStream(bytes.NewReader([]byte(t1)))
		assert.NoError(t, err)

		assert.Equal(t, 1, len(res))
		assert.Equal(t, 1, strings.Count(res[0], "{"))
		assert.Equal(t, 1, strings.Count(res[0], "}"))

		// two struct
		t2 := t1 + t1
		res, err = extractStructFromStream(bytes.NewReader([]byte(t2)))
		assert.NoError(t, err)

		assert.Equal(t, 2, len(res))
		assert.Equal(t, 1, strings.Count(res[0], "{"))
		assert.Equal(t, 1, strings.Count(res[0], "}"))
	})

}
