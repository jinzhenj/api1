package render

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/go-swagger/pkg/types"
)

func TestRenderHandlers(t *testing.T) {
	t.Run("parser_tag", func(t *testing.T) {
		handlers := []types.HttpHandler{
			{Resource: "user", Endpoint: "/users", Name: "postUsers", Method: "post"},
			{Resource: "user", Endpoint: "/users", Name: "getUsers", Method: "get"},
			{Resource: "user", Endpoint: "/users/:id", Name: "getUser", Method: "get"},
			{Resource: "role", Endpoint: "/roles", Name: "getRoles", Method: "get"},
		}

		res, err := generateSwaggerEndpointHandler(testLogger, handlers)
		assert.NoError(t, err)
		assert.Equal(t, len(res), 3)

		jsonData, _ := json.Marshal(&res)
		testLogger.V(6).Info("render handlers", "serilizatized", string(jsonData))
	})
}
