package golang

import (
	"testing"

	"github.com/jinzhenj/api1/pkg/api1"
	"github.com/stretchr/testify/assert"
)

func TestRender(t *testing.T) {
	r := Render{}

	t1 := api1.EnumType{
		HasName: api1.HasName{Name: "Role"},
		HasComments: api1.HasComments{
			Comments: []string{
				"This is a role",
			},
		},
		Options: []api1.EnumOption{
			{HasName: api1.HasName{Name: "ADMIN"}},
			{HasName: api1.HasName{Name: "USER"}},
		},
	}
	r1 := r.renderEnum(&t1)
	assert.Equal(t, r1.Name, "Role")
	assert.Equal(t, len(r1.Options), 2)
	assert.Equal(t, r1.Options[0].Name, "RoleAdmin")
	assert.Equal(t, r1.Options[1].Name, "RoleUser")

	t2 := api1.StructType{
		HasName: api1.HasName{Name: "User"},
		HasComments: api1.HasComments{
			Comments: []string{
				"This is a user",
			},
		},
		Fields: []api1.StructField{
			{
				HasName: api1.HasName{Name: "name"},
				Type: &api1.TypeRef{
					HasName: api1.HasName{Name: "string"},
				},
			},
			{
				HasName: api1.HasName{Name: "age"},
				Type: &api1.TypeRef{
					HasName: api1.HasName{Name: "int"},
				},
			},
			{
				HasName: api1.HasName{Name: "role"},
				Type: &api1.TypeRef{
					HasName: api1.HasName{Name: "Role"},
				},
			},
			{
				HasName: api1.HasName{Name: "address"},
				HasComments: api1.HasComments{
					Comments: []string{
						"Address is nullable",
					},
					SemComments: map[string]interface{}{
						"omitempty": nil,
					},
				},
				Type: &api1.TypeRef{
					HasName:  api1.HasName{Name: "string"},
					Nullable: true,
				},
			},
			{
				HasName: api1.HasName{Name: "password"},
				HasComments: api1.HasComments{
					SemComments: map[string]interface{}{
						"ignore": nil,
					},
				},
				Type: &api1.TypeRef{
					HasName: api1.HasName{Name: "string"},
				},
			},
		},
	}
	r2 := r.renderStruct(&t2)
	t.Log(r2.Code())

}

func TestRenderEnum(t *testing.T) {
	parser := api1.Parser{}
	r := Render{}

	t1 := `
	  group t1

		enum E1 {
			O1
			O2
		}

		enum E2 {
			O1 = 1,
			O2 = 2,
		}
	`

	schema, err := parser.Parse(t1)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	r1 := r.renderEnum(&schema.Groups[0].EnumTypes[0])
	assert.Equal(t, "E1", r1.Name)
	assert.Equal(t, "string", r1.BaseType.Name)
	assert.Equal(t, 2, len(r1.Options))
	assert.Equal(t, "E1O1", r1.Options[0].Name)
	assert.Equal(t, "E1", r1.Options[0].TypeName)
	assert.Equal(t, "O1", *r1.Options[0].Value.StrVal)
	assert.Equal(t, "E1O2", r1.Options[1].Name)
	assert.Equal(t, "E1", r1.Options[1].TypeName)
	assert.Equal(t, "O2", *r1.Options[1].Value.StrVal)

	r2 := r.renderEnum(&schema.Groups[0].EnumTypes[1])
	assert.Equal(t, "E2", r2.Name)
	assert.Equal(t, "int64", r2.BaseType.Name)
	assert.Equal(t, 2, len(r2.Options))
	assert.Equal(t, "E2O1", r2.Options[0].Name)
	assert.Equal(t, "E2", r2.Options[0].TypeName)
	assert.Equal(t, int64(1), *r2.Options[0].Value.IntVal)
	assert.Equal(t, "E2O2", r2.Options[1].Name)
	assert.Equal(t, "E2", r2.Options[1].TypeName)
	assert.Equal(t, int64(2), *r2.Options[1].Value.IntVal)
}
