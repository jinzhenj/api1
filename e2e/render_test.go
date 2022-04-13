package e2e

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/go-swagger/pkg/render"
)

func TestRenderApi(t *testing.T) {
	err := os.Chdir("./testData")
	assert.NoError(t, err)
	rObj := render.NewRenderSwagger(testLogger, "pkg/api", nil)

	swgObj, err := rObj.BuildSwaggerEndpoint()
	assert.NoError(t, err)

	jsonData, _ := json.Marshal(&swgObj)
	testLogger.V(6).Info("render_api", "json data", string(jsonData))
}

func TestRenderEntity(t *testing.T) {
	err := os.Chdir("./testData")
	assert.NoError(t, err)
	rObj := render.NewRenderSwagger(testLogger, "pkg/api", nil)

	swgObj, err := rObj.BuildSwaggerEntity()
	assert.NoError(t, err)

	jsonData, _ := json.Marshal(&swgObj)
	testLogger.V(6).Info("render_api", "json data", string(jsonData))
}
