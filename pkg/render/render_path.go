package render

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-logr/logr"

	"github.com/go-swagger/pkg/parser"
	"github.com/go-swagger/pkg/tools"
	"github.com/go-swagger/pkg/types"
)

func NewRenderSwagger(log logr.Logger, apiDir string, typesDir []string) *RenderSwagger {
	return &RenderSwagger{log: log, apiDir: apiDir, typesDir: typesDir}
}

func (o RenderSwagger) BuildSwaggerEndpoint() (types.SwaggerEndpointsStruct, error) {
	filter := func(fn string) bool {
		return strings.HasSuffix(fn, ".api")
	}
	httpHandlers, err := tools.ListHttpHandlers(o.log, o.apiDir, filter)
	if err != nil {
		return nil, err
	}
	return o.generateSwaggerEndpointHandler(httpHandlers)
}

// underlay functions
func (o RenderSwagger) buildStructDefs() ([]types.StructRecord, error) {
	if !o.isInitStructDefs {
		structDefs, err := parser.ParserStructFromDirs(o.log, []string{"pkg/"})
		if err != nil {
			return nil, err
		}
		o.structDefs = structDefs
	}
	return o.structDefs, nil
}

func (o RenderSwagger) generateSwaggerEndpointHandler(httpHandlers []types.HttpHandler) (types.SwaggerEndpointsStruct, error) {
	ret := make(types.SwaggerEndpointsStruct)

	// build swagger: step0
	// step1: collect struct definitions
	structDefs, err := o.buildStructDefs()
	if err != nil {
		return nil, err
	}
	o.log.V(6).Info("generate swagger handler", "structDefs", structDefs)

	// step2: prepare
	httpHandlerIdxMap := make(map[string][]int, 0)
	swaggerHandlers := make([]types.SwaggerEndpointHandler, len(httpHandlers))
	for i, handler := range httpHandlers {
		swgHandlers, err := o.httpHandler2SwaggerEndpointHandler(&handler, structDefs)
		if err != nil {
			return nil, err
		}
		swaggerHandlers[i] = *swgHandlers
		if val, ok := httpHandlerIdxMap[handler.Endpoint]; ok {
			tmp := make([]int, 0)
			tmp = append(tmp, val...)
			tmp = append(tmp, i)
			httpHandlerIdxMap[handler.Endpoint] = tmp
		} else {
			httpHandlerIdxMap[handler.Endpoint] = []int{i}
		}
	}

	// stepN:
	for key, arr := range httpHandlerIdxMap {
		pathInfo := make(types.SwaggerSingleResourceApi)
		for _, idx := range arr {
			pathInfo[swaggerHandlers[idx].Method] = swaggerHandlers[idx]
		}
		ret[key] = pathInfo
	}

	return ret, nil
}

func (o RenderSwagger) httpHandler2SwaggerEndpointHandler(handler *types.HttpHandler, structDefs []types.StructRecord) (*types.SwaggerEndpointHandler, error) {
	ret := types.SwaggerEndpointHandler{
		Produces: []string{"application/json"},
		Tags:     []string{handler.Resource},
		Method:   handler.Method,
		Endpoint: handler.Endpoint,
	}

	// doc
	if handler.Doc != nil {
		ret.Summary = handler.Doc.Summary
	}

	// parameters
	ret.Parameters = o.httpHandlerReq2SwaggerParameters(handler.Req, structDefs)

	// response
	resp, err := o.httpHandler2ResSwaggerResponse(handler.Res, structDefs)
	if err != nil {
		return nil, err
	}
	ret.Responses = resp

	return &ret, nil
}

