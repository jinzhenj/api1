package parser

import (
	"testing"

	"github.com/go-swagger/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestParserTag(t *testing.T) {
	t.Run("parser_tag", func(t *testing.T) {
		t1 := "`" + `json:"fake_json" binding:"required" form:"-,omitempty"` + "`"
		res := parserTag(t1)
		expected := types.Tag{Json: "fake_json", Binding: true}
		assert.Equal(t, *res, expected)

		t2 := "`" + `json:"fake_json,omitempty" binding:"required,omitempty" form:"fake_form,omitempty"` + "`"
		res = parserTag(t2)
		expected = types.Tag{Json: "fake_json", Binding: true, Form: "fake_form"}
		assert.Equal(t, *res, expected)
	})
}

func TestParserField(t *testing.T) {
	t.Run("parser_field", func(t *testing.T) {
		t1 := "User string " + "`" + `json:"user" binding:"required" form:"user,omitempty"` + "`"
		t2 := "RoleIDs []int"
		t3 := "  RoleIDs []int"

		var ret1 types.Field
		parserField(t1, &ret1)
		assert.Equal(t, *ret1.Tag, types.Tag{Json: "user", Binding: true, Form: "user"})
		ret1.Tag = nil
		assert.Equal(t, ret1, types.Field{Name: "User", Kind: "string"})

		var ret2 types.Field
		parserField(t2, &ret2)
		assert.Equal(t, *ret2.Tag, types.Tag{Json: "RoleIDs"})
		ret2.Tag = nil
		assert.Equal(t, ret2, types.Field{Name: "RoleIDs", Kind: "int", IsArray: true})

		var ret3 types.Field
		parserField(t3, &ret3)
		ret3.Tag = nil
		assert.Equal(t, ret2, types.Field{Name: "RoleIDs", Kind: "int", IsArray: true})
	})

}

func TestParserStruct(t *testing.T) {
	t.Run("parser_struct", func(t *testing.T) {
		t1 := `type User struct {  
  // just id
  Id int
  //hello
  //world
  Age int
  } `
		record := ParseStruct("test/", t1)
		assert.Equal(t, record.Name, "User")
		assert.Equal(t, record.SInfo, types.SourceInfo{FileName: "test/"})
		assert.True(t, len(record.Fields) == 2)
		record.Fields[0].Tag = nil
		assert.Equal(t, record.Fields[0], types.Field{Name: "Id", Kind: "int", Comments: " just id"})
		record.Fields[1].Tag = nil
		assert.Equal(t, record.Fields[1], types.Field{Name: "Age", Kind: "int", Comments: "hello\nworld"})

	})

}
