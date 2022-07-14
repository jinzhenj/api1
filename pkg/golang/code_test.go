package golang

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCodeGen(t *testing.T) {
	t1 := GoEnumCodeBlock{
		Comments: []string{
			"This is user role",
		},
		Name: "Role",
		Options: []GoEnumOption{
			{
				Comments: []string{
					"This is admin",
				},
				Name:     "RoleAdmin",
				TypeName: "Role",
				Value:    "ADMIN",
			},
			{
				Comments: []string{
					"This is normal user",
				},
				Name:     "RoleUser",
				TypeName: "Role",
				Value:    "USER",
			},
		},
	}
	exp1 := `// This is user role
type Role string

const (
  // This is admin
  RoleAdmin Role = "ADMIN"
  // This is normal user
  RoleUser Role = "USER"
)
`
	assert.Equal(t, exp1, t1.Code())

	t2 := GoStructType{
		Comments: []string{
			"This is user",
			"Many operations require user logged in",
		},
		Name: "User",
		Fields: []GoStructField{
			{
				Comments: []string{
					"This is user name",
				},
				Name: "Name",
				Type: &GoType{Name: "string"},
				Tags: map[string]string{
					"json": "name",
				},
			},
			{
				Comments: []string{
					"This is user age",
				},
				Name: "Age",
				Type: &GoType{Name: "int"},
				Tags: map[string]string{
					"json": "age",
				},
			},
		},
	}
	exp2 := "// This is user\n" +
		"// Many operations require user logged in\n" +
		"type User struct {\n" +
		"  // This is user name\n" +
		"  Name string `json:\"name\"`\n" +
		"  // This is user age\n" +
		"  Age int `json:\"age\"`\n" +
		"}\n"
	assert.Equal(t, exp2, t2.Code())

	t3 := GoFunction{
		Comments: []string{
			"This is a function",
		},
		Name: "GetSomeString",
		Params: []GoParam{
			{
				Name: "param1",
				Type: &GoType{Name: "int"},
			},
			{
				Name: "param2",
				Type: &GoType{
					KeyType: &GoType{Name: "int"},
					ItemType: &GoType{
						ItemType: &GoType{Name: "string"},
					},
				},
			},
		},
		RetTypes: []GoType{{
			ItemType: &GoType{
				Name: "string",
			},
		}},
	}
	exp3 := `// This is a function
func GetSomeString(param1 int, param2 map[int][]string) []string {
}
`

	assert.Equal(t, exp3, t3.Code())

	t4 := GoInterface{
		Comments: []string{
			"This is an interface",
		},
		Name: "UserController",
		Functions: []GoFunction{
			{
				InIface: true,
				Comments: []string{
					"This is a function",
				},
				Name: "GetSomeString",
				Params: []GoParam{
					{
						Name: "param1",
						Type: &GoType{Name: "int"},
					},
					{
						Name: "param2",
						Type: &GoType{
							KeyType: &GoType{Name: "int"},
							ItemType: &GoType{
								ItemType: &GoType{Name: "string"},
							},
						},
					},
				},
				RetTypes: []GoType{{
					ItemType: &GoType{
						Name: "string",
					},
				}},
			},
			{
				InIface: true,
				Comments: []string{
					"This is another function",
				},
				Name: "GetAnInt",
				Params: []GoParam{
					{
						Name: "param1",
						Type: &GoType{Name: "int", IsPointer: true},
					},
				},
				RetTypes: []GoType{{
					Name: "int",
				}},
			},
		},
	}
	exp4 := `// This is an interface
type UserController interface {

  // This is a function
  GetSomeString(param1 int, param2 map[int][]string) []string

  // This is another function
  GetAnInt(param1 *int) int
}
`

	assert.Equal(t, exp4, t4.Code())

	file := GoFile{
		Name:    "user.go",
		Package: "api",
		Imports: []string{},
		CodeGens: []CodeGen{
			&t1,
			&t2,
			&t4,
		},
	}
	exp5 := "package api\n\n" +
		exp1 + "\n" +
		exp2 + "\n" +
		exp4

	assert.Equal(t, exp5, file.Code())
}
