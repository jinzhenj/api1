package api1

import (
	"fmt"
	"strings"

	"github.com/api1/pkg/utils"
	"github.com/pkg/errors"
)

var resID = "[A-Za-z][0-9A-Za-z_]*?"
var rePathParam = utils.Compile("^(?::(%s)|\\{(%s)\\})$", resID, resID)
var reRoute = utils.Compile("(?i)^\\s*(get|put|post|delete|options|head|patch|trace)\\s+(.+?)\\s*$")

type PathStyle int

const (
	PathStyleBrace PathStyle = 0
	PathStyleColon PathStyle = 1
)

var tplByStyle = map[PathStyle]string{
	PathStyleBrace: "{%s}",
	PathStyleColon: ":%s",
}

func ParsePath(s string, style PathStyle) (string, []string) {
	var pathParams []string
	parts := strings.Split(s, "/")
	for i, part := range parts {
		if m := rePathParam.FindStringSubmatch(part); m != nil {
			paramName := m[1] + m[2]
			parts[i] = fmt.Sprintf(tplByStyle[style], paramName)
			pathParams = append(pathParams, paramName)
		}
	}
	return strings.Join(parts, "/"), pathParams
}

// return: method, path, pathParams, err
func ParseRoute(route string, style PathStyle) (string, string, []string, error) {
	var m []string
	if m = reRoute.FindStringSubmatch(route); m == nil {
		return "", "", nil, errors.Errorf("invalid route: %s", route)
	}
	method := strings.ToLower(m[1])
	path, pathParams := ParsePath(m[2], style)
	return method, path, pathParams, nil
}

type Position string

const (
	PositionBody   Position = "body"
	PositionPath   Position = "path"
	PositionQuery  Position = "query"
	PositionHeader Position = "header"
	PositionCookie Position = "cookie"
)

type RouteParam struct {
	Param
	In Position `json:"in"`
}

type TypeKind string

const (
	TypeKindScalar TypeKind = "scalar"
	TypeKindEnum   TypeKind = "enum"
	TypeKindStruct TypeKind = "struct"
)

type RouteParser struct {
	types map[string]TypeKind
}

func (p *RouteParser) LoadSchema(schema *Schema) {
	p.types = make(map[string]TypeKind)
	if schema == nil {
		return
	}
	for _, g := range schema.Groups {
		for _, sc := range g.ScalarTypes {
			p.types[sc.Name] = TypeKindScalar
		}
		for _, en := range g.EnumTypes {
			p.types[en.Name] = TypeKindEnum
		}
		for _, st := range g.StructTypes {
			p.types[st.Name] = TypeKindStruct
		}
	}
}

func (p *RouteParser) IsBodyParam(t *TypeRef) bool {
	if t == nil {
		return false
	}
	if t.ItemType != nil {
		return true
	}
	if t.Name == "object" || t.Name == "any" {
		return true
	}
	if t.Name == "int" || t.Name == "float" ||
		t.Name == "string" || t.Name == "boolean" {
		return false
	}
	kind, ok := p.types[t.Name]
	if !ok {
		return false
	}
	switch kind {
	case TypeKindEnum:
		return false
	case TypeKindStruct:
		return true
	case TypeKindScalar:
		return false // TODO: return false for temp
	}
	return false
}

func (p *RouteParser) ParseParams(iface *Iface, fun *Fun, pathParams []string) ([]RouteParam, error) {
	// init unused path params
	paramInPath := make(map[string]bool)
	for _, pathParam := range pathParams {
		paramInPath[pathParam] = true
	}

	var params []RouteParam
	var bodyParams []string
	for _, param := range fun.Params {
		position := PositionQuery
		if paramInPath[param.Name] {
			position = PositionPath
			delete(paramInPath, param.Name)
		} else if p.IsBodyParam(param.Type) {
			position = PositionBody
			bodyParams = append(bodyParams, param.Name)
		}

		// TODO: check path param is simple type

		if position == PositionPath && param.Type.Nullable {
			return nil, errors.Errorf(
				"Function [%s.%s] has nullable path param [%s]",
				iface.Name, fun.Name, param.Name)
		}

		params = append(params, RouteParam{
			Param: param,
			In:    position,
		})
	}

	if len(paramInPath) > 0 {
		var keys []string
		for key := range paramInPath {
			keys = append(keys, key)
		}
		return nil, errors.Errorf(
			"Function [%s.%s] has undefined path params [%s]",
			iface.Name, fun.Name, strings.Join(keys, ", "))
	}
	if len(bodyParams) > 1 {
		return nil, errors.Errorf(
			"Function [%s.%s] has more than one bodyParams [%s]",
			iface.Name, fun.Name, strings.Join(bodyParams, ", "))
	}

	return params, nil
}
