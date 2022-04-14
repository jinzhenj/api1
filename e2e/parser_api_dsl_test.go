package e2e

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/go-swagger/pkg/parser"
)

//TODO: check attribute
func TestParserApiFile(t *testing.T) {
	t.Run("parser_api_file", func(t *testing.T) {
		pfObj := parser.NewParserApiFile(testLogger, "testData/pkg/api/user.api")
		res, err := pfObj.ParsrApiDefFile()
		assert.NoError(t, err)
		assert.Equal(t, len(res), 2)

		jsonData, err := json.Marshal(&res)
		assert.NoError(t, err)
		testLogger.V(6).Info("parse go struct file", "json data", string(jsonData))
	})
}
