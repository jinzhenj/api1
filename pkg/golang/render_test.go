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
