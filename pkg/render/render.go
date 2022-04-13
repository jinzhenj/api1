package render

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"path/filepath"

	"github.com/go-logr/logr"

	"github.com/go-swagger/pkg/types"
	"github.com/go-swagger/pkg/utils"
)

const (
	swaggerVersion = "2.0"
)

var (
	definitonPrefix      = "#/definitions"
	ErrInvalidDefinition = errors.New("not allowed to define go struct in current dir")
)

type RenderSwagger struct {
	log        logr.Logger
	apiDir     string
	typesDir   []string
	structDefs []types.StructRecord

	//
	structDefsUsedInApi []types.StructRecord
	isInitStructDefs    bool
}

func (o RenderSwagger) RenderSwaggerDoc() (*types.SwaggerDoc, error) {
	var ret types.SwaggerDoc
	ret.Swagger = swaggerVersion
	defs, err := o.BuildSwaggerEntity()
	if err != nil {
		return nil, err
	}
	ret.Definitions = defs

	apis, err := o.BuildSwaggerEndpoint()
	if err != nil {
		return nil, err
	}
	ret.Paths = apis
	return &ret, nil
}

func (o RenderSwagger) GenerateSwaggerJson(fileName string) error {
	doc, err := o.RenderSwaggerDoc()
	if err != nil {
		return err
	}
	if err := utils.MayCreateDir(filepath.Dir(fileName)); err != nil {
		return err
	}

	byteData, err := json.Marshal(doc)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(fileName, byteData, 0755)
}
