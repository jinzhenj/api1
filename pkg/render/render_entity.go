package render

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/go-swagger/pkg/types"
)

const (
	disabledToRenderJsonTag = "-"
)

func (o RenderSwagger) BuildSwaggerEntity() (types.SwaggerEntitiesStruct, error) {
	defs, err := o.buildStructDefs()
	if err != nil {
		return nil, err
	}
	ret := make(types.SwaggerEntitiesStruct)

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
			o.log.Info("build entity, field kind is not supported", "kind", field.Kind.Val)
			continue
		}
		var property types.SwaggerItemDef

		if field.Kind.IsArray() {
			property = types.SwaggerItemDef{
				Type: "array",
			}
			if field.Kind.IsBuiltin() {
				if !types.IsInterface(field.Kind.GetKind()) {
					property.Items = &types.EmbedSwaggerItemDef{Type: field.Kind.GetKind()}
				}
			} else {
				//TODO: use functions
				property.Items = &types.EmbedSwaggerItemDef{Ref: getRef(&field.Kind)}
			}
		} else {
			if field.Kind.IsBuiltin() {
				if !types.IsInterface(field.Kind.GetKind()) {
					property.Type = field.Kind.GetKind()
				}
			} else {
				property.Ref = getRef(&field.Kind)
			}
		}
		property.Description = field.Comments
		propertiesM[field.Tag.Json] = property
	}

	ret.Properties = propertiesM
	return ret, nil
}

func getRef(f *types.TypeD) string {
	if f.IsBuiltin() {
		return f.GetKind()
	} else {
		ref := strings.Replace(filepath.Join(f.GetModule(), f.GetKind()), "/", ".", -1)
		return fmt.Sprintf("%s/%s", definitonPrefix, ref)
	}

}
