package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStrings(t *testing.T) {
	assert.Equal(t, "user", CamelCase("user"))
	assert.Equal(t, "user", CamelCase("User"))
	assert.Equal(t, "user", CamelCase("USER"))
	assert.Equal(t, "userName", CamelCase("userName"))
	assert.Equal(t, "userName", CamelCase("UserName"))
	assert.Equal(t, "userName", CamelCase("user_name"))
	assert.Equal(t, "userName", CamelCase("USER_NAME"))
	assert.Equal(t, "userName", CamelCase("User_Name"))
	assert.Equal(t, "userName", CamelCase("UsEr_NaMe"))

	assert.Equal(t, "User", PascalCase("user"))
	assert.Equal(t, "User", PascalCase("User"))
	assert.Equal(t, "User", PascalCase("USER"))
	assert.Equal(t, "UserName", PascalCase("userName"))
	assert.Equal(t, "UserName", PascalCase("UserName"))
	assert.Equal(t, "UserName", PascalCase("user_name"))
	assert.Equal(t, "UserName", PascalCase("USER_NAME"))
	assert.Equal(t, "UserName", PascalCase("User_Name"))
	assert.Equal(t, "UserName", PascalCase("UsEr_NaMe"))
}
