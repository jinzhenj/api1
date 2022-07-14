package openapi

import (
	"strings"
)

func tryGetSchema(typeName string) (*Schema, bool) {
	switch typeName {
	case "int":
		return &Schema{Type: "integer"}, true
	case "float":
		return &Schema{Type: "number"}, true
	case "string":
		return &Schema{Type: "string"}, true
	case "boolean":
		return &Schema{Type: "boolean"}, true
	case "object":
		return &Schema{Type: "object"}, true
	}
	return nil, false
}

func parseMethod(s string) Method {
	switch strings.ToLower(s) {
	case string(MethodGet):
		return MethodGet
	case string(MethodPut):
		return MethodPut
	case string(MethodPost):
		return MethodPost
	case string(MethodDelete):
		return MethodDelete
	case string(MethodOptions):
		return MethodOptions
	case string(MethodHead):
		return MethodHead
	case string(MethodPatch):
		return MethodPatch
	case string(MethodTrace):
		return MethodTrace
	}
	panic("unexpected")
}

func parsePosition(s string) Position {
	switch s {
	case string(PositionPath):
		return PositionPath
	case string(PositionQuery):
		return PositionQuery
	case string(PositionHeader):
		return PositionHeader
	case string(PositionCookie):
		return PositionCookie
	}
	panic("unexpected")
}
