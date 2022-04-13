package render

import (
	"fmt"

	"github.com/go-swagger/pkg/types"
	"github.com/go-swagger/pkg/utils"
)

const (
	disabledToRenderJsonTag = "-"
)

var example = `
"api.Response": {
	"type": "object",
	"properties": {
		"code": {
			"type": "integer"
		},
		"data": {}, --// sik
		"msg": {
			"type": "string"
		},
		"media_type": {
			"description": "文件类型",
			"$ref": "#/definitions/model.MediaType"
		}
	}
},
"model.BatchRequest": {
	"type": "object",
	"properties": {
		"ids": {
			"type": "array",
			"items": {
				"type": "integer"
			}
		}
	}
},

`

func (o RenderSwagger) BuildSwaggerEntity() (map[string]types.SwaggerObjectDef, error) {
	defs, err := o.buildStructDefs()
	if err != nil {
		return nil, err
	}
	ret := make(map[string]types.SwaggerObjectDef)

	for _, def := range defs {
		entity, err := o.buildEntity(def)
		if err != nil {
			return nil, err
		}
		ret[def.RelativePathName] = entity
	}

	return ret, nil
}

func (o RenderSwagger) buildEntity(def types.StructRecord) (types.SwaggerObjectDef, error) {
	ret := types.SwaggerObjectDef{
		Type: "object",
	}
	propertiesM := make(map[string]types.SwaggerItemDef)
	for _, field := range def.Fields {
		if field.Tag != nil && field.Tag.Json == disabledToRenderJsonTag {
			continue
		}
		if !field.Kind.IsSupported() {
			o.log.Info("build entity, field kind is not supported", "kind", field.Kind.Kind)
			continue
		}
		var property types.SwaggerItemDef

		if field.Kind.IsArray() {
			property = types.SwaggerItemDef{
				Type: "array",
			}
			if utils.IsGoBuiltinTypes(field.Kind.GetKind()) {
				if !utils.IsInterface(field.Kind.GetKind()) {
					property.Items = &types.EmbedSwaggerItemDef{Type: field.Kind.GetKind()}
				}
			} else {
				//TODO: use functions
				property.Items = &types.EmbedSwaggerItemDef{Ref: fmt.Sprintf("%s/%s", definitonPrefix, field.Kind.GetKind())}
			}
		} else {
			if utils.IsGoBuiltinTypes(field.Kind.GetKind()) {
				if !utils.IsInterface(field.Kind.GetKind()) {
					property.Type = field.Kind.GetKind()
				}
			} else {
				property.Ref = fmt.Sprintf("%s/%s", definitonPrefix, field.Kind.GetKind())
			}
		}
		property.Description = field.Comments
		propertiesM[field.Tag.Json] = property
	}

	ret.Properties = propertiesM
	return ret, nil
}
