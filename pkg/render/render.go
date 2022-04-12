package render

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/go-logr/logr"

	"github.com/go-swagger/pkg/parser"
	"github.com/go-swagger/pkg/types"
	"github.com/go-swagger/pkg/utils"
)

var (
	definitonPrefix      = "#/definitions"
	ErrInvalidDefinition = errors.New("not allowed to define go struct in current dir")
)

func BuildSwaggerFile(log logr.Logger, dir string) (types.SwaggerEndpointStruct, error) {
	filter := func(fn string) bool {
		return strings.HasSuffix(fn, ".api")
	}

	log.V(6).Info("scan dirs", "dir", dir)
	files, err := utils.ListFiles(dir, filter)
	if err != nil {
		log.Error(err, "failed to list files", "dir", dir)
		return nil, err
	}
	log.V(6).Info("build swagger file", "dir", dir, "found api definition files", files)

	// parse api definions from file
	httpHandlers := make([]types.HttpHandler, 0)

	for _, fn := range files {
		tmp, err := parser.ParsrApiDefFile(log, fn)
		if err != nil {
			log.Error(err, "failed to parser api def file", "fileName", fn)
			return nil, err
		}
		httpHandlers = append(httpHandlers, tmp...)
	}
	return generateSwaggerEndpointHandler(log, httpHandlers)
}

// underlay functions
func generateSwaggerEndpointHandler(log logr.Logger, httpHandlers []types.HttpHandler) (types.SwaggerEndpointStruct, error) {
	ret := make(types.SwaggerEndpointStruct)

	// build swagger: step0
	// step1: collect struct definitions
	structDefs, err := parser.ParserStructFromDirs(log, []string{"pkg/"})
	if err != nil {
		return nil, err
	}
	log.V(6).Info("generate swagger handler", "structDefs", structDefs)

	// step2: prepare
	httpHandlerIdxMap := make(map[string][]int, 0)
	swaggerHandlers := make([]types.SwaggerEndpointHandler, len(httpHandlers))
	for i, handler := range httpHandlers {
		swgHandlers, err := httpHandler2SwaggerEndpointHandler(&handler, structDefs)
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

func httpHandler2SwaggerEndpointHandler(handler *types.HttpHandler, structDefs []types.StructRecord) (*types.SwaggerEndpointHandler, error) {
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
	ret.Parameters = httpHandlerReq2SwaggerParameters(handler.Req, structDefs)

	// response
	resp, err := httpHandler2ResSwaggerResponse(handler.Res, structDefs)
	if err != nil {
		return nil, err
	}
	ret.Responses = resp

	return &ret, nil
}

func httpHandlerReq2SwaggerParameters(params *types.HandlerBodyParams, structDefs []types.StructRecord) []types.SwaggerParameters {
	if params == nil {
		return nil
	}
	structDef := foundStructDef(params.Name, structDefs)
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
				if utils.IsGoBuiltinTypes(field.Kind.Kind) {
					param.Schema = &types.EmbedSwaggerItemDef{
						Ref: mapGoTypesToSwagger(field.Kind.Kind),
					}
				} else {
					param.Schema = &types.EmbedSwaggerItemDef{
						Ref: fmt.Sprintf("%s/%s", definitonPrefix, field.Kind.Kind),
					}
				}
			} else if field.Tag.Position == string(types.ParamQueryPositionKind) {
				if utils.IsGoBuiltinTypes(field.Kind.Kind) {
					param.Type = mapGoTypesToSwagger(field.Kind.Kind)
				} else {
					fmt.Printf("struct:(%s) not supported query parameters, kind: %s\n", structDef.Name, field.Kind.Kind)
				}
			} else if field.Tag.Position == string(types.ParamPathPositionKind) {
				param.Type = mapGoTypesToSwagger(field.Kind.Kind)
			} else {
				fmt.Printf("struct:(%s) not supported query parameters, kind: %s\n", structDef.Name, field.Kind.Kind)
			}
			ret = append(ret, param)
		}
	}

	return ret
}

// currently, only status code 200 repsonse
func httpHandler2ResSwaggerResponse(params *types.HandlerBodyParams, structDefs []types.StructRecord) (map[string]types.SwaggerResponseSchema, error) {
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
	if utils.IsGoBuiltinTypes(param) {
		if strings.HasPrefix(param, "[]") {
			tmp := types.SwaggerItemDef{Type: "array", Items: types.EmbedSwaggerItemDef{Type: mapGoTypesToSwagger(strings.Trim(param, "[]"))}}
			data, _ := json.Marshal(&tmp)
			return data, nil
		} else {
			tmp := types.EmbedSwaggerItemDef{Type: mapGoTypesToSwagger(param)}
			data, _ := json.Marshal(&tmp)
			return data, nil
		}
	} else if utils.IsComposedByBuiltin(param) {
		return nil, fmt.Errorf("not supported response param:(%s) when build swagger response", param)
	} else {
		obj, m := parser.ExtractNestedReplacedStruct(param)
		if strings.HasPrefix(obj, "[]") {
			if m == nil {
				tmp := types.SwaggerItemDef{Type: "array", Items: types.EmbedSwaggerItemDef{Type: strings.Trim(obj, "[]")}}
				return json.Marshal(&tmp)

			} else {
				return nil, fmt.Errorf("not supported now, when marshalHttpResponse, param: %s", param)
			}
		} else {
			if m == nil {
				tmp := types.SwaggerParameters{Schema: &types.EmbedSwaggerItemDef{Ref: fmt.Sprintf("%s/%s", definitonPrefix, obj)}}
				jsonData, _ := json.Marshal(&tmp)
				return jsonData, nil
			} else {
				var allOf types.HelperSwaggerAllOf
				arr := make([]json.RawMessage, 0)
				// obj def
				objDef := types.SwaggerObjectDef{Ref: fmt.Sprintf("%s/%s", definitonPrefix, obj)}
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
