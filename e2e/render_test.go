package e2e

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/go-swagger/pkg/render"
)

//TODO: check attribute
func TestRenderApi(t *testing.T) {
	setUp(t)
	defer tearDown(t)

	rObj := render.NewRenderSwagger(testLogger, "pkg/api", nil)

	swgObj, err := rObj.BuildSwaggerEndpoint()
	assert.NoError(t, err)

	jsonData, _ := json.Marshal(&swgObj)
	testLogger.V(6).Info("render_api", "json data", string(jsonData))
}

// TODO:check attribute
func TestRenderEntity(t *testing.T) {
	setUp(t)
	defer tearDown(t)

	rObj := render.NewRenderSwagger(testLogger, "pkg/api", nil)

	swgObj, err := rObj.BuildSwaggerEntity()
	assert.NoError(t, err)

	jsonData, _ := json.Marshal(&swgObj)
	testLogger.V(6).Info("render_api", "json data", string(jsonData))
}

func TestRenderDoc(t *testing.T) {
	setUp(t)
	defer tearDown(t)

	rObj := render.NewRenderSwagger(testLogger, "pkg/api", nil)
	err := rObj.GenerateSwaggerJson("docs/swagger.json")
	assert.NoError(t, err)
}
