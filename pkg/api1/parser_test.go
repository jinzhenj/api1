package api1

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
	fmt.Println("Start")

	parser := Parser{}
	t1 := `
	
# @go.package example.com/some/pkg
# @go.import example.com/some/pkg
group user

# this type is datetime
scalar DateTime #datetime is represented by string

# this type is season
enum Season { # a
	# b
	Spring #c
	Summer
	# also called Fall
	Autumn
	Winter

	# b
}

# this type is a person
#	@ts.modifier public
struct Person {

	# person's name
	# @a
	# @b a b c
	# @c:yaml|
	#   asome comment
	# @d:json {"a": "b"}
	# @e:json ["abc"]
	# @f:yml|
	#   a:
	#     b: 1
	#     c: 2
	# @g:json|
	#   {
	#     "name": "john"	
	#   }
	# @param pageSize required "10"
	# @param pageNo optional "1"
	name: string

	# person's age
  # @go.tag json:",omitempty"
	age: int

	# person's address
	# @omitEmpty
	address: string

	# @jsonIgnore
	friends: [string]
}

# this is a interface
interface UserController {

	echoUser(
		): Person

	# aaa
	testFunc(param1: [string], param2: int # bbb

			# ccc
			param3: int, param4: float # ddd

			# eee
			param5: string, param6: int ): [Person] #fff

	createUser(p: Person, ): Person

	currentUser(): Person

	changeNames(names: [string]): [Person]

	modifyUser(
		# this is id
		id: int, # this is id

		# p is person
		p: Person, newName: string # this is newName


	): Person
}

# comment at the end
`

	if schema, err := parser.Parse(t1); err != nil {
		t.Fatalf("Parse error: %v", err)
	} else {
		fmt.Println(dump(schema))
	}

	t2 := `

group t2

struct User {
	name: string

	# @omitempty
	address: string
}

struct User2 {
  name: string

  # @ignore
  password: string
}

interface user {

  # @route get /users/:id
  getUser(id: int): User
}

# @go.type uint64
scalar Timestamp

struct User3 {
  
  # @go.type map[string]string
  properties: object
}

interface user_interface {

  # @route get /users/:id/properties
  # @go.type map[string]string
  getProperties(id: int): object
 
  # @route put /users/:id/properties
  setProperties(
    id: int,

    # @go.type map[string]string
    properties: object
  )
}

struct Type {
	
}

interface TypeController {

  # @route get /types
  listTypes(): [Type]

  # @route put /types/:id
  # @go.middleware adminRequired
  modifyType(id: int, type: Type)
}

struct UserLoginRequest {
  # @go.validator: email
  email: string
  password: string
}

interface LoginController {

  # @route post /actions/checkUserExists
  checkUserExists(
    # @go.validator email
    email: string
  ): boolean
}


	`
	if schema, err := parser.Parse(t2); err != nil {
		t.Fatalf("Parse error: %v", err)
	} else {
		fmt.Println(dump(schema))
	}

	fmt.Println("Done")
	// t.Log("Done")

	var err error

	t3 := `
	  # test function without braces is not valid
	  group t3

		interface T3 {
			validFun(): string
			invalidFun: string
		}
	`
	_, err = parser.Parse(t3)
	t.Log(err)
	assert.Error(t, err)

	t4 := `
  group user

	enum Role {
		ADMIN
		USER
	}

	struct UserLoginRequest {
		# @go.validator: email
		email: string
		password: string
	}
	
	struct UserLoginResponse {
		token: string
	}
	
	struct UserRegisterRequest {
		email: string
		name: string
		role: Role
		password: string
	}
	
	struct ResetPasswordRequest {
		resetPasswordCode: string
		password: string
	}
	
	struct User {
		id: int
		name: string
		email: string
		role: Role
	}
	
	interface UserController {
	 
		# @route post /users/login
		userLogin(req: UserLoginRequest): UserLoginResponse
	
		# @route post /users/logout
		userLogout()
	
		# @route post /users/register
		userRegister(req: UserRegisterRequest)
	
		# @route get /users/activate
		userActivate(token: string)
	
		# @route get /users/reSendActivationEmail
		userReSendActivationEmail(email: string)
	
		# @route get /users/forgetPassword
		userForgetPassword(email: string)
	
		# @route post /users/reSetPassword
		userResetPassword(req: ResetPasswordRequest)
	
		# @route get /users
		currentUser(): User?
	
		# status=StatusConflict, if confliction detected.
		# @route get /users/check
		checkUserExists(email: string, name: string)
	}`

	if schema, err := parser.Parse(t4); err != nil {
		t.Fatalf("Parse error: %v", err)
	} else {
		fmt.Println(dump(schema))
	}
}

func dump(o interface{}) string {
	b, _ := json.MarshalIndent(o, "", "  ")
	return string(b)
}