func (o RenderSwagger) httpHandlerReq2SwaggerParameters(params *types.HandlerBodyParams, structDefs []types.StructRecord) []types.SwaggerParameters {
	if params == nil {
		return nil
	}
	//TODO: not expose val here
	structDef := foundStructDef(params.Kind.Val, structDefs)
	if structDef == nil {
		return nil
	}

	ret := make([]types.SwaggerParameters, 0)
	for _, field := range structDef.Fields {
		if field.Tag != nil && field.Tag.Position != "" {
			param := types.SwaggerParameters{
				Description: field.Comments,
				Name:        field.Tag.Json,
				In:          field.Tag.Position,
				Required:    field.Tag.Binding,
			}
			// currently, neither array kind nor map kind parameter in http request
			// TODO: using switch
			if field.Tag.Position == string(types.ParamBodyPositionKind) {
				if field.Kind.IsBuiltin() {
					param.Schema = &types.EmbedSwaggerItemDef{
						Ref: mapGoTypesToSwagger(field.Kind.GetKind()),
					}
				} else {
					param.Schema = &types.EmbedSwaggerItemDef{
						Ref: getRef(&field.Kind),
					}
				}
			} else if field.Tag.Position == string(types.ParamQueryPositionKind) {
				if field.Kind.IsBuiltin() {
					param.Type = mapGoTypesToSwagger(field.Kind.GetKind())
				} else {
					o.log.Info("not supported query params", "struct def", structDef.Name, "kind", field.Kind.GetKind())
				}
			} else if field.Tag.Position == string(types.ParamPathPositionKind) {
				param.Type = mapGoTypesToSwagger(field.Kind.GetKind())
			} else {
				o.log.Info("not supported query params", "struct def", structDef.Name, "kind", field.Kind.GetKind())
			}
			ret = append(ret, param)
		}
	}

	return ret
}

// currently, only status code 200 repsonse
func (o RenderSwagger) httpHandler2ResSwaggerResponse(params *types.HandlerBodyParams, structDefs []types.StructRecord) (map[string]types.SwaggerResponseSchema, error) {
	if params == nil {
		return nil, nil
	}
	ret := make(map[string]types.SwaggerResponseSchema)
	var record types.SwaggerResponseSchema
	rawMsg, err := marshalHttpResponse(strings.Trim(params.Value, "()"))
	if err != nil {
		return nil, err
	}
	record.Schema = rawMsg
	ret["200"] = record
	return ret, nil
}

//too complicated !!!
func marshalHttpResponse(param string) (json.RawMessage, error) {
	if types.IsGoBuiltinTypes(param) {
		if strings.HasPrefix(param, "[]") {
			tmp := types.SwaggerItemDef{Type: "array", Items: &types.EmbedSwaggerItemDef{Type: mapGoTypesToSwagger(strings.Trim(param, "[]"))}}
			data, _ := json.Marshal(&tmp)
			return data, nil
		} else {
			tmp := types.EmbedSwaggerItemDef{Type: mapGoTypesToSwagger(param)}
			data, _ := json.Marshal(&tmp)
			return data, nil
		}
	} else if types.IsComposedByBuiltin(param) {
		return nil, fmt.Errorf("not supported response param:(%s) when build swagger response", param)
	} else {
		obj, m := parser.ExtractNestedReplacedStruct(param)
		kind := types.TypeD{Val: obj}
		if kind.IsArray() {
			if m == nil {
				tmp := types.SwaggerItemDef{
					Type:  "array",
					Items: &types.EmbedSwaggerItemDef{Ref: getRef(&kind)}}
				return json.Marshal(&tmp)
			} else {
				return nil, fmt.Errorf("not supported now, when marshalHttpResponse, param: %s", param)
			}
		} else {
			if m == nil {
				tmp := types.SwaggerParameters{Schema: &types.EmbedSwaggerItemDef{Ref: getRef(&kind)}}
				jsonData, _ := json.Marshal(&tmp)
				return jsonData, nil
			} else {
				var allOf types.HelperSwaggerAllOf
				arr := make([]json.RawMessage, 0)
				// obj def
				objDef := types.SwaggerObjectDef{Ref: getRef(&kind)}
				jsonData, _ := json.Marshal(&objDef)
				arr = append(arr, jsonData)
				// properties related
				properties := types.HelperSwaggerProperties{Type: "object"}
				propertiesMap := make(map[string]json.RawMessage)
				for key, value := range m {
					rawMsg, err := marshalHttpResponse(value)
					if err != nil {
						return nil, err
					}
					propertiesMap[key] = rawMsg
				}
				properties.Properties = propertiesMap
				jsonData, _ = json.Marshal(&properties)
				arr = append(arr, jsonData)
				allOf.AllOf = arr
				return json.Marshal(&allOf)

			}
		}
	}
}
