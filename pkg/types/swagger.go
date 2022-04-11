package types

import "encoding/json"

type SwaggerEndpointStruct struct {
	Endpoints map[string]SwaggerEndpointHandler `json:"endpoints,omitempty"`
}

type SwaggerEndpointHandler struct {
	Produces   []string                         `json:"produces,omitempty"`
	Tags       []string                         `json:"tags,omitempty"`
	Summary    string                           `json:"summary,omitempty"`
	Parameters []SwaggerParameters              `json:"parameters,omitempty"`
	Responses  map[string]SwaggerResponseSchema `json:"responses,omitempty"`
}

type SwaggerResponseSchema struct {
	Description string          `json:"description,omitempty"`
	Schema      json.RawMessage `json:"schema,omitempty"`
}

// swagger http request related
type SwaggerParameters struct {
	Type        string              `json:"type,omitempty"`
	Description string              `json:"description,omitempty"`
	Name        string              `json:"name,omitempty"`
	In          string              `json:"in,omitempty"`
	Required    bool                `json:"required,omitempty"`
	Schema      EmbedSwaggerItemDef `json:"schema,omitempty"`
}

// swagger definitions related
type SwaggerObjectDef struct {
	Ref        string                    `json:"$ref,omitempty"`
	Type       string                    `json:"type,omitempty"`
	Properties map[string]SwaggerItemDef `json:"properties,omitempty"`
}

type SwaggerItemDef struct {
	EmbedSwaggerItemDef `json:"embed_swagger_item_def,omitempty"`
	Items               map[string]EmbedSwaggerItemDef `json:"items,omitempty"`
}

type EmbedSwaggerItemDef struct {
	Type        string `json:"type,omitempty"`
	Description string `json:"description,omitempty"`
	Ref         string `json:"$ref,omitempty"`
}
