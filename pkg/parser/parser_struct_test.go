package parser

import (
	"testing"

	"github.com/go-swagger/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestParserTag(t *testing.T) {
	t.Run("parser_tag", func(t *testing.T) {
		pfobj := ParserStructFile{log: testLogger, FileName: ""}
		t1 := "`" + `json:"fake_json" binding:"required" form:"-,omitempty"` + "`"
		res := pfobj.parserTag(t1)
		expected := types.Tag{Json: "fake_json", Binding: true}
		assert.Equal(t, *res, expected)

		t2 := "`" + `json:"fake_json,omitempty" binding:"required,omitempty" form:"fake_form,omitempty"` + "`"
		res = pfobj.parserTag(t2)
		expected = types.Tag{Json: "fake_json", Binding: true, Form: "fake_form"}
		assert.Equal(t, *res, expected)
	})
}

func TestParserField(t *testing.T) {
	t.Run("parser_field", func(t *testing.T) {
		pfobj := ParserStructFile{log: testLogger, FileName: "pkg/types/user.go", ModulePrefixName: "pkg.types"}
		t1 := "User string " + "`" + `json:"user" binding:"required" position:"path" form:"user,omitempty"` + "`"
		t2 := "RoleIDs []int"
		t3 := "  RoleIDs []int"
		t4 := " Detail UserDetail"
		t5 := " Details []UserDetail"
		t6 := " Detail config.Detail"

		var ret1 types.Field
		pfobj.parserField(t1, &ret1)
		assert.Equal(t, *ret1.Tag, types.Tag{Json: "user", Binding: true, Form: "user", Position: "path"})
		ret1.Tag = nil
		assert.Equal(t, ret1, types.Field{Name: "User", Kind: types.TypeD{Kind: "string"}})

		var ret2 types.Field
		pfobj.parserField(t2, &ret2)
		assert.Equal(t, *ret2.Tag, types.Tag{Json: "RoleIDs"})
		ret2.Tag = nil
		assert.Equal(t, ret2, types.Field{Name: "RoleIDs", Kind: types.TypeD{Kind: "[]int"}})

		var ret3 types.Field
		pfobj.parserField(t3, &ret3)
		ret3.Tag = nil
		assert.Equal(t, ret3, types.Field{Name: "RoleIDs", Kind: types.TypeD{Kind: "[]int"}})

		var ret4 types.Field
		pfobj.parserField(t4, &ret4)
		ret4.Tag = nil
		assert.Equal(t, ret4, types.Field{Name: "Detail", Kind: types.TypeD{Kind: "pkg.types.UserDetail"}})

		var ret5 types.Field
		pfobj.parserField(t5, &ret5)
		ret5.Tag = nil
		assert.Equal(t, ret5, types.Field{Name: "Details", Kind: types.TypeD{Kind: "[]pkg.types.UserDetail"}})

		var ret6 types.Field
		pfobj.parserField(t6, &ret6)
		ret6.Tag = nil
		assert.Equal(t, ret6, types.Field{Name: "Detail", Kind: types.TypeD{Kind: "pkg.config.Detail"}})
	})

}

func TestParserStruct(t *testing.T) {
	t.Run("parser_struct", func(t *testing.T) {
		pfobj := ParserStructFile{log: testLogger, FileName: ""}
		t1 := `type User struct {  
  // just id
  Id int
  //hello
  //world
  Age int
  } `
		record := pfobj.ParseStruct(t1)
		assert.Equal(t, record.Name, "User")
		assert.True(t, len(record.Fields) == 2)
		record.Fields[0].Tag = nil
		assert.Equal(t, record.Fields[0], types.Field{Name: "Id", Kind: types.TypeD{Kind: "int"}, Comments: " just id"})
		record.Fields[1].Tag = nil
		assert.Equal(t, record.Fields[1], types.Field{Name: "Age", Kind: types.TypeD{Kind: "int"}, Comments: "hello\nworld"})

	})

}
