package openapi

// see: https://swagger.io/specification/

type OpenAPI struct {
	OpenAPI    string     `json:"openapi"` // required
	Info       Info       `json:"info"`    // required
	Servers    []Server   `json:"servers,omitempty"`
	Paths      Paths      `json:"paths"` // required
	Components Components `json:"components,omitempty"`
}

type Info struct {
	Title       string `json:"title"`   // required
	Version     string `json:"version"` // required
	Description string `json:"description,omitempty"`
}

type Server struct {
	Url       string                    `json:"url"` // required
	Variables map[string]ServerVariable `json:"variables"`
}

type ServerVariable struct {
	Default string `json:"default"` // required
}

type Paths map[string]PathItem
type PathItem map[Method]Operation
type Method string

const (
	MethodGet     Method = "get"
	MethodPut     Method = "put"
	MethodPost    Method = "post"
	MethodDelete  Method = "delete"
	MethodOptions Method = "options"
	MethodHead    Method = "head"
	MethodPatch   Method = "patch"
	MethodTrace   Method = "trace"
)

type Operation struct {
	Tags        []string     `json:"tags,omitempty"`
	Summary     string       `json:"summary,omitempty"`
	Description string       `json:"description,omitempty"`
	OperationID string       `json:"operationId,omitempty"`
	Parameters  []Parameter  `json:"parameters,omitempty"`
	RequestBody *RequestBody `json:"requestBody,omitempty"`
	Responses   Responses    `json:"responses"` // required
	Deprecated  bool         `json:"deprecated,omitempty"`
}

type Parameter struct {
	Name        string   `json:"name"` // required
	In          Position `json:"in"`   // required
	Description string   `json:"description,omitempty"`
	Required    bool     `json:"required,omitempty"` // required if in path
	Deprecated  bool     `json:"deprecated,omitempty"`
	Schema      *Schema  `json:"schema,omitempty"`
}

type Position string

const (
	PositionPath   Position = "path"
	PositionQuery  Position = "query"
	PositionHeader Position = "header"
	PositionCookie Position = "cookie"
)

type RequestBody struct {
	Description string               `json:"description,omitempty"`
	Content     map[string]MediaType `json:"content"` // required
	Required    bool                 `json:"required,omitempty"`
}

type Responses map[string]Response

type Response struct {
	Description string               `json:"description"` // required
	Content     map[string]MediaType `json:"content,omitempty"`
}

type MediaType struct {
	Schema *Schema `json:"schema,omitempty"`
}

type Components struct {
	Schemas map[string]Schema `json:"schemas,omitempty"`
}

type Schema struct {
	Ref string `json:"$ref,omitempty"`

	Type        string `json:"type,omitempty"`
	Format      string `json:"format,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Default     string `json:"default,omitempty"`

	// object
	Properties    map[string]Schema `json:"properties,omitempty"`
	Required      []string          `json:"required,omitempty"`
	Deprecated    bool              `json:"deprecated,omitempty"`
	MinProperties *int              `json:"minProperties,omitempty"`
	MaxProperties *int              `json:"maxProperties,omitempty"`

	// array
	Items       *Schema `json:"items,omitempty"`
	MinItems    *int    `json:"minItems,omitempty"`
	MaxItems    *int    `json:"maxItems,omitempty"`
	UniqueItems *int    `json:"uniqueItems,omitempty"`

	// string
	Enum      []string `json:"enum,omitempty"`
	Pattern   string   `json:"pattern,omitempty"`
	MinLength *int     `json:"minLength,omitempty"`
	MaxLength *int     `json:"maxLength,omitempty"`

	// integer|number
	MultipleOf       *float64 `json:"multipleOf,omitempty"`
	Minimum          *float64 `json:"minimum,omitempty"`
	Maximum          *float64 `json:"maximum,omitempty"`
	ExclusiveMinimum *float64 `json:"exclusiveMinimum,omitempty"`
	ExclusiveMaximum *float64 `json:"exclusiveMaximum,omitempty"`

	// other
	AllOf []Schema `json:"allOf,omitempty"`
	OneOf []Schema `json:"oneOf,omitempty"`
	AnyOf []Schema `json:"anyOf,omitempty"`
}
