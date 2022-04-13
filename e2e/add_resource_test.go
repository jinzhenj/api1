package e2e

import (
	"testing"

	"github.com/stretchr/testify/assert"

	// "github.com/go-swagger/pkg/parser"
	"github.com/go-swagger/pkg/tools"
)

//TODO: check attribute
func TestAddResource(t *testing.T) {
	setUp(t)
	defer tearDown(t)

	t.Run("parser_add_resource", func(t *testing.T) {
		u := tools.NewUpdateApiInterface(testLogger, "pkg/api", "github.com/go-swagger/e2e/testData")
		err := u.Do()
		assert.NoError(t, err)

	})
}
