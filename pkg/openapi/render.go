package openapi

import (
	"fmt"
	"strings"

	"github.com/api1/pkg/api1"
)

const (
	OpenAPIVersion = "3.0.3"
	refPrefix      = "#/components/schemas/"
	mimeJson       = "application/json"
)

type Render struct {
	rParser *api1.RouteParser
}

func (o *Render) Render(s *api1.Schema) (*OpenAPI, error) {
	o.rParser = &api1.RouteParser{}
	o.rParser.LoadSchema(s)

	var openAPI OpenAPI
	openAPI.OpenAPI = OpenAPIVersion
	openAPI.Info.Title = ""
	openAPI.Info.Version = ""
	c, err := o.renderComponents(s)
	if err != nil {
		return nil, err
	}
	openAPI.Components = *c

	if paths, err := o.renderPaths(s); err != nil {
		return nil, err
	} else {
		openAPI.Paths = paths
	}
	return &openAPI, nil
}

func (o *Render) renderComponents(s *api1.Schema) (*Components, error) {
	c := Components{
		Schemas: make(map[string]Schema),
	}
	any := Schema{
		OneOf: []Schema{
			{Type: "integer"},
			{Type: "number"},
			{Type: "string"},
			{Type: "boolean"},
			{Type: "object"},
			{Type: "array", Items: &Schema{Ref: refPrefix + "any"}},
		},
	}
	c.Schemas["any"] = any
	for _, g := range s.Groups {
		for _, sc := range g.ScalarTypes {
			s := o.renderSchemaScalar(sc)
			if s != nil {
				c.Schemas[sc.Name] = *s
			}
		}
		for _, en := range g.EnumTypes {
			s := o.renderSchemaEnum(en)
			c.Schemas[en.Name] = *s
		}
		for _, st := range g.StructTypes {
			s, err := o.renderSchemaObject(st)
			if err != nil {
				return nil, err
			}
			c.Schemas[st.Name] = *s
		}
	}
	return &c, nil
}

// return schema, required, err
func (o *Render) renderSchemaRef(t *api1.TypeRef) (*Schema, bool) {
	if t.Name != "" {
		if s, ok := tryGetSchema(t.Name); ok {
			return s, !t.Nullable
		}
		// any/scalar/enum/struct
		return &Schema{Ref: fmt.Sprintf("%s%s", refPrefix, t.Name)}, !t.Nullable
	}
	if t.ItemType != nil {
		itemSchema, _ := o.renderSchemaRef(t.ItemType)
		return &Schema{Type: "array", Items: itemSchema}, !t.Nullable
	}
	panic("unreachable")
}

func (o *Render) renderSchemaScalar(sc api1.ScalarType) *Schema {
	if typ, ok := sc.SemComments["openapi.type"].(string); ok {
		return &Schema{Type: typ}
	}
	return nil
}

func (o *Render) renderSchemaEnum(en api1.EnumType) *Schema {
	s := Schema{
		Type:        "string",
		Description: strings.Join(en.Comments, " "),
	}
	if _, ok := en.SemComments["deprecated"]; ok {
		s.Deprecated = true
	}
	for _, o := range en.Options {
		s.Enum = append(s.Enum, o.Name)
		if o.Comments != nil && len(o.Comments) > 0 {
			s.Description +=
				fmt.Sprintf("\n%s: %s", o.Name, strings.Join(o.Comments, " "))
		}
	}
	return &s
}

func (o *Render) renderSchemaObject(st api1.StructType) (*Schema, error) {
	s := Schema{
		Type:        "object",
		Description: strings.Join(st.Comments, " "),
		Properties:  make(map[string]Schema),
	}
	if _, ok := st.SemComments["deprecated"]; ok {
		s.Deprecated = true
	}
	for _, field := range st.Fields {
		if _, ok := field.SemComments["ignore"]; ok {
			continue
		}
		property, required := o.renderSchemaRef(field.Type)
		if required {
			s.Required = append(s.Required, field.Name)
		}
		if property.Ref == "" {
			property.Description = strings.Join(field.Comments, " ")
		}
		s.Properties[field.Name] = *property
	}

	return &s, nil
}

func (o *Render) renderPaths(s *api1.Schema) (Paths, error) {
	paths := make(Paths)
	for _, g := range s.Groups {
		for _, iface := range g.Ifaces {
			for _, fun := range iface.Funs {
				var route string
				var ok bool
				if route, ok = fun.SemComments["route"].(string); !ok {
					continue
				}

				m, path, pathParams, err := api1.ParseRoute(route, api1.PathStyleBrace)
				if err != nil {
					return nil, err
				}

				method := parseMethod(m)
				operation, err := o.renderOperation(&iface, &fun, method, pathParams)
				if err != nil {
					return nil, err
				}
				if paths[path] == nil {
					paths[path] = make(PathItem)
				}
				paths[path][method] = *operation
			}
		}
	}
	return paths, nil
}

func (o *Render) renderOperation(iface *api1.Iface, fun *api1.Fun, method Method, pathParams []string) (*Operation, error) {
	operation := Operation{
		Tags:        []string{iface.Name},
		Description: strings.Join(fun.Comments, " "),
		OperationID: fun.Name,
	}
	if summary, ok := fun.SemComments["summary"].(string); ok {
		operation.Summary = summary
	}
	if _, ok := fun.SemComments["deprecated"]; ok {
		operation.Deprecated = true
	}

	parameters, requestBody, err := o.renderParameters(iface, fun, method, pathParams)
	if err != nil {
		return nil, err
	}
	operation.Parameters = parameters
	operation.RequestBody = requestBody

	responses, err := o.renderResponses(fun.Type)
	if err != nil {
		return nil, err
	}
	operation.Responses = responses

	return &operation, nil
}

func (o *Render) renderParameters(iface *api1.Iface, fun *api1.Fun, method Method, pathParams []string) ([]Parameter, *RequestBody, error) {
	var parameters []Parameter
	var requestBody *RequestBody

	routeParams, err := o.rParser.ParseParams(iface, fun, pathParams)
	if err != nil {
		return nil, nil, err
	}

	for _, param := range routeParams {
		s, required := o.renderSchemaRef(param.Type)

		if param.In == api1.PositionBody {
			content := make(map[string]MediaType)
			content[mimeJson] = MediaType{Schema: s}
			requestBody = &RequestBody{
				Description: strings.Join(param.Comments, " "),
				Content:     content,
				Required:    required,
			}
			continue
		}

		p := Parameter{
			Name:        param.Name,
			In:          parsePosition(string(param.In)),
			Description: strings.Join(param.Comments, " "),
			Required:    required,
			Schema:      s,
		}
		if _, ok := param.SemComments["deprecated"]; ok {
			p.Deprecated = true
		}
		parameters = append(parameters, p)
	}

	return parameters, requestBody, nil
}

// currently, only status code 200 repsonse
func (o *Render) renderResponses(t *api1.TypeRef) (Responses, error) {
	content := make(map[string]MediaType)
	if t != nil {
		s, _ := o.renderSchemaRef(t)
		content[mimeJson] = MediaType{Schema: s}
	}
	responses := make(Responses)
	responses["200"] = Response{
		Description: "Default Response",
		Content:     content,
	}
	return responses, nil
}
